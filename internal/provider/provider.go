package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const API_HOST = "https://api.onboardbase.com"

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"apikey": {
					Type:     schema.TypeString,
					Required: true,
				},
				"passcode": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"onboardbase_resource": resourceScaffolding(),
			},
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

type ApiClient struct {
	userAgent string
	host      string
	client    *http.Client
	apiKey    string
	passcode  string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userAgent := p.UserAgent("terraform-provider-scaffolding", version)

		apiKey := d.Get("apikey").(string)
		passcode := d.Get("passcode").(string)
		var diags diag.Diagnostics
		c := &http.Client{}
		return &ApiClient{
			userAgent: userAgent,
			host:      API_HOST,
			passcode:  passcode,
			apiKey:    apiKey,
			client:    c,
		}, diags
	}
}
