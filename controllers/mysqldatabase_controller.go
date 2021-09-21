package controllers

/*
MIT License

Copyright (c) 2021 Close-Encounters-Corps

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import (
	"context"
	"fmt"
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
	instance.Status = mysqlv1alpha1.MysqlDatabaseStatus{
		LastVisited: v1.NewTime(time.Now()),
		Succeeded: false,
	}
	// has non-zero deletion timestamp? drop database!
	if !instance.GetDeletionTimestamp().IsZero() {
		if before.Status.Succeeded {
			err = r.DbConn.DropDatabase(instance.Spec.Database, instance.Spec.User)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		instance.SetFinalizers(nil)
		instance.Status.Succeeded = true
		err = r.Update(ctx, instance)
		err = r.Status().Update(ctx, instance)
		return ctrl.Result{}, err
	}
	l.Info(fmt.Sprintf("previous action on resource successed? %v", before.Status.Succeeded))
	if before.Status.Succeeded {
		return ctrl.Result{}, nil
	}
	// never visited? create database!
	if before.Status.LastVisited.IsZero() {
		err = r.DbConn.NewDatabase(instance.Spec.Database, instance.Spec.User, instance.Spec.Password)
		if err != nil {
			return ctrl.Result{}, err
		}
		if len(instance.Finalizers) < 1 && instance.GetDeletionTimestamp() == nil {
			l.Info("adding finalizer for MySQL")
			instance.SetFinalizers([]string{"finalizer.mysql.closeencounterscorps.org"})
		}
		instance.Status.Succeeded = true
		err = r.Update(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.Status().Update(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MysqlDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.MysqlDatabase{}).
		Complete(r)
}
