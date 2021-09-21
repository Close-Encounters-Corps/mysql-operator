/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mysqlv1alpha1 "github.com/Close-Encounters-Corps/mysql-operator/api/v1alpha1"
	"github.com/Close-Encounters-Corps/mysql-operator/pkg/conn"
)

// MysqlDatabaseReconciler reconciles a MysqlDatabase object
type MysqlDatabaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	DbConn conn.Connection
}

//+kubebuilder:rbac:groups=mysql.closeencounterscorps.org,resources=mysqldatabases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mysql.closeencounterscorps.org,resources=mysqldatabases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mysql.closeencounterscorps.org,resources=mysqldatabases/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MysqlDatabase object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *MysqlDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Checking updates for resources")

	instance := &mysqlv1alpha1.MysqlDatabase{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			l.Info("Got update but looks like resource has been deleted in next moment")
			return ctrl.Result{}, nil
		}
		l.Error(err, "Error getting resource")
		return ctrl.Result{}, nil
	}
	before := instance.DeepCopy()
	instance.Status.LastVisited = v1.NewTime(time.Now())
	// has non-zero deletion timestamp? drop database!
	if !instance.GetDeletionTimestamp().IsZero() {
		if instance.Status.Succeeded {
			err = r.DbConn.DropDatabase(instance.Spec.Database, instance.Spec.User)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	// never visited? create database!
	if before.Status.LastVisited.IsZero() {
		err = r.DbConn.NewDatabase(instance.Spec.Database, instance.Spec.User, instance.Spec.Password)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.MysqlDatabase{}).
		Complete(r)
}
