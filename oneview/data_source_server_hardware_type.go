// (C) Copyright 2019 Hewlett Packard Enterprise Development LP
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
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceServerHardwareType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerHardwareTypeRead,

		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_capability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"controller_modes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"drive_technologies": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"raid_levels": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServerHardwareTypeRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	name := d.Get("name").(string)

	server_hardware_type, err := config.ovClient.GetServerHardwareTypeByName(name)
	if err != nil || server_hardware_type.URI.IsNil() {
		d.SetId("")
		return nil
	}
	d.SetId(name)
	d.Set("name", server_hardware_type.Name)
	d.Set("description", server_hardware_type.Description.String())
	d.Set("category", server_hardware_type.Category)
	d.Set("etag", server_hardware_type.ETAG)
	d.Set("uri", server_hardware_type.URI.String())

	controllerModes := make([]interface{}, len(server_hardware_type.StorageCapabilities.ControllerModes))
	for i, controllerMode := range server_hardware_type.StorageCapabilities.ControllerModes {
		controllerModes[i] = controllerMode
	}
	driveTechnologies := make([]interface{}, len(server_hardware_type.StorageCapabilities.DriveTechnologies))
	for i, driveTechnology := range server_hardware_type.StorageCapabilities.DriveTechnologies {
		driveTechnologies[i] = driveTechnology
	}
	raidLevels := make([]interface{}, len(server_hardware_type.StorageCapabilities.RaidLevels))
	for i, raidLevel := range server_hardware_type.StorageCapabilities.RaidLevels {
		raidLevels[i] = raidLevel
	}

	storageCapability := make([]map[string]interface{}, 0, 1)
	storageCapability = append(storageCapability, map[string]interface{}{
		"controller_modes":   schema.NewSet(schema.HashString, controllerModes),
		"drive_technologies": schema.NewSet(schema.HashString, driveTechnologies),
		"raid_levels":        schema.NewSet(schema.HashString, raidLevels),
	})
	d.Set("storage_capability", storageCapability)
	return nil
}
