package main

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
	"time"

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
			"todos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"completed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTodoListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*http.Client)

	// Call your REST API
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

	var todos []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return diag.FromErr(err)
	}

	// Store in Terraform state
	if err := d.Set("todos", todos); err != nil {
		return diag.FromErr(err)
	}

	// Unique ID for Terraform state
	d.SetId(time.Now().Format(time.RFC3339))

	return nil
}
