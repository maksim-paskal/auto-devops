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
package templates_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/maksim-paskal/auto-devops/pkg/config"
	"github.com/maksim-paskal/auto-devops/pkg/templates"
)

const fileName = "testdata/test.yaml"

func TestTemplateFile(t *testing.T) {
	t.Parallel()

	info, err := os.Stat(fileName)
	if err != nil {
		t.Fatal(err)
	}

	config.Boostrap.Name = "test"
	config.Boostrap.DelimsLeft = "<<"
	config.Boostrap.DelimsRight = ">>"
	config.Boostrap.Pwd = path.Join(os.TempDir(), "auto-devops-test")

	tmplFile, err := templates.TemplateFile(fileName, info, 0o766, false)
	if err != nil {
		t.Fatal(err)
	}

	tmplFileBytes, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(config.Boostrap.Pwd)
	if err != nil {
		t.Fatal(err)
	}

	if result := string(tmplFileBytes); result != "test" {
		t.Fatal("templating file is not correct result=" + result)
	}
}
