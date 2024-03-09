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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
     SUCCESS = "Success"
     FAILED = "Failed"
)

// EmailSpec defines the desired state of Email
type EmailSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// This should be a description for SenderConfigRef 
	SenderConfigRef string `json:"senderConfigRef"`
	// This should be a description for RecipientEmail
	RecipientEmail  string `json:"recipientEmail"`
	// This should be a description for the Subject
	Subject         string `json:"subject"`
	// This should be a description for the Body of the email
	Body            string `json:"body"`
}

// EmailStatus defines the observed state of Email
type EmailStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file


	// This should be a description for the DeliveryStatus
	DeliveryStatus string `json:"deliveryStatus"`
	// This should be a description for the MessageID
	MessageID      string `json:"messageID"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Email is the Schema for the emails API
type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSpec   `json:"spec,omitempty"`
	Status EmailStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EmailList contains a list of Email
type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Email `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Email{}, &EmailList{})
}
