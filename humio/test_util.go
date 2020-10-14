// Copyright © 2020 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package humio

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider

func init() {
	testAccProviders = map[string]*schema.Provider{
		"humio": Provider(),
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("HUMIO_ADDR") == "" {
		t.Fatal("HUMIO_ADDR must be set for acceptance tests")
	}
	if os.Getenv("HUMIO_API_TOKEN") == "" {
		t.Fatal("HUMIO_API_TOKEN must be set for acceptance tests")
	}
}

func accTestCase(t *testing.T, steps []resource.TestStep, checkDestroyFunc resource.TestCheckFunc) {
	resource.Test(t, resource.TestCase{
		CheckDestroy: checkDestroyFunc,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps:     steps,
	})
}
