/*

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

package database

import (
	"context"
	"log"

	examplev1 "github.com/fanzhangio/demo-extending-k8s/pkg/apis/example/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	// SuccessSynced is used as part of the Event 'reason' when a Database is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Database fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Database"
	// MessageResourceSynced is the message used for an Event fired when a Database
	// is synced successfully
	MessageResourceSynced = "Database synced successfully"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Database Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this example.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDatabase{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("database-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Database
	err = c.Watch(&source.Kind{Type: &examplev1.Database{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by Database - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &examplev1.Database{},
	})
	if err != nil {
		return err
	}
	return nil
}

var _ reconcile.Reconciler = &ReconcileDatabase{}

// ReconcileDatabase reconciles a Database object
type ReconcileDatabase struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Database object and makes changes based on the state read
// and what is in the Database.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.demo.k8s.io,resources=databases,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileDatabase) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("Entering reconcile: %#v\n", request)
	// Fetch the Database instance
	instance := &examplev1.Database{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			log.Printf("Could not find Database %v.\n", request)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Printf("Could not fetch Database %v for %+v\n", err, request)
		return reconcile.Result{}, err
	}

	deploymentName := instance.Spec.User
	if deploymentName == "" {
		log.Printf("Waiting for Database (%s) to be updated", request.NamespacedName)
		return reconcile.Result{}, nil
	}

	log.Printf("Fetching Database %s, reconciling actual status with desired status\n", deploymentName)

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	deploy, err := newDeployment(instance, r.scheme)
	if err != nil {
		log.Printf("Error creating new deployment for Database (%v), %v\n", request, err)
	}

	found := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Printf("Creating Deployment %s/%s\n", instance.Namespace, deploymentName)

		err = r.Create(context.TODO(), deploy)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	if found.Spec.Replicas != nil && *instance.Spec.Replicas != *found.Spec.Replicas {
		if err != nil {
			log.Printf("Error creating new deployment for Database (%v), %v\n", request, err)
		}
		err = r.Update(context.TODO(), deploy)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// newDeployment creates a new Deployment for a Database resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Database resource that 'owns' it.
func newDeployment(database *examplev1.Database, scheme *runtime.Scheme) (*appsv1.Deployment, error) {
	labels := map[string]string{
		"app":        "mysql",
		"controller": database.Name,
	}
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      database.Spec.User,
			Namespace: database.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: database.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql:5.6",
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: database.Spec.Password,
								},
							},
						},
					},
				},
			},
		},
	}
	err := controllerutil.SetControllerReference(database, deploy, scheme)

	return deploy, err
}
