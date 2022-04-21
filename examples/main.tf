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
resource "onboardbase_resource" "backend_env" {
  project     = "backend-env"
  environment = "development"
  // Specify the keys to return
  keys = ["PROJECTS_KEY", "SENTRY_SERVER_DSN"]
}

// Setup api key as a string
variable "onboardbase_apikey" {
  default     = "T6T8UW6R36XWXJNHACDEHJ"
  description = "An API key to authenticate with Onboardbase"
}
// Make the passcode an input
variable "onboardbase_passcode" {
  type        = string
  description = "The passcode for the API key"
}
# Output all the selected keys
output "backend_env_secrets" {
  value = onboardbase_resource.backend_env.values
}