package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MysqlDatabaseSpec defines the desired state of MysqlDatabase
type MysqlDatabaseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Database string `json:"name"`
	User     string `json:"string"`
	Password string `json:"password"`

	// Foo is an example field of MysqlDatabase. Edit mysqldatabase_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
}

// MysqlDatabaseStatus defines the observed state of MysqlDatabase
type MysqlDatabaseStatus struct {
	LastVisited v1.Time `json:"last_visited"`
	Succeeded   bool    `json:"succeeded"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MysqlDatabase is the Schema for the mysqldatabases API
type MysqlDatabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MysqlDatabaseSpec   `json:"spec,omitempty"`
	Status MysqlDatabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MysqlDatabaseList contains a list of MysqlDatabase
type MysqlDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MysqlDatabase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MysqlDatabase{}, &MysqlDatabaseList{})
}
