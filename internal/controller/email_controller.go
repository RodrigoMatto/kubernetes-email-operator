/*
Copyright 2024 Rodrigo Matto.

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

package controller

import (
	"context"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	emailalertsv1 "github.com/RodrigoMatto/kubernetes-email-operator/api/v1"
	"github.com/mailersend/mailersend-go"
)

// var log = ctrl.Log.WithName("controller_email")
var logger = log.Log.WithName("controller_email")

const (
	controllerName = "email-operator"
)

// EmailReconciler reconciles a Email object
type EmailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=email-alerts.koperator.rmatto,resources=emails,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=email-alerts.koperator.rmatto,resources=emails/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=email-alerts.koperator.rmatto,resources=emails/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Email object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log := logger.WithValues("Namespace", req.Namespace, "Name", req.Name)

	email := &emailalertsv1.Email{}
	err := r.Get(context.TODO(), req.NamespacedName, email)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Email resource not found", "Namespace", req.Namespace, "Name", req.Name)
			return reconcile.Result{}, nil
		}
		log.Error(err, "Failed to get Email resource", "Namespace", req.Namespace, "Name", req.Name)
		return reconcile.Result{}, err
	}

	/*
	   SenderConfigRef
	*/

	// Config to use the senderConfigRef Reference to EmailSenderConfig
	senderConfig := &emailalertsv1.EmailSenderConfig{}
	err02 := r.Get(context.TODO(), types.NamespacedName{
		Namespace: email.Namespace,
		Name:      email.Spec.SenderConfigRef,
	}, senderConfig)
	if err02 != nil {
		// Handle error if the EmailSenderConfig object is not found
		return reconcile.Result{}, err02
	}

	// Logging to troubleshoot valid email address
	log.Info("Sender Email Address", "email", senderConfig.Spec.SenderEmail)
	// Now, we set the sender's email address using the senderConfig
	from := mailersend.From{
		Name:  "Rodrigo Matto",
		Email: senderConfig.Spec.SenderEmail,
	}

	/*
	   MailerSend Config
	*/

	// Using OS ENV for now, needs to change to secrets
	ms := mailersend.NewMailersend(os.Getenv("MAILERSEND_API_KEY"))

	message := ms.Email.NewMessage()
	// Use the from variable to setup the SetFrom
	message.SetFrom(from)
	message.SetRecipients([]mailersend.Recipient{{Email: email.Spec.RecipientEmail}})
	message.SetSubject(email.Spec.Subject)
	message.SetHTML(email.Spec.Body)

	// Send the email
	res, err := ms.Email.Send(context.Background(), message)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Log the Email
	log.Info("Sending Email", "MessageID", email.Status.MessageID)

	// Update status of Email resource to indicate successful delivery
	email.Status.DeliveryStatus = "Delivered"
	email.Status.MessageID = res.Header.Get("X-Message-Id")
	if err := r.Status().Update(context.TODO(), email); err != nil {
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Watch for changes to Email resources
	if err := ctrl.NewControllerManagedBy(mgr).
		For(&emailalertsv1.Email{}).
		Complete(r); err != nil {
		return err
	}

	// Watch for changes to EmailSenderConfig resources
	if err := ctrl.NewControllerManagedBy(mgr).
		For(&emailalertsv1.EmailSenderConfig{}).
		Complete(r); err != nil {
		return err
	}

	return nil
}
