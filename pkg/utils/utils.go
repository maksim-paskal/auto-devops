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
package utils

import (
	"archive/zip"
	"context"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var errIlligalFilePath = errors.New("illigal file path")

func downloadFile(filepath string, url string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "error in crating new requests")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error in doing requests")
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "error in os.Create")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return errors.Wrap(err, "error in io.Copy")
}

func Unzip(src string, dest string) ([]string, error) { //nolint: funlen,cyclop
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		log.Debug("Downloading file")

		downloadedFile := path.Join(dest, "bootstrap.zip")

		if err := downloadFile(path.Join(dest, "bootstrap.zip"), src); err != nil {
			return nil, errors.Wrap(err, "error downloading file")
		}

		src = downloadedFile
	}

	filenames := make([]string, 0)

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, errors.Wrap(err, "error open zip reader "+src)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name) //nolint:gosec

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, errors.Wrap(errIlligalFilePath, fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
				return nil, errors.Wrap(err, "error creating directory")
			}

			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, errors.Wrap(err, "error creating directory")
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, errors.Wrap(err, "error creating file")
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, errors.Wrap(err, "error open file")
		}

		_, err = io.Copy(outFile, rc) //nolint: gosec

		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, errors.Wrap(err, "error copy file")
		}
	}

	return filenames, nil
}

const (
	RandPortMin = 10000
	RandPortMax = 50000
)

func GoTemplateFunc(t *template.Template) map[string]interface{} {
	f := sprig.TxtFuncMap()

	f["toYaml"] = func(v interface{}) string {
		data, err := yaml.Marshal(v)
		if err != nil {
			return ""
		}

		return string(data)
	}

	f["randPort"] = func() int {
		rand.Seed(time.Now().UnixNano())

		min := 10000
		max := 50000

		randPort := rand.Intn(max-min) + min //nolint:gosec

		return randPort
	}

	return f
}
