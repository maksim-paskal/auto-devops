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
package types_test

import (
	"strings"
	"testing"

	"github.com/maksim-paskal/auto-devops/pkg/types"
)

func TestString(t *testing.T) {
	t.Parallel()

	boostrap := types.Boostrap{
		Name: "test",
	}

	if !strings.Contains(boostrap.String(), "name: test") {
		t.Fatal("String method working not correct")
	}
}

func TestTemplate(t *testing.T) {
	t.Parallel()

	boostrap := types.Boostrap{
		Name: "test",
	}

	r, err := boostrap.Template("{{ .Name }}")
	if err != nil {
		t.Fatal(err)
	}

	if r != "test" {
		t.Fatal("Template method working not correct")
	}
}
