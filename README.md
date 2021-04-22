# Terraform NodePing Provider

This is [NodePing](https://nodeping.com) provider for Terraform. It can be used to manage 
[NodePing resources](https://nodeping.com/docs-api-overview.html) like Checks or Contacts.

## Quick start

First of all, build and install the provider into your local Terraform installation:

```bash
git clone https://github.com/softkraftco/terraform-nodeping.git
cd terraform-nodeping
make install
```

Next create new Terraform project:

```bash
mkdir mynodeping
cd mynodeping
touch mynodeping.tf 
```

As a next step add NodePing provider to `required_providers` section of `mynodeping.tf ` file:

```hcl
terraform {
  required_providers {
    nodeping = {
      version = "0.0.1"
      source  = "softkraft.co/terraform/nodeping"
    }
  }  
}

provider "nodeping" {
    token = "00AAAAAA-A0A0-AAAA-0000-A0AA0A000AAA"
}
```

Then add NodePing check resource definition to `mynodeping.tf ` file:

```hcl
resource "nodeping_check" "mycheck"{
	label = "mycheck"
	type = "HTTP"
	target = "http://example.eu/"
}
```

And finally execute Terraform!

```bash
terraform apply
```


## Usage

To authenticate to NodePing API, one needs to use API token. It can be passed to provider plugin in two ways: 
set an environment variable:
```
export NODEPING_API_TOKEN=00AAAAAA-A0A0-AAAA-0000-A0AA0A000AAA
```

or define it in a `provider` block in your terraform file:

```
provider "nodeping" {
  token = "00AAAAAA-A0A0-AAAA-0000-A0AA0A000AAA"
}
```

### Resources

#### Checks

This is an simple example declaration of `check` resource:
```
resource "nodeping_check" "some_check"{
	label = "SomeCheck"
	type = "HTTP"
	target = "http://example.eu/"
	enabled = "inactive"
}
```

As of now, only HTTP, SSH, and SSL checks were tested, but parameters used by other check types are implemented and should work, although without any guaranties.

The implementation generally follows [API documentation](https://nodeping.com/docs-api-checks.html), with some differencess:
 - `enabled` parameter accepts only strings: "active", "inactive"
 - `homeloc` will only accept strings, if you need to run the check on a random probe, set this tarameter to `"false"` (string).
 - `public` is a bool, not a string.

Remember that by default nodeping creates checks as inactive.

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

One difference to the original API is that this terraform provider is aware of the order in with addresses are declared. This is due to the use of address id when connecting checks and addresses. See more robust example below.

#### Schedules

This is an ecample declaration of a `schedule` resource:

```
resource "nodeping_schedule" "my_schedule"{
	name = "MySchedule"
	data {
		day = "monday"
		time1 = "16:00"
		time2 = "17:00"
		exclude = false
	}
	data {
		day = "sunday"
		allday = true
	}
}

output "my_schedule_name" {
	value = nodeping_schedule.my_schedule.name
}
```

Note that the `data` parameter is declared differently then in the official [API documentation](https://nodeping.com/docs-api-schedules.html).

Note that this resource uses a `name` parameter, that is not present in official documentation. In practice the API requires an `id` parameter when creating a new schedule, but then returns a response like this: `{"ok":true,"id":"100000000000A0A0A"}`, that indicates there is some other id. To avoid this confusion, schedule name is used as it's id.

### More robust example

This example shows how to use all these resources together.

```
terraform {
  required_providers {
    nodeping = {
      version = "0.0.1"
      source  = "softkraft.co/terraform/nodeping"
    }
  }
}

resource "nodeping_contact" "my_contact"{
	custrole = "owner"
	name = "MyContact"
	addresses {
		address = "my-email@o1.com"
		type = "email"
	}
	addresses {
		address = "some-other-email@o1.com"
		type = "email"
	}
}

resource "nodeping_schedule" "my_schedule"{
	name = "MySchedule"
	data {
		day = "monday"
		time1 = "16:00"
		time2 = "17:00"
		exclude = false
	}
	data {
		day = "sunday"
		allday = true
	}
}

resource "nodeping_check" "my_check"{
	label = "MyCheck"
	type = "HTTP"
	target = "http://example.eu/"
	enabled = "inactive"
	notifications {
		contact = nodeping_contact.my_contact.addresses[0].id
		delay = 1
		schedule = "Weekdays"
	}
    notifications {
		contact = nodeping_contact.my_contact.addresses[1].id
		delay = 10
		schedule = nodeping_schedule.my_schedule.name
	}
	homeloc = "false"
	follow = true
	ipv6 = true
}


output "my_check_id" {
	value = nodeping_check.my_check.id
}

output "my_contact_id" {
	value = nodeping_contact.my_contact.id
}

output "my_first_address_id"{
	value = nodeping_contact.my_contact.addresses[0].id
}

output "my_schedule_name" {
	value = nodeping_schedule.my_schedule.name
}
```

### Data sourcesss

Currently the only implemented data source is contact. Here's an example use:

```
data "nodeping_contact" "my_contact" {
	id = "202103031206A0A0A-A0A0A"
}

output "contact" {
  value = data.nodeping_contact.my_contact
}
```

## Development

It will be useful, from a developers stand point, to set `TF_LOG=DEBUG`. More info [here](https://www.terraform.io/docs/internals/debugging.html).

This project includes a Makefile to ease standard every day tasks. Currently this includes three commands:
- `make build` - builds the package,
- `make install` - builds the package, and moves it to terraform plugins,
- `make run_tests` - installs the provider, and runs tests.

For `make install` (and by extension `make run_tests`) to work, an `OS_ARCH` environment variable should be set. If it's not present, then "linux_amd64" is assumed.

The `make run_tests` command requires `NODEPING_API_TOKEN` environmental variable needs to be set, even if it is declared in terraform files. This is because tests query nodeping API to check that expected resources were created. 

Terraform keeps a checksum of providers in projects state, so after every plugin re-installation terraform state needs to be reset (removing .terraform, .terraform.lock.hcl, terraform.tfstate).

At this point, the only tests in the project are integration tests. Note that if a tests fails with an error, you may be left with some resources created at NodePing. In that case there might also be a need to restart terraform state.
