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

To authenticate to NodePing API, one needs to use API token. It can be passed to provider plugin in two ways. The preferred approach is to define a token in a `provider` block in your Terraform file:

```hcl
provider "nodeping" {
  token = "00AAAAAA-A0A0-AAAA-0000-A0AA0A000AAA"
}
```

As alternative you can specify token using an environment variable:
```bash
export NODEPING_API_TOKEN=00AAAAAA-A0A0-AAAA-0000-A0AA0A000AAA
terraform apply
```

### Resources

This provider supports the following type of resources:
- Check
- Contact
- Schedule 

#### Check

This is an simple example declaration of `check` resource:

```hcl
resource "nodeping_check" "some_check"{
	label = "SomeCheck"
	type = "HTTP"
	target = "http://example.eu/"
	enabled = "inactive"
}
```

The implementation generally follows [API documentation](https://nodeping.com/docs-api-checks.html), with some differences.

 - `customerid` customerid of the subaccount the check belongs to.
 - `type` string check type, one of: AGENT, AUDIO, CLUSTER, DOHDOT, DNS, FTP, HTTP, HTTPCONTENT, HTTPPARSE, HTTPADV, IMAP4, MYSQL, NTP, PING, POP3, PORT, PUSH, RBL, RDP, SIP, SMTP, SNMP, SPEC10DNS, SPEC10RDDS, SSH, SSL, WEBSOCKET, WHOIS
 - `target` check target. Note that for HTTP, HTTPCONTENT, HTTPPARSE, and HTTPADV checks this must begin with "http://" or "https://".
 - `label` give this check a label. If this is absent, the target will be used as the label.
 - `interval` how often this check runs in minutes. Can be any integer starting at 1 for one minute checks. Once a day is 1440. Defaults to 15.
 - `enabled` set to 'active' to enable this check to run, set to 'inactive' to disable.
> `enabled` differs from API in that it only accepts strings: 'active', 'inactive', and it defaults to 'active'.
 - `public` set to `true` to enable public reports for this check.
> `public` differs from API in that it uses boolean values, and not strings.
 - `runlocations` an array of geographical regions where probe servers are located. This can be an one or more of the following: 'nam' for North America, 'eur' for Europe, 'eao' for East Asia/Oceania, or 'wlw' for worldwide. If omitted, the account's default location setting is used. To run checks on an AGENT, set the location to the ID of the agent check you want to run on.
> Running checks on an AGENT was never tested

> `runlocations` differs from API in that it only allows array as an input value.
 - `homeloc` set the preferred probe location, home location, for this check. The default is 'false' and will run the check on a random probe in the selected region (see runlocations) or the account default region if no region is specified on the check. The probe two letter indicators are listed in our FAQ (example: 'ca' would run a check from our California probe). Set this value to 'roam' to have the check change probe location on each interval.
> `homeloc` differs from API in that it will only accept strings. If you need to run the check on a random probe, set this parameter to `"false"` (string).
 - `threshold` the timeout for this check in seconds, defaults to 5 for a five second timeout. Can be any integer starting at 1. For CLUSTER checks, this indicates how many checks listed in the 'data' element must be passing in order for this check to pass.
 - `sens` number of rechecks before this check is considered 'down' or 'up'. Defaults to 2 (2 rechecks means it will be checked by 3 different servers). Rechecks are run immediately, not based on the interval and happen when a status change occurs.
 - `notifications` array of objects containing the contact id, delay, and scheduling for notifications. The IDs can be obtained by listing the relevant contacts. See more robust example below for usage. 
> Unlike in case of API, there is no need to set the contact key value to "None", in order to remove an address from a check's notifications. Simply removing the address item is enough.
 - `dep` - optional string - the id of the check used for the notification dependency. If the check this is set to is failing, no notifications will be sent. For example, set this to the check id of a PING check on an edge router for all services that depend on that router. It helps reduce the number of alerts you receive when core networks or services go offline. To remove this functionality, set this to false.
 - `description` - optional string - you can put arbitrary text, JSON, XML, etc. Size limit is 1000 characters.

The following are only relevant for certain types:

 - `checktoken` read-only field on AGENT and PUSH checks - can request server-side re-generation by setting this field to 'reset'.
 - `clientcert` string to specify the ID of a client certificate/key to be used in the DOHDOT check.
 - `contentstring` string for DOHDOT, DNS, HTTPCONTENT, HTTPADV, FTP, SSH, WEBSOCKET, WHOIS type checks - the string to match the response against.
 - `dohdot` string used to specify DoH or DoT in the DOHDOT check. Valid value is either 'doh' or 'dot' - defaults to 'doh'.
 - `dnstype` string for DNS and DOHDOT checks to indicate the type of DNS entry to query - String set to one of: 'ANY', 'A', 'AAAA', 'CNAME', 'MX, 'NS, 'PTR', 'SOA', 'SRV', 'TXT'.
 - `dnstoresolve` optional string for DNS and DOHDOT checks - The FQDN of the DNS query
 - `dnsrd` boolean for DNS RD (Recursion Desired) bit - defaults to true. If you're using CloudFlare DNS servers, set this to false.
 - `transport` string for DNS check transport protocol - defaults to 'udp' but 'tcp' is also supported.
 - `follow` boolean used for HTTP, HTTPCONTENT and HTTPADV checks. If true, the check will follow up to four redirects. The HTTPADV check only supports follow for GET requests.
 - `email` string used for IMAP4 and SMTP checks.
 - `port` positive integer for DNS, FTP, NTP, PORT, SSH type checks - used for check types that support port fields separate from the target address. HTTP and HTTPCONTENT will ignore this field as the port must be set in the target in standard URL format. This field is required by PORT and NTP checks.
 - `username` string for FTP, IMAP4, POP3, SMTP and SSH type checks - HTTP and HTTPCONTENT will ignore this field as the username must be set in the target in standard URL format.
 - `password` string for FTP, IMAP4, POP3, SMTP and SSH type checks - Note that this is going to be passed back and forth in the data, so you should always be sure that credentials used for checks are very limited in their access level. See our Terms of Service. HTTP and HTTPCONTENT will ignore this field as the password must be set in the target in standard URL format.
 - `secure` string to specify whether the IMAP4, SMTP, and POP3 checks should use TLS for the check. Can be set to "false" or "ssl".
 - `verify` string to set whether or not to verify the SSL certificate (SMTP, IMAP4, POP3, DOHDOT check types) or DNSSEC (DNS check type only). Can be "true" or "false".
 - `ignore` string for the RBL check type, specifies RBL lists to ignore. Multiple lists can be added by including them in the string, separated by commas.
 - `invert` string for FTP, HTTPCONTENT, HTTPADV, NTP, PORT, SSH type checks - used for 'Does not contain' functionality in checks. Default is 'false' - Set to 'true' to invert the content type match.
 - `warningdays` positive integer for SSL, WHOIS, POP3, IMAP4, and SMTP checks - number of days before certificate (or domain for WHOIS) expiration to fail the check and send a notification.
 - `fields` used for fields to parse from the HTTPADV, HTTPPARSE, and SNMP response. Each block should include elements: name, min and max.
 - `postdata` string that can be used in the HTTPADV check as an alternative to the data object. postdata should be a single string to post.
 - `data, receiveheaders, sendheaders` these are optional maps used by HTTPADV ('data' can also be used for CLUSTER checks - see blow - 'sendheaders' can be used for DOHDOT and HTTPPARSE). They are formatted as key:value pairs.
 - `edns` object used to send EDNS(0) OPT psudo-records in a DNS query in the DOHDOT check type. Must be formatted as key:value pairs
 - `method` string used by the HTTPADV and DOHDOT checks to specify the HTTP method. For the HTTPADV check, value can be one of: GET, POST, PUT, HEAD, TRACE, or CONNECT. For DOHDOT, value can be GET or POST.
 - `statuscode` positive integer specifying the expected HTTP status code in the response to an HTTPADV or DOHDOT check.
 - `ipv6` boolean specifying the check should run against an ipv6 address - PING, HTTP, HTTPCONTENT, HTTPADV, WHOIS, and DOHDOT checks.
 - `regex` boolean/string to set whether the 'contentstring' element is a regular expression or just a string to be matched. Can be "true"/true/1 or "false"/false/0. HTTPCONTENT and HTTPADV only.
 - `servername` string to specify the FQDN sent to SNI services in the SSL check.
 - `snmpv` string specifying the SNMP version the check should use. Valid values are "1" and "2c". Defaults to "1" - SNMP check only.
 - `snmpcom` string specifying the SNMP community indicator that should be used. Defaults to 'public' - SNMP check only.
 - `verifyvolume` boolean to enable the volume detection feature - AUDIO check only.
 - `volumemin` integer (acceptable range -90 to 0) used by the volume detection feature - AUDIO check only.
 - `whoisserver` string specifying the WHOIS server FQDN or IPv(4/6) to query - WHOIS check only.


#### Contact

This is an example declaration of `contact` resource:
```hcl
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
        headers = {"Content-Type": "application/json"}
        querystrings = {"querykey1":"value1"}
	}
}
```

The provider closely follows [NodePing API](https://nodeping.com/docs-api-contacts.html), so all parameters like `suppressup`, `suppressfirst` and additional parameters for the notification types are available.

 - `customerid` customerid of the subaccount to which the contact belongs. Not needed if the contact is in your primary account.
 - `name` the name of the contact, used as a label.
 - `custrole` set to 'edit,' 'view' or 'notify' to set permissions for this contact. Default is 'view'. Contacts created with an 'edit' or 'view' access level will receive a welcome email message, which by default includes a random password and the suggestion that they change the password immediately. You can avoid these messages by creating the contact with a 'notify' custrole, and then updating them to 'edit' or 'view'. You can also change the content of the welcome email message in the Branding preferences.
 - `addresses` optional collection of email addresses and phone numbers. See example above.
> Significant difference to the original API is that this terraform provider is aware of the order in with addresses are declared. This is due to the use of address id when connecting checks and addresses. See more robust example below.
 - `priority` additional parameter for the pushover notification type. Defaults to 2, with ther valid values are 1, 0, -1, or -2 per the priority API documentation for Pushover.

Additional parameters for the webhook notification type:

 - `action` defaults to 'get' but can be 'put', 'post', 'head', or 'delete'
 - `data` payload or body of an HTTP POST or PUT request. This can be JSON, XML, or any arbitrary string.
 - `headers` set HTTP request headers
 - `querystrings` set HTTP query string key/values, that will be appended to your webhook URL as part of the query string.

#### Schedule

This is an example declaration of a `schedule` resource:

```hcl
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

 - `customerid` - customerid of the subaccount to which the schedule belongs.
 - `name` - this parameter is not present in official documentation. In practice the API requires an `id` parameter when creating a new schedule, but then returns a response like this: `{"ok":true,"id":"100000000000A0A0A"}`, that indicates there is some other id. To avoid this confusion, schedule name is used as it's id.
 - `data` - required object containing the schedule (see the example above). Properties inside each day object are as follows:
    - `time1` - start of time span
    - `time2` - end of time span
    - `exclude` - inverts the time span so it is all day except for the time between time1 and time2
    - `disabled` - disables notifications for this day.
    - `allday` - enables notifications for the entire day.

### More robust example

This example shows how to use all these resources together.

```hcl
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

### Data sources

Currently three data sources are implemented: `nodeping_contact`, `nodeping_check`, `nodeping_schedule`.

#### Contact

Here's an example use:

```hcl
data "nodeping_contact" "my_contact" {
	id = "202103031206A0A0A-A0A0A"
}

output "contact" {
  value = data.nodeping_contact.my_contact
}
```

#### Check

```hcl
data "nodeping_check" "the_check" {
	id = var.check_id
}

output "check_id" {
	value = data.nodeping_check.the_check.id
}
```

#### Schedule

Schedule data source uses `name` attribute instead of `id` when running queries.

```hcl
data "nodeping_schedule" "the_schedule" {
	name = "Weekdays"
}

output "schedule_name" {
	value = data.nodeping_schedule.the_schedule.name
}

output "schedule" {
	value = data.nodeping_schedule.the_schedule
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

## License
 
This project is distributed under [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0.html).