package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"terraform-provider-onboardbase/internal/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScaffolding() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   "Sample resource in the Terraform provider scaffolding.",
		ReadContext:   resourceScaffoldingRead,
		DeleteContext: resourceScaffoldingDelete,
		UpdateContext: resourceScaffoldingCreate,
		CreateContext: resourceScaffoldingCreate,
		Schema: map[string]*schema.Schema{
			"data": {
				Description: "The value of the resource",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"project": {
				Description: "The project to fetch secret from.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"keys": {
				Description: "The keys to return their values from the secret data",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"values": {
				Description: "The values of the keys",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"environment": {
				Description: "The environment from the project.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceScaffoldingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	project := d.Get("project").(string)
	keys := d.Get("keys").([]interface{})
	environment := d.Get("environment").(string)

	diags, encoded, err := fetchData(project, environment, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	return setValues(encoded, keys, diags, d, meta)
}

func resourceScaffoldingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	d.Set("data", "")
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceScaffoldingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	project := d.Get("project").(string)
	keys := d.Get("keys").([]interface{})
	environment := d.Get("environment").(string)

	diags, encoded, err := fetchData(project, environment, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	return setValues(encoded, keys, diags, d, meta)
}

func setValues(encoded string, keys []interface{}, diags diag.Diagnostics, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	var secrets []string
	json.Unmarshal([]byte(encoded), &secrets)

	clientStruct := meta.(*ApiClient)
	decryptedSecrets := utils.DecryptSecrets(secrets, clientStruct.passcode)

	var secretValues = make(map[string]string)

	for _, key := range keys {
		value := decryptedSecrets[key.(string)]
		if value == "" {
			return diag.FromErr(fmt.Errorf("%s is not a recognized secret", key))
		}
		secretValues[key.(string)] = value
	}

	d.Set("values", secretValues)
	return diags
}

func fetchData(project string, environment string, meta interface{}) (diag.Diagnostics, string, error) {
	var diags diag.Diagnostics

	clientStruct := meta.(*ApiClient)

	query := fmt.Sprintf(`query {
			generalPublicProjects(filterOptions: { title: "%v", disableCustomSelect: true }) {
			  list {
				id
				title
				value
				publicEnvironments(filterOptions: { title: "%v" }) {
				  list {
					id
					key
					title
				  }
				}
			  }
			}
		  }
		`, project, environment)

	host := clientStruct.host
	client := clientStruct.client
	apiKey := clientStruct.apiKey

	reqBody := map[string]interface{}{}
	reqBody["query"] = query
	graphqlQuery, err := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/graphql", host), strings.NewReader(string(graphqlQuery)))
	if err != nil {
		return diag.FromErr(err), "", nil
	}
	req.Header.Set("KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err), "", nil
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return diag.FromErr(err), "", nil
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return diag.FromErr(err), "", nil

	}
	if result["errors"] != nil {
		query_error := result["errors"].([]interface{})[0].(map[string]interface{})
		log.Println(query_error)
		extensions := query_error["extensions"].(map[string]interface{})
		exceptions := extensions["exception"].(map[string]interface{})
		return diag.FromErr(errors.New(exceptions["message"].(string))), "", nil
	}
	resultData := result["data"].(map[string]interface{})
	encoded, err := utils.Parseresult(resultData)
	return diags, string(encoded), err
}
