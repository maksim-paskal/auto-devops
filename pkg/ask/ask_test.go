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
package ask_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/maksim-paskal/auto-devops/pkg/ask"
	"github.com/maksim-paskal/auto-devops/pkg/config"
	"github.com/maksim-paskal/auto-devops/pkg/types"
)

func TestOnce(t *testing.T) {
	t.Parallel()

	var tpl bytes.Buffer

	ask.Output = &tpl
	ask.Reader = bufio.NewReader(strings.NewReader("user-test-input"))

	config.Boostrap.Name = "test"

	q := types.BoostrapQuestion{
		Prompt:     "test",
		Result:     "default-answer",
		Validation: ".+",
	}

	ok, err := ask.Once(&q)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("not correct")
	}

	if q.Result != "user-test-input" {
		t.Fatal("not correct result")
	}

	if show := tpl.String(); show != "test .+\n[default-answer]: " {
		t.Fatal("question does not show corrected: " + show)
	}
}
