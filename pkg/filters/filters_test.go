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
package filters_test

import (
	"testing"

	"github.com/maksim-paskal/auto-devops/pkg/config"
	"github.com/maksim-paskal/auto-devops/pkg/filters"
	"github.com/maksim-paskal/auto-devops/pkg/types"
)

func TestGetFilter(t *testing.T) {
	t.Parallel()

	config.Boostrap = types.Boostrap{
		Filters: []types.BoostrapFilter{
			{
				Match:     "test1.yaml",
				Condition: "{{ false  }}",
			},
			{
				Match:     "test2.yaml",
				Condition: "{{ true  }}",
			},
			{
				Match:  "test3.yaml",
				Ignore: true,
			},
		},
	}

	f1, err := filters.GetFilter("/test/test1.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if len(f1.Match) != 0 {
		t.Fatal("filter must be empty", f1)
	}

	f2, err := filters.GetFilter("/test/test2.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if len(f2.Match) == 0 {
		t.Fatal("filter must be not empty")
	}

	f3, err := filters.GetFilter("/test/test3.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if f3.Ignore != true {
		t.Fatal("filter must be return Ignore=true")
	}
}
