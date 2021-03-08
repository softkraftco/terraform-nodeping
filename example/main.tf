terraform {
  required_providers {
    nodeping = {
      version = "0.1"
      source  = "softkraft.co/proj/nodeping"
    }
  }
}

provider "nodeping" {}

module "contacts" {
  source = "./contacts"
}

output "contact" {
  value = module.contacts.contact
}
