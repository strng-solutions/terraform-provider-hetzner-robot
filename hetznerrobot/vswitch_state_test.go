package hetznerrobot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestVSwitchState(t *testing.T) {
	for _, operation := range []string{"read", "import", "create", "data"} {
		for _, empty := range []bool{false, true} {
			t.Run(operation+map[bool]string{false: "/populated", true: "/empty"}[empty], func(t *testing.T) {
				response := HetznerRobotVSwitch{ID: 42, Name: "test-switch", Vlan: 4000, Cancelled: true}
				if !empty {
					response.Server = []HetznerRobotVSwitchServer{
						{ServerNumber: 200, ServerIP: "192.0.2.2", ServerIPv6Net: "2001:db8:2::/64", Status: "ready"},
						{ServerNumber: 100, ServerIP: "192.0.2.1", ServerIPv6Net: "2001:db8:1::/64", Status: "ready"},
					}
					response.Subnet = []HetznerRobotVSwitchSubnet{{IP: "198.51.100.0", Mask: 24, Gateway: "198.51.100.1"}}
					response.CloudNetwork = []HetznerRobotVSwitchCloudNetwork{{ID: 7, IP: "10.0.0.0", Mask: 24, Gateway: "10.0.0.1"}}
				}
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					wantMethod, wantPath := http.MethodGet, "/vswitch/42"
					if operation == "create" {
						wantMethod, wantPath = http.MethodPost, "/vswitch"
					}
					if r.Method != wantMethod || r.URL.Path != wantPath {
						t.Errorf("request = %s %s, want %s %s", r.Method, r.URL.Path, wantMethod, wantPath)
					}
					w.Header().Set("Content-Type", "application/json")
					if err := json.NewEncoder(w).Encode(response); err != nil {
						t.Error(err)
					}
				}))
				defer server.Close()
				client := NewHetznerRobotClient("", "", server.URL)
				resource := resourceVSwitch()
				if operation == "data" {
					resource = dataVSwitch()
				}
				d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{"name": "test-switch", "vlan": 4000})
				d.SetId("42")
				ctx := context.Background()
				switch operation {
				case "import":
					states, err := resourceVSwitchImportState(ctx, d, client)
					if err != nil || len(states) != 1 || states[0] != d {
						t.Fatalf("import states = %v, error = %v", states, err)
					}
				case "create":
					d.SetId("")
					if diags := resourceVSwitchCreate(ctx, d, client); diags.HasError() {
						t.Fatal(diags)
					}
				default:
					if diags := resource.ReadContext(ctx, d, client); diags.HasError() {
						t.Fatal(diags)
					}
				}
				if d.Id() != "42" {
					t.Errorf("ID = %q, want 42", d.Id())
				}
				if d.Get("is_cancelled") != true {
					t.Errorf("is_cancelled = %v, want true", d.Get("is_cancelled"))
				}
				wantServers, wantSubnets, wantNetworks := []interface{}{}, []interface{}{}, []interface{}{}
				if !empty {
					wantServers = []interface{}{
						map[string]interface{}{"server_number": 100, "server_ip": "192.0.2.1", "server_ipv6_net": "2001:db8:1::/64", "status": "ready"},
						map[string]interface{}{"server_number": 200, "server_ip": "192.0.2.2", "server_ipv6_net": "2001:db8:2::/64", "status": "ready"},
					}
					wantSubnets = []interface{}{map[string]interface{}{"ip": "198.51.100.0", "mask": 24, "gateway": "198.51.100.1"}}
					wantNetworks = []interface{}{map[string]interface{}{"id": 7, "ip": "10.0.0.0", "mask": 24, "gateway": "10.0.0.1"}}
				}
				for key, want := range map[string][]interface{}{"servers": wantServers, "subnets": wantSubnets, "cloud_networks": wantNetworks} {
					if got := d.Get(key); !reflect.DeepEqual(got, want) {
						t.Errorf("%s = %#v, want %#v", key, got, want)
					}
				}
			})
		}
	}
}

func TestVSwitchReadPreservesConfiguredServerOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(HetznerRobotVSwitch{
			ID: 42, Name: "test-switch", Vlan: 4000,
			Server: []HetznerRobotVSwitchServer{{ServerNumber: 100}, {ServerNumber: 300}, {ServerNumber: 200}},
		})
	}))
	defer server.Close()
	d := schema.TestResourceDataRaw(t, resourceVSwitch().Schema, map[string]interface{}{
		"name": "test-switch", "vlan": 4000,
		"servers": []interface{}{
			map[string]interface{}{"server_number": 200},
			map[string]interface{}{"server_number": 100},
		},
	})
	d.SetId("42")
	client := NewHetznerRobotClient("", "", server.URL)
	for i := 0; i < 2; i++ {
		if diags := resourceVSwitchRead(context.Background(), d, client); diags.HasError() {
			t.Fatal(diags)
		}
		servers := d.Get("servers").([]interface{})
		got := make([]int, len(servers))
		for i, server := range servers {
			got[i] = server.(map[string]interface{})["server_number"].(int)
		}
		if want := []int{200, 100, 300}; !reflect.DeepEqual(got, want) {
			t.Errorf("refresh server order = %v, want %v", got, want)
		}
	}
}
