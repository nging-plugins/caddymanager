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
	"sort"
	"strings"

	"github.com/admpub/caddy/dnsproviders"
	"github.com/admpub/log"
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

func (v Values) getExtraIndexes(addon string) []string {
	return v.Values[addon+`_extra_index[]`]
}

func (v Values) GetExtraList(addon string) []*ExtraItem {
	indexes := v.getExtraIndexes(addon)
	results := make([]*ExtraItem, len(indexes)+1)
	results[0] = &ExtraItem{
		addon:  addon,
		Values: v,
	}
	for i, index := range indexes {
		results[i+1] = &ExtraItem{
			addon:  addon,
			prefix: addon + `_extra`,
			index:  index,
			Values: v,
		}
	}
	return results
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

func (v Values) SplitBySpaceWithPrefixAndSuffix(content string, prefix string, suffix string) []string {
	return SplitBySpace(content, func(v string) string {
		return prefix + v + suffix
	})
}

func (v Values) SplitBySpaceWithRegexpQuote(content string) []string {
	return SplitBySpace(content, v.RegexpQuote)
}

func (v Values) SplitBySpaceWithPathWildcardSuffix(content string) []string {
	return SplitBySpace(content, v.AddPathWildcardSuffix)
}

func (v Values) SplitBySpaceWithExtWildcardPrefix(content string) []string {
	return SplitBySpace(content, v.AddExtWildcardPrefix)
}

func (v Values) RegexpQuote(content string) string {
	return regexp.QuoteMeta(content)
}

func (v Values) SliceAddPathWildcardSuffix(content []string) []string {
	for index, value := range content {
		content[index] = v.AddPathWildcardSuffix(value)
	}
	return content
}

func (v Values) SliceAddExtWildcardPrefix(content []string) []string {
	for index, value := range content {
		content[index] = v.AddExtWildcardPrefix(value)
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
	values := v.Values[key]
	return param.StringSlice(values)
}

func (v Values) AddSlashes(val string) string {
	return AddCSlashesIngoreSlash(val, '"')
}

func (v Values) AddSlashesSingleQuote(val string) string {
	return AddCSlashesIngoreSlash(val, '\'')
}

func (v Values) IteratorKV(addon string, item string, prefix string, withQuotes ...bool) interface{} {
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys := v.Values[k]

	k = addon + item + `_v`
	values := v.Values[k]
	return v.iteratorKV(keys, values, prefix, withQuotes...)
}

func (v Values) iteratorKV(keys []string, values []string, prefix string, withQuotes ...bool) interface{} {
	var r, t string
	var withQuote bool
	var ignoreSlash bool
	if len(withQuotes) > 0 {
		withQuote = withQuotes[0]
		if len(withQuotes) > 1 {
			ignoreSlash = withQuotes[1]
		}
	}
	l := len(values)
	var suffix string
	if v.Config.GetEngine() == `nginx` {
		suffix = `;`
	}
	for i, k := range keys {
		if len(k) == 0 {
			continue
		}
		if i < l {
			val := values[i]
			if withQuote {
				if val == `"` || (!strings.HasPrefix(val, `"`) && !strings.HasSuffix(val, `"`)) {
					if ignoreSlash {
						val = `"` + v.AddSlashes(val) + `"`
					} else {
						val = `"` + com.AddCSlashes(val, '"') + `"`
					}
				}
			}
			r += t + prefix + k + `   ` + val + suffix
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
	keys := v.Values[k]

	k = addon + item + `_v`
	values := v.Values[k]

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

func (v Values) GetKVData(addon string, itemOr ...string) *echo.KVData {
	var item string
	if len(itemOr) > 0 {
		item = itemOr[0]
	}
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys := v.Values[k]

	k = addon + item + `_v`
	values := v.Values[k]

	l := len(values)
	result := echo.NewKVData()
	for i, k := range keys {
		if len(k) == 0 {
			continue
		}
		if i < l {
			result.Add(k, values[i])
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
	values := v.Values[k]
	return v.getValueList(values)
}

func (v Values) getValueList(values []string) []string {
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
	values := v.Values[k]
	return v.iterator(values, prefix, withQuotes...)
}

func (v Values) iterator(values []string, prefix string, withQuotes ...bool) interface{} {
	var r, t string
	var withQuote bool
	var ignoreSlash bool
	if len(withQuotes) > 0 {
		withQuote = withQuotes[0]
		if len(withQuotes) > 1 {
			ignoreSlash = withQuotes[1]
		}
	}
	for _, val := range values {
		if len(val) == 0 {
			continue
		}
		if withQuote {
			if ignoreSlash {
				val = `"` + v.AddSlashes(val) + `"`
			} else {
				val = `"` + com.AddCSlashes(val, '"') + `"`
			}
		}
		r += t + prefix + val
		t = "\n"
	}
	if withQuote {
		return template.HTML(r)
	}
	return r
}

func (v Values) GroupByLocations(fields []string) Locations {
	var staticPathList []string
	var regexpPathList []string
	groupByPath := map[string][]*LocationDef{}
	var addPath = func(pathKey, moduleName, path string, extraIndex int, extraItem *ExtraItem, isRegexp bool) {
		data := &LocationDef{
			PathKey:    pathKey,
			Module:     moduleName,
			Location:   path,
			ExtraIndex: extraIndex,
			ExtraItem:  extraItem,
		}
		if _, ok := groupByPath[path]; !ok {
			groupByPath[path] = []*LocationDef{}
			if isRegexp {
				regexpPathList = append(regexpPathList, path)
			} else {
				staticPathList = append(staticPathList, path)
			}
		}
		groupByPath[path] = append(groupByPath[path], data)
	}
	var addPath2 = func(pathKey, moduleName, path string, extraIndex int, extraItem *ExtraItem) {
		if _, ok := groupByPath[path]; !ok {
			groupByPath[path] = []*LocationDef{}
			staticPathList = append(staticPathList, path)
		}
		data := &LocationDef{
			PathKey:    pathKey,
			Module:     moduleName,
			Location:   path,
			ExtraIndex: extraIndex,
			ExtraItem:  extraItem,
		}
		groupByPath[path] = append(groupByPath[path], data)
	}
	for _, pathKey := range fields {
		parts := strings.SplitN(pathKey, `_`, 2)
		moduleName := parts[0]
		if !v.IsEnabled(moduleName) {
			continue
		}
		if len(parts) == 2 {
			pathKey = parts[1]
		}
		if strings.HasSuffix(pathKey, `[]`) {
			pathKey = strings.TrimSuffix(pathKey, `[]`)
			var isRegexp bool
			if moduleName+`_`+pathKey == `expires_match_k` {
				isRegexp = true
			}
			extraList := v.GetExtraList(moduleName)
			for index, extraItem := range extraList {
				values := extraItem.GetValues(pathKey)
				if len(values) == 1 {
					addPath(extraItem.GenName(pathKey), moduleName, values[0], index, extraItem, isRegexp)
				}
			}
		} else {
			extraList := v.GetExtraList(moduleName)
			for index, extraItem := range extraList {
				path := extraItem.Get(pathKey)
				addPath2(extraItem.GenName(pathKey), moduleName, path, index, extraItem)
			}
		}
	}
	sort.Sort(SortByLen(regexpPathList))
	sort.Sort(SortByLen(staticPathList))
	return Locations{
		SortedStaticPath: staticPathList,
		SortedRegexpPath: regexpPathList,
		GroupByPath:      groupByPath,
	}
}

func (v Values) ParseHost(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Error(err)
		return ``
	}
	return u.Host
}

func (v Values) DNSProvider(provider string) []string {
	inputs := dnsproviders.GetInputs(provider).Clone()
	for i, input := range inputs {
		inputs[i].Value = v.Values.Get(`tls_acme_dns_` + input.Name)
	}
	return inputs.RenderCaddyfile()
}

func (v Values) CA(name string) *echo.KV {
	return CAList.GetItem(name)
}
