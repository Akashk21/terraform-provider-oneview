// (C) Copyright 2021 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package oneview

import (
	"github.com/HewlettPackard/oneview-golang/ov"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceIPv4Subnets() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPv4SubnetCreate,
		Read:   resourceIPv4SubnetRead,
		Update: resourceIPv4SubnetUpdate,
		Delete: resourceIPv4SubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"allocator_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associated_resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"association_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"collector_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_servers": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"range_uris": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subnet_mask": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIPv4SubnetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	subnet := ov.Ipv4Subnet{
		Name:       d.Get("name").(string),
		Domain:     d.Get("domain").(string),
		Gateway:    d.Get("gateway").(string),
		NetworkId:  d.Get("network_id").(string),
		SubnetMask: d.Get("subnet_mask").(string),
	}

	err := config.ovClient.CreateIPv4Subnet(subnet)
	if err != nil {
		d.SetId("")
		return err
	}
	sNet, err := config.ovClient.GetSubnetByNetworkId(d.Get("network_id").(string))
	uri := sNet.URI.String()
	subnet_id := strings.Split(uri, "/")[5]
	d.SetId(subnet_id)

	return resourceIPv4SubnetRead(d, meta)
}

func resourceIPv4SubnetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	subnet, err := config.ovClient.GetIPv4SubnetbyId(d.Id())
	if err != nil || subnet.URI.IsNil() {
		d.SetId("")
		return nil
	}

	associatedRes := make([]map[string]interface{}, 0, len(subnet.AssociatedResources))
	for _, res := range subnet.AssociatedResources {
		associatedRes = append(associatedRes, map[string]interface{}{
			"association_type":  res.AssociationType,
			"resource_category": res.ResourceCategory,
			"resource_name":     res.ResourceName,
			"resource_uri":      res.ResourceUri,
		})
	}

	d.Set("allocator_uri", subnet.AllocatorUri)
	d.Set("associated_resources", associatedRes)
	d.Set("category", subnet.Category)
	d.Set("collector_uri", subnet.CollectorUri)
	d.Set("created", subnet.Created)
	d.Set("dns_servers", subnet.DnsServers)
	d.Set("domain", subnet.Domain)
	d.Set("etag", subnet.ETAG)
	d.Set("gateway", subnet.Gateway)
	d.Set("modified", subnet.Modified)
	d.Set("name", subnet.Name)
	d.Set("network_id", subnet.NetworkId)
	d.Set("range_uris", subnet.RangeUris)
	d.Set("subnet_mask", subnet.SubnetMask)
	d.Set("type", subnet.Type)
	d.Set("uri", subnet.URI.String())
	d.SetId(d.Id())
	return nil
}

func resourceIPv4SubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	subnet := ov.Ipv4Subnet{
		Name:       d.Get("name").(string),
		Gateway:    d.Get("gateway").(string),
		NetworkId:  d.Get("network_id").(string),
		SubnetMask: d.Get("subnet_mask").(string),
	}

	err := config.ovClient.UpdateIpv4Subnet(d.Id(), subnet)

	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(d.Id())
	return resourceIPv4SubnetRead(d, meta)
}

func resourceIPv4SubnetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	err := config.ovClient.DeleteIpv4Subnet(d.Id())

	if err != nil {
		d.SetId("")
		return err
	}

	return nil
}
