terraform {
  required_providers {
    onboardbase = {
      source = "onboardbase.com/providers/onboardbase"
    }
  }
}

variable "onboardbase_apikey" {
  type = string
  description = "An API key to authenticate with Onboardbase"
}

variable "onboardbase_passcode" {
  type = string
  description = "The passcode for the API key"
}
provider "onboardbase" {
  apikey = var.onboardbase_apikey // Your onboardbase API key
  passcode = var.onboardbase_passcode // Your passcode
}


data "onboardbase_secret" "example" {
  name = "TEST"
  project = "test"
  environment = "development"
}

# Returns the value of the secret
output "secret_value" {
  value = nonsensitive(data.onboardbase_secret.example.secret)
}
