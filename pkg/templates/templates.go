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
package templates

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/maksim-paskal/auto-devops/pkg/config"
	"github.com/maksim-paskal/auto-devops/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func TemplateFile(bootstrapPath string, info os.FileInfo, fileMode uint32, dryRun bool) (string, error) {
	dest := path.Join(config.Boostrap.Dir, "templates")
	destAbs := strings.TrimPrefix(bootstrapPath, dest)
	destAbs = path.Join(config.Boostrap.Pwd, destAbs)

	log.Debug(destAbs)

	t := template.New(info.Name())

	tmpl, err := t.
		Funcs(utils.GoTemplateFunc(t)).
		Delims(config.Boostrap.DelimsLeft, config.Boostrap.DelimsRight).
		ParseFiles(bootstrapPath)
	if err != nil {
		return "", errors.Wrap(err, "error templating file")
	}

	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, config.Boostrap)
	if err != nil {
		return "", errors.Wrap(err, "error executing templating")
	}

	if !dryRun {
		if err = os.MkdirAll(filepath.Dir(destAbs), os.ModePerm); err != nil {
			return "", errors.Wrap(err, "error create directory")
		}

		err = ioutil.WriteFile(destAbs, tpl.Bytes(), os.FileMode(fileMode))
		if err != nil {
			return "", errors.Wrap(err, "error writing file")
		}
	}

	return destAbs, nil
}
