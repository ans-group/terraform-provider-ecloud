package ecloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
)

func dataSourceAPIRequestFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"property": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func buildDataSourceAPIRequestFilters(set *schema.Set) []connection.APIRequestFiltering {
	var filters []connection.APIRequestFiltering
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		filters = append(filters, *connection.NewAPIRequestFiltering(m["property"].(string), connection.EQOperator, filterValues))
	}
	return filters
}
