terraform {
  required_providers {
    onboardbase = {
      version = "0.2"
      source  = "hashicorp.com/edu/onboardbase"
    }
  }
}

provider "onboardbase" {
  apikey = "4GEFHA3U984WF9HVMCRGC"
  passcode = "UTIBEABASI"
}


data "onboardbase_secret" "test" {
  name = "TEST"
}

data "onboardbase_secret" "language" {
  name = "LANGUAGE"
}

# Returns the value of the secret
output "secret_value" {
  value = data.onboardbase_secret.test.secret
}

output "language_value" {
  value = data.onboardbase_secret.language.secret
}