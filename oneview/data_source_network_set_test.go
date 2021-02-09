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
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccNetworkSet_2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSetData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("oneview_network_set.test", "name", "network_set_1"),
				),
			},
		},
	})
}

var testAccNetworkSetData = `
  data "oneview_network_set" "test" {	  
    name = "network_set_1"
  }`
