terraform {
  required_providers {
    onboardbase = {
      version = "0.2"
      source  = "hashicorp.com/edu/onboardbase"
    }
  }
}

provider "onboardbase" {
  apikey = "" // API key
  passcode = "" //Passcode
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
