package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/jbayfield/mythic-client-go"
)

func dataSourceVPSProducts() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample data source in the Terraform provider scaffolding.",

		ReadContext: dataSourceVPSProductsRead,
		Schema: map[string]*schema.Schema{
			"vps_products": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"cores": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVPSProductsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*mythic.Client)

	var diags diag.Diagnostics

	products, err := c.GetVPSProducts(&c.Token)
	if err != nil {
		return diag.FromErr(err)
	}

	productsSlice := []mythic.VPSProduct{}
	for _, value := range products {
		productsSlice = append(productsSlice, value)
	}

	productsFlattened := flattenProductsData(&productsSlice)
	tflog.Info(ctx, fmt.Sprintf("%v", productsFlattened))
	if err := d.Set("vps_products", productsFlattened); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenProductsData(products *[]mythic.VPSProduct) []interface{} {
	if products != nil {
		ps := make([]interface{}, len(*products))

		for i, product := range *products {
			pi := make(map[string]interface{})

			pi["name"] = product.Name
			pi["description"] = product.Description
			pi["cores"] = product.Specs.Cores
			pi["ram"] = product.Specs.RAM
			pi["bandwidth"] = product.Specs.Bandwidth

			ps[i] = pi
		}

		return ps
	}

	return make([]interface{}, 0)
}
