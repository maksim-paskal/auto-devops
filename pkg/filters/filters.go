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
package filters

import (
	"regexp"

	"github.com/maksim-paskal/auto-devops/pkg/config"
	"github.com/maksim-paskal/auto-devops/pkg/types"
	"github.com/pkg/errors"
)

func GetFilter(path string) (types.BoostrapFilter, error) {
	filterRule := types.BoostrapFilter{}

	for _, filter := range config.Boostrap.Filters {
		matched, err := regexp.MatchString(filter.Match, path)
		if err != nil {
			return types.BoostrapFilter{}, errors.Wrap(err, "error in matching")
		}

		if !matched {
			continue
		}

		if len(filter.Condition) > 0 {
			conditionResults, err := config.Boostrap.Template(filter.Condition)
			if err != nil {
				return types.BoostrapFilter{}, errors.Wrap(err, "error templating")
			}

			if conditionResults == "true" {
				filterRule = filter
			}
		} else {
			filterRule = filter
		}

		break
	}

	return filterRule, nil
}
