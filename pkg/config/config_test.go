/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package config_test

import (
	"testing"

	"github.com/maksim-paskal/auto-devops/pkg/config"
)

func TestInit(t *testing.T) {
	t.Parallel()

	if err := config.Init(); err != nil {
		t.Fatal(err)
	}

	if err := config.Boostrap.CleanTempDir(); err != nil {
		t.Fatal(err)
	}

	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	if v := config.GetVersion(); v != "dev" {
		t.Fatal("version not correct")
	}
}
