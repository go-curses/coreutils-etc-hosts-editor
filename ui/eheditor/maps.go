// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eheditor

import (
	"sort"

	"github.com/maruel/natural"
)

func SortedKeys[V interface{}](data map[string]V) (keys []string) {
	for key, _ := range data {
		keys = append(keys, key)
	}
	sort.Sort(natural.StringSlice(keys))
	return
}

func StringInSlices(src string, dst ...[]string) bool {
	for _, v := range dst {
		for _, vv := range v {
			if src == vv {
				return true
			}
		}
	}
	return false
}