package onboardbase

import (
	"context"
	// "encoding/json"
	"fmt"
	// "net/http"
	// "os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// host := os.Getenv("API_HOST")
	// req, err := http.NewRequest("GET", fmt.Sprintf("%s/Secret", host), nil)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// r, err := client.Do(req)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	// defer r.Body.Close()
	type response struct {
		Name string
	}

	// err = json.NewDecoder(r.Body).Decode(&Secret)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	name := d.Get("name")
	if err := d.Set("secret", fmt.Sprintf("Hello, %v", name)); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"secret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
