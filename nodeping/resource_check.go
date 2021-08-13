package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"terraform-nodeping/nodeping_api_client"
)

func resourceCheck() *schema.Resource {

	// prepare accepted values for validation
	checkTypes := []string{
		"AGENT", "AUDIO", "CLUSTER", "DOHDOT", "DNS", "FTP", "HTTP",
		"HTTPCONTENT", "HTTPPARSE", "HTTPADV", "IMAP4", "MYSQL", "NTP", "PING",
		"POP3", "PORT", "PUSH", "RBL", "RDP", "SIP", "SMTP", "SNMP", "SPEC10DNS",
		"SPEC10RDDS", "SSH", "SSL", "WEBSOCKET", "WHOIS",
	}
	httpAdvMethods := []string{"GET", "POST", "PUT", "HEAD", "TRACE", "CONNECT"}
	trueFalseStrings := []string{"false", "true"}

	return &schema.Resource{
		CreateContext: resourceCheckCreate,
		ReadContext:   resourceCheckRead,
		UpdateContext: resourceCheckUpdate,
		DeleteContext: resourceCheckDelete,
		Schema: map[string]*schema.Schema{
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true, Optional: true},
			"type": &schema.Schema{Type: schema.TypeString, Required: true,
				ValidateFunc: validation.StringInSlice(checkTypes, false)},
			"target":       &schema.Schema{Type: schema.TypeString, Optional: true},
			"label":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"interval":     &schema.Schema{Type: schema.TypeInt, Optional: true, Default: 15},
			"enabled":      &schema.Schema{Type: schema.TypeString, Optional: true, ValidateFunc: validation.StringInSlice([]string{"active", "inactive"}, false), Default: "active"},
			"public":       &schema.Schema{Type: schema.TypeBool, Optional: true, Default: false},
			"runlocations": &schema.Schema{Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"homeloc":      &schema.Schema{Type: schema.TypeString, Optional: true},
			"threshold":    &schema.Schema{Type: schema.TypeInt, Optional: true, Default: 5},
			"sens":         &schema.Schema{Type: schema.TypeInt, Optional: true},
			"notifications": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contact":  &schema.Schema{Type: schema.TypeString, Required: true},
						"delay":    &schema.Schema{Type: schema.TypeInt, Required: true},
						"schedule": &schema.Schema{Type: schema.TypeString, Required: true},
					},
				},
			},
			"dep":         &schema.Schema{Type: schema.TypeString, Optional: true},
			"description": &schema.Schema{Type: schema.TypeString, Optional: true},
			// the following are only relevant for certain types:
			"checktoken":    &schema.Schema{Type: schema.TypeString, Computed: true},
			"clientcert":    &schema.Schema{Type: schema.TypeString, Optional: true},
			"contentstring": &schema.Schema{Type: schema.TypeString, Optional: true},
			"dohdot": &schema.Schema{Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"doh", "dot"}, false)},
			"dnstype": &schema.Schema{Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ANY", "A", "AAAA", "CNAME", "MX", "NS", "PTR", "SOA", "SRV", "TXT",
				}, false)},
			"dnstoresolve": &schema.Schema{Type: schema.TypeString, Optional: true},
			"dnsrd":        &schema.Schema{Type: schema.TypeBool, Optional: true},
			"transport": &schema.Schema{Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"udp", "tcp"}, false)},
			"follow":   &schema.Schema{Type: schema.TypeBool, Optional: true},
			"email":    &schema.Schema{Type: schema.TypeString, Optional: true},
			"port":     &schema.Schema{Type: schema.TypeInt, Optional: true},
			"username": &schema.Schema{Type: schema.TypeString, Optional: true},
			"password": &schema.Schema{Type: schema.TypeString, Optional: true},
			"secure": &schema.Schema{Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"false", "ssl"}, false)},
			"verify": &schema.Schema{Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringInSlice(trueFalseStrings, false)},
			"ignore": &schema.Schema{Type: schema.TypeString, Optional: true},
			"invert": &schema.Schema{Type: schema.TypeBool, Optional: true, Default: false},
			"warningdays": &schema.Schema{Type: schema.TypeInt, Optional: true,
				ValidateFunc: validation.IntAtLeast(0)},
			"fields": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   &schema.Schema{Type: schema.TypeString, Optional: true},
						"name": &schema.Schema{Type: schema.TypeString, Required: true},
						"min":  &schema.Schema{Type: schema.TypeInt, Required: true},
						"max":  &schema.Schema{Type: schema.TypeInt, Required: true},
					},
				},
			},
			"postdata":       &schema.Schema{Type: schema.TypeString, Optional: true},
			"data":           &schema.Schema{Type: schema.TypeMap, Optional: true},
			"receiveheaders": &schema.Schema{Type: schema.TypeMap, Optional: true},
			"sendheaders":    &schema.Schema{Type: schema.TypeMap, Optional: true},
			"edns":           &schema.Schema{Type: schema.TypeMap, Optional: true},
			"method": &schema.Schema{Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringInSlice(httpAdvMethods, false)},
			"statuscode": &schema.Schema{Type: schema.TypeInt, Optional: true,
				ValidateFunc: validation.IntAtLeast(0)},
			"ipv6":       &schema.Schema{Type: schema.TypeBool, Optional: true, Default: false},
			"regex":      &schema.Schema{Type: schema.TypeBool, Optional: true},
			"servername": &schema.Schema{Type: schema.TypeString, Optional: true},
			"snmpv":      &schema.Schema{Type: schema.TypeString, Optional: true},
			"snmpcom": &schema.Schema{Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"1", "2c"}, false)},
			"verifyvolume": &schema.Schema{Type: schema.TypeBool, Optional: true},
			"volumemin": &schema.Schema{Type: schema.TypeInt, Optional: true,
				ValidateFunc: validation.IntBetween(-90, 0)},
			"whoisserver": &schema.Schema{Type: schema.TypeString, Optional: true},
		},
	}
}

