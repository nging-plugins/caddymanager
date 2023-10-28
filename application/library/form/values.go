/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package form

import (
	"html/template"
	"net/url"
	"regexp"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/library/webdav"
)

func NewValues(values url.Values, cfg engine.Configer) *Values {
	return &Values{Values: values, Config: cfg}
}

type Values struct {
	url.Values
	Config engine.Configer
}

func (v Values) GetWebdavGlobal() []*webdav.WebdavPerm {
	return webdav.ParseGlobalForm(v.Values)
}

func (v Values) VhostConfigDir() string {
	if v.Config.GetEnviron() == engine.EnvironContainer {
		return v.Config.GetVhostConfigContainerDir()
	}
	dir, _ := v.Config.GetVhostConfigLocalDirAbs()
	return dir
}

func (v Values) GetWebdavUser() []*webdav.WebdavUser {
	return webdav.ParseUserForm(v.Values)
}

func (v Values) GetDomainList() []string {
	domain := v.Values.Get(`domain`)
	return SplitBySpace(domain)
}

func (v Values) SplitBySpace(content string) []string {
	return SplitBySpace(content)
}

func (v Values) SplitBySpaceWithRegexpQuote(content string) []string {
	return SplitBySpace(content, v.RegexpQuote)
}

func (v Values) SplitBySpaceWithWildcardSuffix(content string) []string {
	return SplitBySpace(content, v.AddWildcardSuffix)
}

func (v Values) RegexpQuote(content string) string {
	return regexp.QuoteMeta(content)
}

func (v Values) SliceAddWildcardSuffix(content []string) []string {
	for index, value := range content {
		content[index] = v.AddWildcardSuffix(value)
	}
	return content
}

func (v Values) SliceRegexpQuote(content []string) []string {
	for index, value := range content {
		content[index] = v.RegexpQuote(value)
	}
	return content
}

func (v Values) GetSlice(key string) param.StringSlice {
	values, _ := v.Values[key]
	return param.StringSlice(values)
}

func (v Values) AddSlashes(val string) string {
	return com.AddCSlashes(val, '"')
}

func (v Values) AddSlashesSingleQuote(val string) string {
	return com.AddCSlashes(val, '\'')
}

func (v Values) IteratorKV(addon string, item string, prefix string, withQuotes ...bool) interface{} {
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys, _ := v.Values[k]

	k = addon + item + `_v`
	values, _ := v.Values[k]

	var r, t string
	var withQuote bool
	if len(withQuotes) > 0 {
		withQuote = withQuotes[0]
	}
	l := len(values)
	var suffix string
	if v.Config.GetEngine() == `nginx` {
		suffix = `;`
	}
	for i, k := range keys {
		if i < l {
			v := values[i]
			if withQuote {
				v = `"` + com.AddCSlashes(v, '"') + `"`
			}
			r += t + prefix + k + `   ` + v + suffix
			t = "\n"
		}
	}
	if withQuote {
		return template.HTML(r)
	}
	return r
}

func (v Values) GetAttrVal(addon string, item string, defaults ...string) string {
	if len(addon) > 0 {
		addon += `_`
	}
	k := addon + item
	val := v.Values.Get(k)
	if len(val) == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ``
	}
	return val
}

func (v Values) AddonAttr(addon string, item string, defaults ...string) string {
	if len(addon) > 0 {
		addon += `_`
	}
	k := addon + item
	val := v.Values.Get(k)
	if len(val) == 0 {
		if len(defaults) > 0 && len(defaults[0]) > 0 {
			return item + `   ` + defaults[0]
		}
		return ``
	}
	return item + `   ` + val
}

func (v Values) IsEnabled(key string, expectedValue ...string) bool {
	var expected string
	if len(expectedValue) > 0 {
		expected = expectedValue[0]
	} else {
		expected = `1`
	}
	return v.Get(key) == expected
}

func (v Values) GetKVList(addon string, itemOr ...string) []echo.KV {
	var item string
	if len(itemOr) > 0 {
		item = itemOr[0]
	}
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys, _ := v.Values[k]

	k = addon + item + `_v`
	values, _ := v.Values[k]

	l := len(values)
	result := make([]echo.KV, 0, len(keys))
	for i, k := range keys {
		if len(k) == 0 {
			continue
		}
		if i < l {
			result = append(result, echo.KV{K: k, V: values[i]})
		}
	}
	return result
}

func (v Values) GetValueList(addon string, itemOr ...string) []string {
	var item string
	if len(itemOr) > 0 {
		item = itemOr[0]
	}
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item
	values, _ := v.Values[k]
	result := make([]string, 0, len(values))
	for _, v := range values {
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result
}

func (v Values) Iterator(addon string, item string, prefix string, withQuotes ...bool) interface{} {
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item
	values, _ := v.Values[k]
	var r, t string
	var withQuote bool
	if len(withQuotes) > 0 {
		withQuote = withQuotes[0]
	}
	for _, v := range values {
		if withQuote {
			v = `"` + com.AddCSlashes(v, '"') + `"`
		}
		r += t + prefix + v
		t = "\n"
	}
	if withQuote {
		return template.HTML(r)
	}
	return r
}
