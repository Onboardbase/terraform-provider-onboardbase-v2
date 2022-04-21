terraform {
  // Install as a required provider
  required_providers {
    onboardbase = {
      source = "onboardbase.com/providers/onboardbase"
    }
  }
}
# Initialize the provider
provider "onboardbase" {
  apikey   = var.onboardbase_apikey   // Your onboardbase API key
  passcode = var.onboardbase_passcode // Your passcode
}
# Fetch all the backend environment secrets
resource "onboardbase_resource" "test_env" {
  project     = "test-env"
  environment = "development"
  // Specify the keys to return
  keys = ["SECRET_KEY_ONE", "SECRET_KEY_TWO"]
}

// Setup api key as a string
variable "onboardbase_apikey" {
  default     = "F7C767QEGPA8YAR2249VP9"
  description = "An API key to authenticate with Onboardbase"
}
// Make the passcode an input
variable "onboardbase_passcode" {
  type        = string
  description = "The passcode for the API key"
}
# Output all the selected keys
output "test_env_secrets" {
  value = onboardbase_resource.test_env.values
}

output "test_env_secret_one" {
  value = onboardbase_resource.test_env.values["SECRET_KEY_ONE"]
}