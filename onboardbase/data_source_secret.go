package onboardbase

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"

	// "os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientStruct := m.(APIClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	query := fmt.Sprintf(`query {
		generalPublicProjects(filterOptions: { title: "test", disableCustomSelect: true }) {
		  list {
			id
			title
			value
			publicEnvironments(filterOptions: { title: "development" }) {
			  list {
				id
				key
				title
			  }
			}
		  }
		}
	  }
	`)

	host := clientStruct.host
	client := clientStruct.client
	apiKey := clientStruct.apiKey
	passcode := clientStruct.passcode
	name := d.Get("name").(string)

	reqBody := map[string]interface{}{}
	reqBody["query"] = query
	graphqlQuery, err := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/graphql", host), strings.NewReader(string(graphqlQuery)))
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return diag.FromErr(err)
	}
	if result["errors"] != nil {
		error := result["errors"].([]interface{})[0].(map[string]interface{})
		extensions := error["extensions"].(map[string]interface{})
		exceptions := extensions["exception"].(map[string]interface{})
		return diag.FromErr(errors.New(exceptions["message"].(string)))
	}
	resultData := result["data"].(map[string]interface{})

	encoded, err := Parseresult(resultData)
	if err != nil {
		return diag.FromErr(err)
	}
	secret, err := DecryptSecrets(encoded, passcode, name, client)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("secret", secret)
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
