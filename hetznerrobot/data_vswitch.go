package hetznerrobot

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataVSwitch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVSwitchRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vSwitch name",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VLAN ID",
			},
			"is_cancelled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Cancellation status",
			},
			"servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached server list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server_ipv6_net": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server_number": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"status": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached subnet list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mask": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"cloud_networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached cloud network list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mask": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Id()
	vSwitch, err := c.getVSwitch(ctx, vSwitchID)
	if err != nil {
		return diag.Errorf("Unable to find VSwitch with ID %s:\n\t %q", vSwitchID, err)
	}

	d.Set("name", vSwitch.Name)
	d.Set("vlan", vSwitch.Vlan)
	d.Set("is_cancelled", vSwitch.Cancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)
	d.SetId(vSwitchID)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
