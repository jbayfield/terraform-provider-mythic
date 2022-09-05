package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/jbayfield/mythic-client-go"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"keyid": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MYTHIC_API_KEYID", nil),
				},
				"secret": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("MYTHIC_API_SECRET", nil),
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"mythic_vpsproducts": dataSourceVPSProducts(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"mythic_vps": resourceVPS(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the providerd name for yours):
		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent

		keyid := d.Get("keyid").(string)
		secret := d.Get("secret").(string)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if (keyid != "") && (secret != "") {
			c, err := mythic.NewClient(nil, &keyid, &secret)
			if err != nil {
				return nil, diag.FromErr(err)
			}

			return c, diags
		} else {
			return nil, diag.Errorf("specify keyid and secret")
		}
	}
}
