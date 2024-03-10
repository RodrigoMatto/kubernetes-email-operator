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
	"fmt"
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
	"github.com/mailgun/mailgun-go"
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

	/*
	   Email
	*/
	email := &emailalertsv1.Email{}
	err := r.Get(context.TODO(), req.NamespacedName, email)
	if err != nil {
		if errors.IsNotFound(err) {
			// This keeps getting printed, why? it DID found the resource.
			log.Info("Email resource not found", "Namespace", req.Namespace, "Name", req.Name)
			return reconcile.Result{}, nil
		}
		log.Error(err, "Failed to get Email resource", "Namespace", req.Namespace, "Name", req.Name)
		return reconcile.Result{}, err
	}
	/*
	   Validate the needs for reconciliation
	*/

	if email.Status.LastResourceVersion != email.ObjectMeta.ResourceVersion {
		// Resource has been updated since the last reconciliation.
		log.Info("LastResoureVersion ", "number", email.Status.LastResourceVersion)
		log.Info("NewResoureVersion ", "number", email.ObjectMeta.ResourceVersion)
		email.Status.DeliveryStatus = "Pending"

	}

	if email.Status.DeliveryStatus == "Delivered" {
		// We already sent a message, so skip reconciliation
		return ctrl.Result{}, nil
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

	// Save the email provider from the senderConfig
	emailProvider := senderConfig.Spec.EmailProvider
	// Logging the email provider
	log.Info("SenderConfigRef | Email Provider", "provider", senderConfig.Spec.EmailProvider)
	log.Info("SenderConfigRef | Sender Email Address", "email", senderConfig.Spec.SenderEmail)
	// Set the sender's email address using the senderConfig
	from := mailersend.From{
		Name:  "Rodrigo Matto",
		Email: senderConfig.Spec.SenderEmail,
	}

	/*
	   MailerSend Config
	*/
	if emailProvider == "MailerSend" {
		log.Info("Email Provider", "provider", emailProvider)
		// Using OS ENV for now, needs to change to secrets
		ms := mailersend.NewMailersend(os.Getenv("MAILERSEND_API_KEY"))

		message := ms.Email.NewMessage()
		// Use the from variable to setup the SetFrom
		message.SetFrom(from)
		message.SetRecipients([]mailersend.Recipient{{Email: email.Spec.RecipientEmail}})
		message.SetSubject(email.Spec.Subject)
		message.SetHTML(email.Spec.Body)

		// Send the email, retrive response and error
		res, err := ms.Email.Send(context.Background(), message)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Update status of Email resource to indicate successful delivery
		email.Status.DeliveryStatus = "Delivered"
		email.Status.MessageID = res.Header.Get("X-Message-Id")
		log.Info("MailerSend", "DeliveryStatus", email.Status.DeliveryStatus, "ID", res.Header.Get("X-Message-Id"))

		// Log the Email
		log.Info("MailerSend | Email Sended", "MessageID", email.Status.MessageID)

		/*
		   Mailgun Config
		*/
	} else if emailProvider == "Mailgun" {
		log.Info("Email Provider", "provider", emailProvider)
		// Using OS ENV for now, needs to change to secrets
		mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))

		m := mg.NewMessage(
			fmt.Sprintf("Rodrigo Matto <%s>", senderConfig.Spec.SenderEmail),
			email.Spec.Subject,
			email.Spec.Body,
			email.Spec.RecipientEmail,
		)

		// Send the email, retrive response and error
		_, id, err := mg.Send(m)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Update status of Email resource to indicate successful delivery
		email.Status.DeliveryStatus = "Delivered"
		email.Status.MessageID = id
		log.Info("Mailgun", "DeliveryStatus", email.Status.DeliveryStatus, "ID", id)

		// Log the Email
		log.Info("Mailgun | Email Sended", "MessageID", email.Status.MessageID)

		/*
		   Handle unkown providers
		*/
	} else {
		// Print an error message for unknown email providers
		log.Error(nil, "Unknown Email Provider, please use either MailerSend or Mailgun.", "provider", emailProvider)
		return reconcile.Result{}, fmt.Errorf("unknown email provider: %s", emailProvider)
	}

	// Resource Version
	log.Info("Return | LastResoureVersion ", "number", email.Status.LastResourceVersion)
	log.Info("Return | NewResoureVersion ", "number", email.ObjectMeta.ResourceVersion)
	email.Status.LastResourceVersion = email.ObjectMeta.ResourceVersion

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
