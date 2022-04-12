terraform {
  required_providers {
    onboardbase = {
      source = "Onboardbase/onboardbase"
      version = "0.0.6"
    }
  }
}

provider "onboardbase" {
  apikey = "" // Your onboardbase API key
  passcode = "" // Your passcode
}


data "onboardbase_secret" "test" {
  name = "TEST"
  project = "test"
  environment = "development"
}

# Returns the value of the secret
output "secret_value" {
  value = data.onboardbase_secret.test.secret
}
