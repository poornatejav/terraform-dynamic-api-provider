package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return &schema.Provider{
				DataSourcesMap: map[string]*schema.Resource{
					"todo_list": dataSourceTodoList(),
				},
			}
		},
	})
}

func dataSourceTodoList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTodoListRead,
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTodoListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", "http://localhost:9000/todo", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("failed to fetch todos: %s", resp.Status)
	}

	var jsonData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("data", jsonData); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))

	return nil
}
