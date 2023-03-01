/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cri

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"sort"
	"strconv"
	"strings"
)

// List returns the contents as a sorted T slice.
//
// This is a separate function and not a method because not all types supported
// by Generic are ordered and only those can be sorted.
func List(s sets.Set[string]) []string {
	res := make(sortableSliceOfGeneric, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return res
}

type sortableSliceOfGeneric []string

func (g sortableSliceOfGeneric) Len() int { return len(g) }
func (g sortableSliceOfGeneric) Less(i, j int) bool {
	left := g[i]
	right := g[j]

	lstr := strings.Split(string(left), ".")[2]
	rstr := strings.Split(string(right), ".")[2]
	lint, _ := strconv.Atoi(lstr)
	rint, _ := strconv.Atoi(rstr)
	return lint < rint
}
func (g sortableSliceOfGeneric) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
