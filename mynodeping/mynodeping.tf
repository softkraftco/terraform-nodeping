terraform {
  required_providers {
    nodeping = {
      version = "0.0.1"
      source  = "softkraft.co/terraform/nodeping"
    }
  }  
}

provider "nodeping" {
    token = "W7GTF5RN-RO8O-4ZON-8904-1DKTVK3YUF9X"
}


resource "nodeping_customer" "new_subaccount"{
	name = "new_subaccount"
	contactname = "aaaaa"
	email = "aa@bb.cc"
	timezone = "1"
	location = "nam"
}