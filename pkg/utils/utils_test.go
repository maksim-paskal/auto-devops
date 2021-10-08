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
package utils_test

import (
	"testing"
	"text/template"

	"github.com/maksim-paskal/auto-devops/pkg/utils"
)

func TestGoTemplateFunc(t *testing.T) {
	t.Parallel()

	tmpl := template.New("test")

	fun := utils.GoTemplateFunc(tmpl)

	_, ok := fun["toYaml"].(func(v interface{}) string)
	if !ok {
		t.Fatal("toYaml not correct")
	}

	_, ok = fun["include"].(func(name string, data interface{}) (string, error))
	if !ok {
		t.Fatal("include not correct")
	}
}
