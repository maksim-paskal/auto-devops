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
package types

import (
	"bytes"
	"os"
	"text/template"

	"github.com/maksim-paskal/auto-devops/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type BoostrapQuestion struct {
	Key        string
	Prompt     string
	Result     string
	Validation string
	Condition  string
}

type BoostrapFilter struct {
	Match     string
	Condition string
	Ignore    bool
	FileMode  uint32
}

type BoostrapGitInfo struct {
	Host         string
	Path         string
	PathFormated string
}

type Boostrap struct {
	Version     string
	Name        string
	Dir         string
	Pwd         string
	DelimsLeft  string
	DelimsRight string
	Filters     []BoostrapFilter
	Questions   []BoostrapQuestion
	Answers     map[string]string
	GitInfo     BoostrapGitInfo
	User        map[string]string
}

func (b Boostrap) String() string {
	out, err := yaml.Marshal(b)
	if err != nil {
		log.Error(err)
	}

	return string(out)
}

func (b Boostrap) Template(templateText string) (string, error) {
	t := template.New("template")

	tmpl, err := t.Funcs(utils.GoTemplateFunc(t)).Parse(templateText)
	if err != nil {
		return "", errors.Wrap(err, "error in template.parse")
	}

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, b)
	if err != nil {
		return "", errors.Wrap(err, "error in template.execute")
	}

	return tpl.String(), nil
}

func (b Boostrap) CleanTempDir() error {
	log.Debug("Cleaning tmp dir...")

	err := os.RemoveAll(b.Dir)

	return errors.Wrap(err, "error removing temp dir")
}
