package onboardbase

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type APIClient struct {
	host     string
	client   *http.Client
	apiKey   string
	passcode string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("apikey").(string)
	passcode := d.Get("passcode").(string)
	host := "https://api.onboardbase.com"

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	c := &http.Client{}
	clientStruct := APIClient{
		client:   c,
		host:     host,
		apiKey:   apiKey,
		passcode: passcode,
	}
	return clientStruct, diags
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// DefaultFunc: schema.EnvDefaultFunc("ONBOARDBASE_APIKEY", nil),
			},
			"passcode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"onboardbase_secret": dataSourceSecret(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}
