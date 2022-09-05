package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jbayfield/mythic-client-go"
)

func resourceVPS() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceVPSCreate,
		ReadContext:   resourceVPSRead,
		UpdateContext: resourceVPSUpdate,
		DeleteContext: resourceVPSDelete,

		// TODO: Work out which of the forcenew/computed/required are redundant
		Schema: map[string]*schema.Schema{
			"vps": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							ForceNew: false,
						},
						"identifier": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Computed: false,
							Required: true,
						},
						"product": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Computed: false,
							ForceNew: false,
						},
						"dormant": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: false,
							Optional: true,
							Default:  false,
							ForceNew: false,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostserver": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpumode": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Optional: true,
							Default:  "performance",
							ForceNew: false,
						},
						"netdevice": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Optional: true,
							Default:  "virtio",
							ForceNew: false,
						},
						"diskbus": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Optional: true,
							Default:  "virtio",
							ForceNew: false,
						},
						"price": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"isoimage": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Optional: true,
							Default:  "automated-install-config",
							ForceNew: false,
						},
						"bootdevice": &schema.Schema{
							Type:     schema.TypeString,
							Computed: false,
							Optional: true,
							Default:  "hd",
							ForceNew: false,
						},
						"ipv4": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ipv6": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"spec_disktype": &schema.Schema{
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Default:  "ssd",
						},
						"spec_disksize": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							Computed: false,
							ForceNew: false,
						},
						"spec_cores": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spec_ram": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"macs": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVPSCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	// idFromAPI := "my-id"
	// d.SetId(idFromAPI)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	// tflog.Trace(ctx, "created a resource")

	return diag.Errorf("not implemented")
}

func resourceVPSRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*mythic.Client)

	var diags diag.Diagnostics

	vps, err := c.GetVPS(d.Id(), &c.Token)
	if err != nil {
		return diag.FromErr(err)
	}

	vpsFlattened := flattenVPSData(vps)
	tflog.Info(ctx, fmt.Sprintf("%v", vpsFlattened))

	if err := d.Set("vps", vpsFlattened); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceVPSUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceVPSDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func flattenVPSData(vps *mythic.VPS) []interface{} {
	if vps != nil {
		ps := make([]interface{}, 0)
		info := make(map[string]interface{})

		info["name"] = vps.Name
		info["identifier"] = vps.Identifier
		info["product"] = vps.Product
		info["dormant"] = vps.Dormant
		info["status"] = vps.Status
		info["hostserver"] = vps.HostServer
		info["zone_code"] = vps.Zone.Code
		info["zone_name"] = vps.Zone.Name
		info["cpumode"] = vps.CPUMode
		info["netdevice"] = vps.NetDevice
		info["diskbus"] = vps.DiskBus
		info["price"] = vps.Price
		info["isoimage"] = vps.ISOImage
		info["bootdevice"] = vps.BootDevice
		info["ipv4"] = vps.IPv4
		info["ipv6"] = vps.IPv6
		info["spec_disktype"] = vps.Specs.DiskType
		info["spec_disksize"] = vps.Specs.DiskSize
		info["spec_cores"] = vps.Specs.Cores
		info["spec_ram"] = vps.Specs.RAM
		info["macs"] = vps.MACs

		ps = append(ps, info)

		return ps
	}

	return make([]interface{}, 0)
}
