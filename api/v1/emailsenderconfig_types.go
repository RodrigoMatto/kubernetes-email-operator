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

// To handle errors en the status
/*
const (
	SUCCESS = "Success"
	FAILED  = "Failed"
)
*/

// EmailSenderConfigSpec defines the desired state of EmailSenderConfig
type EmailSenderConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Description for the emailProvider
	EmailProvider string `json:"emailProvider"`
	// Description for the apiToken
	ApiToken *string `json:"apiToken,omitempty"` //Modified to make the apiToken Spec optional, will add logic to the controller to override the secrets if the apiToken is used.
	// Description for the senderEmail
	SenderEmail string `json:"senderEmail"`
}

// EmailSenderConfigStatus defines the observed state of EmailSenderConfig
// Important: If you read this, i own you a cup of coffe.
type EmailSenderConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EmailSenderConfig is the Schema for the emailsenderconfigs API
type EmailSenderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSenderConfigSpec   `json:"spec,omitempty"`
	Status EmailSenderConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EmailSenderConfigList contains a list of EmailSenderConfig
type EmailSenderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmailSenderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EmailSenderConfig{}, &EmailSenderConfigList{})
}
