package hetznerrobot

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func dataServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"server_number": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Server number",
			},
			"server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server name",
			},
			"server_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server IP",
			},
			"server_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server IPv6 Net",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data center",
			},
			"is_cancelled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status of server cancellation",
			},
			"paid_until": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Paid until date",
			},
			"product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server product name",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of assigned single IP addresses",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"server_subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of assigned subnets",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mask": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server status (\"ready\" or \"in process\")",
			},
			"traffic": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Free traffic quota, 'unlimited' in case of unlimited traffic",
			},
			"linked_storagebox": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Linked Storage Box ID",
			},
			"reset": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag of reset system availability",
			},
			"rescue": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag of Rescue System availability",
			},
			"vnc": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag of VNC installation availability",
			},
			"windows": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag of Windows installation availability",
			},
			"plesk": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag of Plesk installation availability",
			},
			"cpanel": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag of cPanel installation availability",
			},
			"wol": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag of Wake On Lan availability",
			},
			"hot_swap": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag of Hot Swap availability",
			},
		},
	}
}

func dataSourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverNumber := d.Get("server_number").(int)

	server, err := c.getServer(ctx, serverNumber)
	if err != nil {
		return diag.Errorf("Unable to find Server with IP %d:\n\t %q", serverNumber, err)
	}
	d.Set("datacenter", server.DataCenter)
	d.Set("is_cancelled", server.Cancelled)
	d.Set("paid_until", server.PaidUntil)
	d.Set("product", server.Product)
	d.Set("ip_addresses", server.IPs)
	d.Set("server_ip", server.ServerIP)
	d.Set("server_ipv6", server.ServerIPv6)
	d.Set("server_name", server.ServerName)
	d.Set("server_subnets", server.Subnets)
	d.Set("status", server.Status)
	d.Set("traffic", server.Traffic)
	d.Set("linked_storagebox", server.LinkedStoragebox)
	d.Set("reset", server.Reset)
	d.Set("rescue", server.Rescue)
	d.Set("vnc", server.VNC)
	d.Set("windows", server.Windows)
	d.Set("plesk", server.Plesk)
	d.Set("cpanel", server.CPanel)
	d.Set("wol", server.Wol)
	d.Set("hot_swap", server.HotSwap)
	d.SetId(strconv.Itoa(server.ServerNumber))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func dataServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServersRead,
		Schema: map[string]*schema.Schema{
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"server_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_cancelled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"paid_until": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"server_subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mask": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"linked_storagebox": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"reset": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rescue": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vnc": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"windows": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"plesk": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cpanel": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"wol": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hot_swap": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(HetznerRobotClient)

	servers, err := client.getServers(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	serverList := make([]map[string]interface{}, len(servers))
	for i, server := range servers {
		subnets := make([]map[string]interface{}, len(server.Subnets))
		for j, subnet := range server.Subnets {
			subnets[j] = map[string]interface{}{
				"ip":   subnet.IP,
				"mask": subnet.Mask,
			}
		}

		serverMap := map[string]interface{}{
			"server_number":     server.ServerNumber,
			"server_ip":         server.ServerIP,
			"server_ipv6":       server.ServerIPv6,
			"datacenter":        server.DataCenter,
			"is_cancelled":      server.Cancelled,
			"paid_until":        server.PaidUntil,
			"product":           server.Product,
			"ip_addresses":      server.IPs,
			"server_subnets":    subnets,
			"status":            server.Status,
			"traffic":           server.Traffic,
			"linked_storagebox": server.LinkedStoragebox,
			"reset":             server.Reset,
			"rescue":            server.Rescue,
			"vnc":               server.VNC,
			"windows":           server.Windows,
			"plesk":             server.Plesk,
			"cpanel":            server.CPanel,
			"wol":               server.Wol,
			"hot_swap":          server.HotSwap,
		}
		serverList[i] = serverMap
	}

	if err := d.Set("servers", serverList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("servers")

	var diags diag.Diagnostics
	return diags
}
