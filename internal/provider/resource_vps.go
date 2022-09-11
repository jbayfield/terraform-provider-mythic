package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jbayfield/mythic-client-go"
)

func resourceVPS() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPSCreate,
		ReadContext:   resourceVPSRead,
		UpdateContext: resourceVPSUpdate,
		DeleteContext: resourceVPSDelete,

		// TODO: Work out which of the forcenew/computed/required are redundant
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
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
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"netdevice": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"diskbus": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"isoimage": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			"bootdevice": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
			"disktype": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "ssd",
			},
			"disksize": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				Computed: false,
				ForceNew: false,
			},
			"cores": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram": &schema.Schema{
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

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVPSCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*mythic.Client)

	var diags diag.Diagnostics

	// TODO: Provide all of the things
	vpsspec := mythic.VPSCreateSpec{
		Identifier: d.Id(),
		Product:    d.Get("product").(string),
		DiskSize:   d.Get("disksize").(int),
	}

	vps, err := c.CreateVPS(vpsspec, &c.Token)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vps.Identifier)
	resourceVPSRead(ctx, d, meta)

	return diags
}

func resourceVPSRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*mythic.Client)

	var diags diag.Diagnostics

	vps, err := c.GetVPS(d.Id(), &c.Token)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vps.Identifier)
	d.Set("name", vps.Name)
	d.Set("product", vps.Product)
	d.Set("dormant", vps.Dormant)
	d.Set("status", vps.Status)
	d.Set("hostserver", vps.HostServer)
	d.Set("zone_code", vps.Zone.Code)
	d.Set("zone_name", vps.Zone.Name)
	d.Set("cpumode", vps.CPUMode)
	d.Set("netdevice", vps.NetDevice)
	d.Set("diskbus", vps.DiskBus)
	d.Set("isoimage", vps.ISOImage)
	d.Set("bootdevice", vps.BootDevice)
	d.Set("ipv4", vps.IPv4)
	d.Set("ipv6", vps.IPv6)
	d.Set("macs", vps.MACs)
	d.Set("disktype", vps.Specs.DiskType)
	d.Set("disksize", vps.Specs.DiskSize)
	d.Set("cores", vps.Specs.Cores)
	d.Set("ram", vps.Specs.RAM)

	// TODO: SetConnInfo with IPv4 if available or v6

	return diags
}

func resourceVPSUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*mythic.Client)

	var diags diag.Diagnostics

	// TODO: Provide all of the things
	vpsspec := mythic.VPSUpdateSpec{
		Identifier: d.Id(),
		Name:       d.Get("name").(string),
		Product:    d.Get("product").(string),
		DiskSize:   d.Get("disksize").(int),
	}

	err := c.UpdateVPS(vpsspec, &c.Token)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceVPSRead(ctx, d, meta)

	return diags
}

func resourceVPSDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*mythic.Client)

	var diags diag.Diagnostics

	err := c.DestroyVPS(d.Id(), &c.Token)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
