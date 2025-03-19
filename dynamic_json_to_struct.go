package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return &schema.Provider{
				DataSourcesMap: map[string]*schema.Resource{
					"dynamic_todo_list": dataSourceTodoList(),
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
						"data": {
							Type:     schema.TypeMap, // Allow any dynamic key-value pairs from the JSON response
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

	// Convert dynamic map data to struct
	var todoStructs []interface{}
	for _, todo := range todos {
		todoStruct, err := mapToStruct(todo)
		if err != nil {
			return diag.FromErr(err)
		}
		todoStructs = append(todoStructs, todoStruct)
	}

	// Set the struct data as output in Terraform
	if err := d.Set("todos", todoStructs); err != nil {
		return diag.FromErr(err)
	}

	// Set a unique ID for the resource
	d.SetId(time.Now().Format(time.RFC3339))

	return nil
}

// Convert a dynamic map into a struct using reflection
func mapToStruct(data map[string]interface{}) (interface{}, error) {
	// Create a dynamic struct type
	var structType reflect.Type
	var structFields []reflect.StructField

	// Dynamically create fields for the struct
	for key, value := range data {
		field := reflect.StructField{
			Name: fmt.Sprintf("%s", key), // Field name from map key
			Type: reflect.TypeOf(value),  // Field type from map value type
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, key)),
		}
		structFields = append(structFields, field)
	}

	// Construct the struct type
	structType = reflect.StructOf(structFields)

	// Create an instance of the struct
	structPtr := reflect.New(structType).Elem()

	// Populate the struct with the data
	for key, value := range data {
		field := structPtr.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}
	}

	// Return the struct pointer as an interface{}
	return structPtr.Interface(), nil
}
