package hetznerrobot

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceVSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVSwitchCreate,
		ReadContext:   resourceVSwitchRead,
		UpdateContext: resourceVSwitchUpdate,
		DeleteContext: resourceVSwitchDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceVSwitchImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vSwitch name",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "VLAN ID",
			},
			// computed / read-only fields
			"is_cancelled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Cancellation status",
			},
			"servers": {
				Type:        schema.TypeList,
				Description: "Attached server list",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_number": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"server_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_ipv6_net": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
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
func resourceVSwitchImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Id()
	vSwitch, err := c.getVSwitch(ctx, vSwitchID)
	if err != nil {
		return nil, fmt.Errorf("Unable to find VSwitch with ID %d:\n\t %q", vSwitchID, err)
	}

	d.Set("name", vSwitch.Name)
	d.Set("vlan", vSwitch.Vlan)
	d.Set("is_cancelled", vSwitch.Cancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	name := d.Get("name").(string)
	vlan := d.Get("vlan").(int)
	vSwitch, err := c.createVSwitch(ctx, name, vlan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to create VSwitch :\n\t %q", err))
	}

	d.Set("is_cancelled", vSwitch.Cancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)
	d.SetId(strconv.Itoa(vSwitch.ID))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceVSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Id()
	vSwitch, err := c.getVSwitch(ctx, vSwitchID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to find VSwitch with ID %s:\n\t %q", vSwitchID, err))
	}

	d.Set("name", vSwitch.Name)
	d.Set("vlan", vSwitch.Vlan)
	d.Set("cancelled", vSwitch.Cancelled)
	d.Set("servers", vSwitch.Server)
	d.Set("subnets", vSwitch.Subnet)
	d.Set("cloud_networks", vSwitch.CloudNetwork)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceVSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Id()
	name := d.Get("name").(string)
	vlan := d.Get("vlan").(int)
	err := c.updateVSwitch(ctx, vSwitchID, name, vlan)
	if err != nil {
		return diag.Errorf("Unable to update VSwitch:\n\t %q", err)
	}

	if d.HasChange("servers") {
		o, n := d.GetChange("servers")

		oldServers := o.([]interface{})
		newServers := n.([]interface{})

		mb := make(map[int]struct{}, len(newServers))
		for _, x := range newServers {
			srv := x.(map[string]interface{})
			mb[srv["server_number"].(int)] = struct{}{}
		}
		var serversToRemove []HetznerRobotVSwitchServer
		for _, x := range oldServers {
			srv := x.(map[string]interface{})
			srvNum := srv["server_number"].(int)
			if _, found := mb[srvNum]; !found {
				serversToRemove = append(serversToRemove, HetznerRobotVSwitchServer{ServerNumber: srvNum})
			}
		}

		if err := c.removeVSwitchServers(ctx, vSwitchID, serversToRemove); err != nil {
			diag.Errorf("Unable to remove servers from VSwitch:\n\t %q", err)
		}

		ma := make(map[int]struct{}, len(oldServers))
		for _, x := range oldServers {
			srv := x.(map[string]interface{})
			ma[srv["server_number"].(int)] = struct{}{}
		}
		var serversToAdd []HetznerRobotVSwitchServer
		for _, x := range newServers {
			srv := x.(map[string]interface{})
			srvNum := srv["server_number"].(int)
			if _, found := ma[srvNum]; !found {
				serversToAdd = append(serversToAdd, HetznerRobotVSwitchServer{ServerNumber: srvNum})
			}
		}

		if err := c.addVSwitchServers(ctx, vSwitchID, serversToAdd); err != nil {
			diag.Errorf("Unable to add servers to VSwitch:\n\t %q", err)
		}
	}

	return resourceVSwitchRead(ctx, d, meta)
}

func resourceVSwitchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	vSwitchID := d.Id()
	err := c.deleteVSwitch(ctx, vSwitchID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Unable to find VSwitch with ID %s:\n\t %q", vSwitchID, err))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
