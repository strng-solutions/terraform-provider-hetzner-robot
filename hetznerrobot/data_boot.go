package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataBoot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBootRead,
		Schema: map[string]*schema.Schema{
			// read-only / computed
			"active_profile": {
				Type:        schema.TypeString, // Enum should be better (linux/rescue/...)
				Computed:    true,
				Description: "Active boot profile",
			},
			"ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server main IPv4 address",
			},
			"ipv6_network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server main IPv6 net address",
			},
			"language": {
				Type:        schema.TypeString, // Enum should be better (amd64/...)
				Computed:    true,
				Description: "Language",
			},
			"operating_system": {
				Type:        schema.TypeString, // Enum should be better (ubuntu_20.04/...)
				Computed:    true,
				Description: "Active Operating System / Distribution",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current Rescue System root password / Linux installation password or null",
				Sensitive:   true,
			},
		},
		/*
			AuthorizedKeys []string		    authorized_key (Array)	Authorized public SSH keys
			HostKeys []string				host_key (Array)	Host keys
		*/
	}
}
func dataSourceBootRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverNumber, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	boot, err := c.getBoot(ctx, serverNumber)
	if err != nil {
		return diag.Errorf("Unable to find Boot Profile for server ID %d:\n\t %q", serverNumber, err)
	}

	d.Set("active_profile", boot.ActiveProfile)
	d.Set("ipv4_address", boot.ServerIPv4)
	d.Set("ipv6_network", boot.ServerIPv6)
	d.Set("language", boot.Language)
	d.Set("operating_system", boot.OperatingSystem)
	d.Set("password", boot.Password)
	d.SetId(strconv.Itoa(serverNumber))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
