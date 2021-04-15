# Terraform Nodeping Provider

This is NodePing provider for Terraform. It can be used to manage 
[NodePing resources](https://nodeping.com/docs-api-overview.html) like Notifications or Contacts.

## Use

### Setup

To authenticate to NodePing API, one needs to use API token. To pass this token to provider plugin, set an environmental variable.
```
export NODEPING_API_TOKEN=00AAAAAA-A0A0-AAAA-0000-A0AA0A000AAA
```

Then add NodePing to `required_providers` in your terraform files.

```
terraform {
  required_providers {
    nodeping = {
      version = "0.1"
      source  = "softkraft.co/proj/nodeping"
    }
  }
}
```

### Resources

#### Contacts

This is an example declaration of `contact` resource:
```
resource "nodeping_contact" "example"{
	custrole = "view"
	name = "Example"
	addresses {
		address = "example@expl.com"
		type = "email"
	}
	addresses {
		address = "expl.com"
		type = "webhook"
		action = "PUT"
		data = {"grr": "rrr"}
	}
}
```

The provider closely follows [NodePing API](https://nodeping.com/docs-api-contacts.html), so all parameters like `suppressup`, `suppressfirst` and additional parameters for the notification types are available.

There is also a NodePing contact data source.
```
data "nodeping_contact" "my_contact" {
	id = "202103031206A0A0A-A0A0A"
}

output "contact" {
  value = data.nodeping_contact.my_contact
}
```

## Development

Just like in case of normal use, `NODEPING_API_TOKEN` environmental variable needs to be set.

It will be usefull, from a developers stand point, to also set `TF_LOG=DEBUG`. More info (here)[https://www.terraform.io/docs/internals/debugging.html].

This project includes a Makefile to ease standard every day tasks. Currently this includes three commands:
- `make build` - builds the package,
- `make install` - builds the package, and moves it to terraform plugins,
- `make run_tests` - installs the provider, and runs tests.

For `make install` (and by extension `make run_tests`) to work, an `OS_ARCH` environment variable should be set. If it's not present, then "linux_amd64" is assumed.

Note that terraform keeps a checksum of providers in projects state, so after every plugin re-installation terraform state needs to be reset.

At this point, the only tests in the project are integration tests. Note that if tests fail, you may be left with some resources created at NodePing. 
