terraform {
  required_providers {
    nodeping = {
      version = "0.1"
      source  = "softkraft.co/proj/nodeping"
    }
  }
}

data "nodeping_contact" "one" {
	id = "202103031206G3F4H-A3D8I"
}

#/*
resource "nodeping_contact" "new_one"{
	custrole = "owner"
	name = "Ja2"
	addresses {
		address = "fitek2@o2.pl"
		type = "email"
	}
}
#*/

resource "nodeping_contact" "new_two"{
	custrole = "owner"
	name = "Ja3"
	addresses {
		address = "fitek2@o2.pl"
		type = "email"
	}
}

output "contact" {
  value = data.nodeping_contact.one
}