func applyCheckToSchema(check *nodeping_api_client.Check, d *schema.ResourceData) error {
	d.SetId(check.ID)

	err := d.Set("customer_id", check.CustomerId)
	if err != nil {
		return err
	}
	err = d.Set("type", check.Type)
	if err != nil {
		return err
	}
	err = d.Set("label", check.Label)
	if err != nil {
		return err
	}
	err = d.Set("enabled", check.Enable)
	if err != nil {
		return err
	}
	err = d.Set("public", check.Public)
	if err != nil {
		return err
	}
	err = d.Set("runlocations", check.Runlocations)
	if err != nil {
		return err
	}
	if check.HomeLoc == false {
		err = d.Set("homeloc", "false")
	} else {
		err = d.Set("homeloc", check.HomeLoc)
	}
	if err != nil {
		return err
	}
	err = d.Set("notifications", flattenNotifications(&check.Notifications))
	if err != nil {
		return err
	}
	err = d.Set("interval", check.Interval)
	if err != nil {
		return err
	}

	err = d.Set("description", check.Description)
	if err != nil {
		return err
	}

	for key, val := range check.Parameters {
		err = d.Set(key, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func flattenNotifications(notifications *[]map[string]nodeping_api_client.Notification) []map[string]interface{} {
	if notifications == nil { // return fast if nothing to do
		return make([]map[string]interface{}, 0)
	}

	list := make([]map[string]interface{}, len(*notifications))
	for idx, notificationMap := range *notifications {
		for notificationId, notification := range notificationMap {
			flattened := make(map[string]interface{})
			flattened["contact"] = notificationId
			flattened["delay"] = notification.Delay
			flattened["schedule"] = notification.Schedule

			list[idx] = flattened
		}
	}

	return list
}

func flattenFields(fields *map[string]nodeping_api_client.CheckField) []interface{} {
	if fields == nil { // return fast if nothing to do
		return make([]interface{}, 0)
	}

	flattened := make([]interface{}, len(*fields), len(*fields))
	for fieldId, field := range *fields {
		flattened := make(map[string]interface{})
		flattened["id"] = fieldId
		flattened["name"] = field.Name
		flattened["min"] = field.Min
		flattened["max"] = field.Max
	}

	return flattened
}

func getCheckUpdateFromSchema(d *schema.ResourceData) *nodeping_api_client.CheckUpdate {
	var checkUpdate nodeping_api_client.CheckUpdate

	checkUpdate.ID = d.Id()
	checkUpdate.CustomerId = d.Get("customer_id").(string)
	checkUpdate.Type = d.Get("type").(string)
	checkUpdate.Target = d.Get("target").(string)
	checkUpdate.Label = d.Get("label").(string)
	checkUpdate.Interval = d.Get("interval").(int)
	checkUpdate.Enable = d.Get("enabled").(string)
	// this silly conversion is here because API expects public to be a string,
	// but returns a bool
	if d.Get("public").(bool) {
		checkUpdate.Public = "true"
	} else {
		checkUpdate.Public = "false"
	}
	for _, runLocation := range d.Get("runlocations").(*schema.Set).List() {
		checkUpdate.RunLocations = append(checkUpdate.RunLocations, runLocation.(string))
	}
	if d.Get("homeloc").(string) == "false" {
		checkUpdate.HomeLoc = false
	} else {
		checkUpdate.HomeLoc = d.Get("homeloc").(string)
	}
	checkUpdate.Threshold = d.Get("threshold").(int)
	checkUpdate.Sens = d.Get("sens").(int)

	notificationsSchemaList := d.Get("notifications").(*schema.Set).List()
	checkUpdate.Notifications = make([]map[string]nodeping_api_client.Notification, len(notificationsSchemaList))
	for idx, nS := range notificationsSchemaList {
		notisicationSchema := nS.(map[string]interface{})
		notificationMap := make(map[string]nodeping_api_client.Notification, 1)
		notificationMap[notisicationSchema["contact"].(string)] = nodeping_api_client.Notification{
			notisicationSchema["delay"].(int),
			notisicationSchema["schedule"].(string),
		}
		checkUpdate.Notifications[idx] = notificationMap
	}

	checkUpdate.Dep = d.Get("dep").(string)
	checkUpdate.Description = d.Get("description").(string)
	checkUpdate.CheckToken = d.Get("checktoken").(string)
	checkUpdate.ClientCert = d.Get("clientcert").(string)
	checkUpdate.ContentString = d.Get("contentstring").(string)
	checkUpdate.Dohdot = d.Get("dohdot").(string)
	checkUpdate.DnsType = d.Get("dnstype").(string)
	checkUpdate.DnsToResolve = d.Get("dnstoresolve").(string)
	checkUpdate.Dnsrd = d.Get("dnsrd").(bool)
	checkUpdate.Transport = d.Get("transport").(string)
	checkUpdate.Follow = d.Get("follow").(bool)
	checkUpdate.Email = d.Get("email").(string)
	checkUpdate.Port = d.Get("port").(int)
	checkUpdate.Username = d.Get("username").(string)
	checkUpdate.Password = d.Get("password").(string)
	checkUpdate.Secure = d.Get("secure").(string)
	checkUpdate.Verify = d.Get("verify").(string)
	checkUpdate.Ignore = d.Get("ignore").(string)
	checkUpdate.Invert = d.Get("invert").(bool)
	checkUpdate.WarningDays = d.Get("warningdays").(int)

	fields := d.Get("fields").(*schema.Set).List()
	for _, field := range fields {
		f := field.(map[string]interface{})
		checkUpdate.Fields[f["id"].(string)] = nodeping_api_client.CheckField{
			f["name"].(string), f["min"].(int), f["max"].(int),
		}
		checkUpdate.Fields = d.Get("fields").(map[string]nodeping_api_client.CheckField)
	}

	checkUpdate.Postdata = d.Get("postdata").(string)
	if d.Get("data") != nil {
		checkUpdate.Data = d.Get("data").(map[string]interface{})
	}
	if d.Get("receiveheaders") != nil {
		checkUpdate.Data = d.Get("receiveheaders").(map[string]interface{})
	}
	if d.Get("sendheaders") != nil {
		checkUpdate.Data = d.Get("sendheaders").(map[string]interface{})
	}
	if d.Get("ends") != nil {
		checkUpdate.Data = d.Get("ends").(map[string]interface{})
	}
	checkUpdate.Method = d.Get("method").(string)
	checkUpdate.Statuscode = d.Get("statuscode").(int)
	checkUpdate.Ipv6 = d.Get("ipv6").(bool)
	checkUpdate.Regex = d.Get("regex").(bool)
	checkUpdate.ServerName = d.Get("servername").(string)
	checkUpdate.Snmpv = d.Get("snmpv").(string)
	checkUpdate.Snmpcom = d.Get("snmpcom").(string)
	checkUpdate.VerifyVolume = d.Get("verifyvolume").(bool)
	checkUpdate.VolumeMin = d.Get("volumemin").(int)
	checkUpdate.WhoisServer = d.Get("whoisserver").(string)

	return &checkUpdate
}

func resourceCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	checkUpdate := getCheckUpdateFromSchema(d)

	savedCheck, err := client.CreateCheck(ctx, checkUpdate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(savedCheck.ID)
	d.Set("customer_id", savedCheck.CustomerId)
	return resourceCheckRead(ctx, d, m)
}

func resourceCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	checkId := d.Id()
	customerId := d.Get("customer_id").(string)

	check, err := client.GetCheck(ctx, customerId, checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = applyCheckToSchema(check, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	checkUpdate := getCheckUpdateFromSchema(d)

	_, err := client.UpdateCheck(ctx, checkUpdate)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCheckRead(ctx, d, m)
}

func resourceCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	checkId := d.Id()
	customerId := d.Get("customer_id").(string)
	err := client.DeleteCheck(ctx, customerId, checkId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("") // added here for explicitness

	return nil

}
