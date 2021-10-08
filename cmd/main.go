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
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/maksim-paskal/auto-devops/pkg/ask"
	"github.com/maksim-paskal/auto-devops/pkg/config"
	"github.com/maksim-paskal/auto-devops/pkg/filters"
	"github.com/maksim-paskal/auto-devops/pkg/templates"
	log "github.com/sirupsen/logrus"
)

var version = flag.Bool("version", false, "version")

func main() { //nolint:cyclop,funlen
	flag.Parse()

	if *version {
		fmt.Println(config.GetVersion()) //nolint:forbidigo
		os.Exit(0)
	}

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	if log.GetLevel() >= log.DebugLevel {
		log.Debugf("Loaded config:\n%s", config.Boostrap.String())
	}

	for i, q := range config.Boostrap.Questions {
		if len(q.Condition) > 0 {
			show, err := config.Boostrap.Template(q.Condition)
			if err != nil {
				log.Fatal(err)
			}

			if show == "false" {
				continue
			}
		}

		for {
			ok, err := ask.Once(&config.Boostrap.Questions[i])
			if err != nil {
				log.Fatal(err)
			}

			if ok {
				config.Boostrap.Answers[q.Key] = config.Boostrap.Questions[i].Result

				break
			}
		}
	}

	// test all templates
	processAnswers(true)

	// save templates
	processAnswers(false)

	readmePath := path.Join(config.Boostrap.Dir, "README")

	if _, err := os.Stat(readmePath); err == nil {
		readmeBytes, err := ioutil.ReadFile(readmePath)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("README:\n%s", string(readmeBytes))
	}

	if err := config.Boostrap.CleanTempDir(); err != nil {
		log.Fatal(err)
	}
}

func processAnswers(dryRun bool) {
	templatesDir := path.Join(config.Boostrap.Dir, "templates")

	err := filepath.Walk(templatesDir, func(filePath string, info os.FileInfo, _ error) error {
		log.Debug(filePath)

		if info.IsDir() {
			return nil
		}

		dest := path.Join(config.Boostrap.Dir, "templates")
		filePathAbs := strings.TrimPrefix(filePath, dest)

		filterRule, err := filters.GetFilter(filePathAbs)
		if err != nil {
			log.Fatal(err)
		}

		if len(filterRule.Match) > 0 {
			if filterRule.Ignore {
				return nil
			}
		}

		var fileMode uint32 = 0o644

		if filterRule.FileMode > 0 {
			fileMode = filterRule.FileMode
		}

		_, err = templates.TemplateFile(filePath, info, fileMode, dryRun)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
