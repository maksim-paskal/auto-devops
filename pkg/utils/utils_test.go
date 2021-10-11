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
	"io/ioutil"
	"os"
	"testing"
	"text/template"

	"github.com/maksim-paskal/auto-devops/pkg/utils"
)

func TestGoTemplateFunc(t *testing.T) {
	t.Parallel()

	tmpl := template.New("test")

	fun := utils.GoTemplateFunc(tmpl)

	toYaml, ok := fun["toYaml"].(func(v interface{}) string)
	if !ok {
		t.Fatal("toYaml not correct")
	}

	type Test struct {
		Test string
	}

	if r := toYaml(Test{Test: "test"}); r != "test: test\n" {
		t.Fatal("toYaml returns wrong results")
	}

	randPort, ok := fun["randPort"].(func() int)
	if !ok {
		t.Fatal("randPort not correct")
	}

	if r := randPort(); r < utils.RandPortMin || r > utils.RandPortMax {
		t.Fatal("randPort rerurns wrong results")
	}

	tpl, ok := fun["tpl"].(func(tpl string, data interface{}) string)
	if !ok {
		t.Fatal("tpl not correct")
	}

	if r := tpl("1{{ .Test }}2", Test{Test: "test"}); r != "1test2" {
		t.Fatal("tpl not correct")
	}
}

func TestUnzipFile(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "auto-devops-tests")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	testZip, err := utils.Unzip("testdata/test.zip", dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(testZip) != 1 {
		t.Fatal("zip must contain 1 file")
	}

	testZipContents, err := ioutil.ReadFile(testZip[0])
	if err != nil {
		t.Fatal(err)
	}

	if string(testZipContents) != "test\n" {
		t.Fatal("test file incorrect")
	}

	_, err = utils.Unzip("https://github.com/maksim-paskal/auto-devops/archive/refs/heads/main.zip", dir)
	if err != nil {
		t.Fatal(err)
	}
}
