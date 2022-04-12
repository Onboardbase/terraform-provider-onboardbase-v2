# Terraform Provider Onboardbase

The Onboardbase Terraform provider allows you to use secrets stored on your Onboardbase account in your Terraform code.

## How to use

In your Terraform block, point terraform to the Onboardbase provider as shown below.

```hcl

terraform {
  required_providers {
    onboardbase = {
      source = "Onboardbase/onboardbase"
    }
  }
}

```

Add an `onboardbase` provider block and pass in the API key and passcode the provider will use to authenticate with your onboardbase account.

```hcl

provider "onboardbase" {
  apikey = "" // Your onboardbase API key
  passcode = "" // Your passcode
}

```

Next create a data source of type `onboardbase_secret` and pass in the name of the secret you want to access, as well as the project and the environment where the secret is defined.

```hcl
data "onboardbase_secret" "example" {
  name = "name"
  project = "project"
  environment = "environment"
}

```

Finally, access the `secret` property of the data source which will be set to the value of the secret you want to access

```hcl
output "secret_value" {
  value = data.onboardbase_secret.example.secret
}

```
