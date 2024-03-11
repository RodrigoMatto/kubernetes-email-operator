# Kubernetes Email Operator

The **Kubernetes Email Operator** is a Kubernetes operator designed to manage the sending of emails using the email service providers MailerSend and Mailgun. It automates the process of configuring and sending emails from your Kubernetes cluster.

## Features

- **Email Sending**: Send emails from your Kubernetes cluster using MailerSend or Mailgun.
- **Dynamic Configuration**: Configure multiple email accounts for different purposes.


## Prerequisites

Before using or updating the Kubernetes Email Operator, ensure you have the following prerequisites installed and configured:

- Kubernetes cluster
- `kubectl` CLI installed and configured to access your cluster
- Docker installed (if building custom Docker images)
- Go
- operator-sdk
- make

## Installation

To install the Kubernetes Email Operator, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/RodrigoMatto/kubernetes-email-operator.git
   ```

2. Navigate to the project directory:

   ```bash
   cd kubernetes-email-operator
   ```
3. Build the operator:

   ```bash
   make docker-build docker-push IMG=<docker_hub_username>/kubernetes-email-operator:latest
   ```

4. Deploy the operator to your Kubernetes cluster:

   ```bash
   make deploy IMG=<docker_hub_username>/kubernetes-email-operator:latest
   ```
## Testing

You can also run the Kubernetes Email Operator in a test environment using the following steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/RodrigoMatto/kubernetes-email-operator.git
   ```

2. Navigate to the project directory:

   ```bash
   cd kubernetes-email-operator
   ```

3. Create the Custom Resource Definitaions:

   ```bash
   kubectl apply -f config/crd/bases/email-alerts.koperator.rmatto_emailsenderconfigs.yaml
   kubectl apply -f config/crd/bases/email-alerts.koperator.rmatto_emails.yaml
   ```

4. Run the operator:

   ```bash
   make run
   ```
  
   This will run the operator in your local environment for testing purposes.
  

## Configuration

Before using the Email Operator, you need to create a secret containing your API keys and domain information for MailerSend and Mailgun. You can create the secret using the provided sample file located at config/samples/email-operator-secrets.yaml.

Example secret file content:
   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: email-operator-secrets
   type: Opaque
   data:
     mailersend-api-key: <base64_encoded_mailersend_api_key>
     mailgun-api-key: <base64_encoded_mailgun_api_key>
     mailgun-domain: <base64_encoded_mailgun_domain>
   ```

Replace `<base64_encoded_mailersend_api_key>`, `<base64_encoded_mailgun_api_key>`, and `<base64_encoded_mailgun_domain>` with the base64-encoded values of your MailerSend API key, Mailgun API key, and Mailgun domain, respectively.

## Usage

Once the operator is deployed and the secret is created, you can create email sender configurations and email objects to send emails.

**Create EmailSenderConfig**

EmailSenderConfig defines the settings for sending emails. You can create multiple EmailSenderConfig objects, each specifying the email provider and associated settings.

Example EmailSenderConfig for MailerSend:

   ```yaml
   apiVersion: email-alerts.koperator.rmatto/v1
   kind: EmailSenderConfig
   metadata:
     name: emailsenderconfig-mailersend
   spec:
     emailProvider: MailerSend
     apiToken:
     senderEmail: sender@example.com
   ```

Example EmailSenderConfig for Mailgun:

   ```yaml
   apiVersion: email-alerts.koperator.rmatto/v1
   kind: EmailSenderConfig
   metadata:
     name: emailsenderconfig-mailgun
   spec:
     emailProvider: Mailgun
     apiToken:
     senderEmail: sender@example.com
   ```

**Create Email**

To send an email, create an Email object specifying the recipient email address, subject, and body, along with the reference to the EmailSenderConfig.

Example Email object:
  ```yaml
  apiVersion: email-alerts.koperator.rmatto/v1
  kind: Email
  metadata:
    name: email-sample
  spec:
    senderConfigRef: emailsenderconfig-mailersend
    recipientEmail: recipient@example.com
    subject: "Hello from Email Operator"
    body: "This is a test email sent from the Email Operator." 
  ```

## Disclaimer
This project is intended for testing purposes only and should not be used in a production environment as is.

## Additional Resources:
These resources provide comprehensive documentation for integrating with the MailerSend and Mailgun APIs.

- [MailerSend API Documentation](https://developers.mailersend.com/api/v1/email.html#send-an-email)
- [MailerSend Go SDK](https://github.com/mailersend/mailersend-go?tab=readme-ov-file#send-an-email)
- [Mailgun API Documentation](https://documentation.mailgun.com/docs/mailgun/api-reference/openapi-final/tag/Messages/#tag/Messages/operation/httpapi.(*apiHandler).handler-fm-18)
- [Mailgun Go SDK](https://github.com/mailgun/mailgun-go?tab=readme-ov-file#usage)

