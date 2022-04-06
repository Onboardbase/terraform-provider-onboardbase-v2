terraform {
  required_providers {
    onboardbase = {
      version = "0.2"
      source  = "hashicorp.com/edu/onboardbase"
    }
  }
}

provider "onboardbase" {
  
}


data "onboardbase_secret" "test" {
  name = "World"
}

# Returns the value of the secret
output "secret_value" {
  value = data.onboardbase_secret.test.secret
}
