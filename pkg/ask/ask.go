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
package ask

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/maksim-paskal/auto-devops/pkg/config"
	"github.com/maksim-paskal/auto-devops/pkg/types"
	"github.com/pkg/errors"
)

//nolint: gochecknoglobals
var (
	Reader           = bufio.NewReader(os.Stdin)
	Output io.Writer = os.Stdout
)

func printf(format string, args ...interface{}) error {
	line := fmt.Sprintf(format, args...)

	if _, err := Output.Write([]byte(line)); err != nil {
		return errors.Wrap(err, "error writing")
	}

	return nil
}

func Once(q *types.BoostrapQuestion) (bool, error) { //nolint:cyclop
	var err error

	q.Result, err = config.Boostrap.Template(q.Result)
	if err != nil {
		return false, errors.Wrap(err, "error templating")
	}

	if len(q.Validation) > 0 {
		if err = printf("%s %s\n", q.Prompt, q.Validation); err != nil {
			return false, errors.Wrap(err, "error printing")
		}
	} else {
		if err = printf("%s\n", q.Prompt); err != nil {
			return false, errors.Wrap(err, "error printing")
		}
	}

	if len(q.Result) > 0 {
		if err = printf("[%s]: ", q.Result); err != nil {
			return false, errors.Wrap(err, "error printing")
		}
	}

	data, _, err := Reader.ReadLine()
	if err != nil {
		return false, errors.Wrap(err, "error reading line")
	}

	newResult := string(data)
	newResult = strings.TrimSpace(newResult)

	if len(newResult) == 0 {
		newResult = q.Result
	}

	if len(q.Validation) > 0 {
		matched, err := regexp.MatchString(q.Validation, newResult)
		if err != nil {
			return false, errors.Wrap(err, "error validation")
		}

		if matched {
			q.Result = newResult
		}

		return matched, nil
	}

	q.Result = newResult

	return true, nil
}
