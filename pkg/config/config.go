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
package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/maksim-paskal/auto-devops/pkg/types"
	"github.com/maksim-paskal/auto-devops/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	giturls "github.com/whilp/git-urls"
	"gopkg.in/yaml.v3"
)

const (
	autoDevopsYaml = ".auto-devops.yaml"
)

//nolint: gochecknoglobals
var (
	gitVersion   = "dev"
	BootstrapZip = flag.String("bootstrap", os.Getenv("AUTO_DEVOPS_BOOTSTRAP"), "path to archive")
	LogLevel     = flag.String("log.level", "INFO", "log level")
	Boostrap     = types.Boostrap{
		DelimsLeft:  "{%",
		DelimsRight: "%}",
		Answers:     make(map[string]string),
	}
)

func GetVersion() string {
	return gitVersion
}

func Init() error { //nolint:cyclop,funlen
	log.SetReportCaller(true)

	level, err := log.ParseLevel(*LogLevel)
	if err != nil {
		return errors.Wrap(err, "error parsing level")
	}

	log.SetLevel(level)

	log.Debugf("Starting %s...", GetVersion())

	dir, err := ioutil.TempDir("", "auto-devops")
	if err != nil {
		return errors.Wrap(err, "error creating temp folder")
	}

	// load initial .auto-devops.yaml
	err = loadYAML(autoDevopsYaml)
	if err != nil {
		log.WithError(err).Debug("not loading initital config")
	}

	Boostrap.Dir = dir
	if len(Boostrap.Version) == 0 {
		Boostrap.Version = GetVersion()
	}

	if len(*BootstrapZip) > 0 {
		Boostrap.Bootstrap = *BootstrapZip
	}

	if len(Boostrap.Bootstrap) == 0 {
		log.Fatalf("neen bootstrap zip. Use -bootstrap argument or bootstrap attribute in %s", autoDevopsYaml)
	}

	_, err = utils.Unzip(Boostrap.Bootstrap, Boostrap.Dir)
	if err != nil {
		return errors.Wrapf(err, "error unzip %s", Boostrap.Bootstrap)
	}

	// load server .auto-devops.yaml
	autoDevopsFile := filepath.Join(Boostrap.Dir, autoDevopsYaml)

	err = loadYAML(autoDevopsFile)
	if err != nil {
		return errors.Wrap(err, "error loading server file")
	}

	// load user .auto-devops.yaml
	if _, err := os.Stat(autoDevopsYaml); err == nil {
		err = loadYAML(autoDevopsYaml)
		if err != nil {
			return errors.Wrap(err, "error loading user file")
		}
	}

	err = loadYAML(autoDevopsYaml)
	if err != nil {
		log.WithError(err).Debug("error reading user ", autoDevopsYaml)
	}

	if len(Boostrap.Pwd) == 0 {
		Boostrap.Pwd, err = os.Getwd()
		if err != nil {
			return errors.Wrap(err, "error getting current folder")
		}
	}

	if len(Boostrap.Name) == 0 {
		info, err := os.Stat(Boostrap.Pwd)
		if err != nil {
			return errors.Wrap(err, "error getting current folder stat")
		}

		Boostrap.Name = info.Name()
	}

	err = loadGitInfo()
	if err != nil {
		log.WithError(err).Warn("error loading git info")
	}

	return nil
}

func loadYAML(yamlPath string) error {
	configByte, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return errors.Wrap(err, "error reading file")
	}

	err = yaml.Unmarshal(configByte, &Boostrap)
	if err != nil {
		return errors.Wrap(err, "error parse yaml")
	}

	return nil
}

var (
	errHasNoKey     = errors.New("has no key")
	errWrongVersion = errors.New("wrong version")
)

func Validate() error {
	for _, q := range Boostrap.Questions {
		log.Debugf("key=%s", q.Key)

		if len(q.Key) == 0 {
			return errors.Wrap(errHasNoKey, "question "+q.Prompt)
		}

		if _, err := Boostrap.Template(q.Condition); err != nil {
			return errors.Wrap(err, q.Condition)
		}
	}

	for _, f := range Boostrap.Filters {
		if _, err := Boostrap.Template(f.Condition); err != nil {
			return errors.Wrap(err, f.Condition)
		}
	}

	matched, err := regexp.MatchString(Boostrap.Version, GetVersion())
	if err != nil {
		return errors.Wrap(err, "error in version matching")
	}

	if !matched {
		return errors.Wrap(
			errWrongVersion,
			fmt.Sprintf("required version %s,current version %s. Please install required version", Boostrap.Version, GetVersion()), //nolint:lll
		)
	}

	return nil
}

func loadGitInfo() error {
	r, err := git.PlainOpen(".")
	if err != nil {
		return errors.Wrap(err, "error opening folder")
	}

	list, err := r.Remotes()
	if err != nil {
		return errors.Wrap(err, "error listing remotes")
	}

	if len(list) == 0 {
		return errors.New("no remote")
	}

	if len(list[0].Config().URLs) == 0 {
		return errors.New("no remote urls")
	}

	remoteURL := list[0].Config().URLs[0]

	u, err := giturls.Parse(remoteURL)
	if err != nil {
		return errors.Wrap(err, "error parsing string")
	}

	Boostrap.GitInfo.Host = u.Host
	Boostrap.GitInfo.Path = u.Path
	Boostrap.GitInfo.PathFormated = strings.TrimSuffix(Boostrap.GitInfo.Path, ".git")

	return nil
}
