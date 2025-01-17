package form

import (
	"strings"
)

type ExtraItem struct {
	addon  string
	prefix string
	index  string
	Values
}

func (t *ExtraItem) GenName(name string) string {
	if len(t.prefix) == 0 {
		if len(t.addon) == 0 {
			return name
		}
		name = t.addon + `_` + name
		return name
	}
	if len(t.addon) > 0 {
		name = t.addon + `_` + name
	}
	name = t.FullPrefix() + GenNameSuffix(name)
	return name
}

func (t *ExtraItem) FullPrefix() string {
	if len(t.prefix) == 0 {
		return ``
	}
	return t.prefix + `[` + t.index + `]`
}

func GenNameSuffix(name string) string {
	if strings.HasSuffix(name, `]`) {
		pos := strings.Index(name, `[`)
		if pos > 0 {
			name = `[` + name[0:pos] + `]` + name[pos:]
		} else if pos != 0 {
			name = `[` + name
		}
	} else {
		name = `[` + name + `]`
	}
	return name
}

func (t *ExtraItem) Get(name string) string {
	name = t.GenName(name)
	return t.Values.Get(name)
}

func (t *ExtraItem) GetValues(name string) []string {
	name = t.GenName(name)
	return t.Values.Values[name]
}

func (t *ExtraItem) GetAttrVal(item string, defaults ...string) string {
	val := t.Get(item)
	if len(val) == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ``
	}
	return val
}

func (t *ExtraItem) GetValueList(item string) []string {
	k := t.GenName(item)
	values := t.Values.Values[k]
	return t.Values.getValueList(values)
}

func (t *ExtraItem) AddonAttrFromKey(key string, item string, defaults ...string) string {
	val := t.Get(key)
	if len(val) == 0 {
		if len(defaults) > 0 && len(defaults[0]) > 0 {
			return item + `   ` + defaults[0]
		}
		return ``
	}
	return item + `   ` + val
}

func (t *ExtraItem) AddonAttr(item string, defaults ...string) string {
	val := t.Get(item)
	if len(val) == 0 {
		if len(defaults) > 0 && len(defaults[0]) > 0 {
			return item + `   ` + defaults[0]
		}
		return ``
	}
	return item + `   ` + val
}

func (t *ExtraItem) IteratorKV(item string, prefix string, withQuotes ...bool) interface{} {
	k := t.GenName(item + `_k`)
	keys := t.Values.Values[k]

	k = t.GenName(item + `_v`)
	values := t.Values.Values[k]
	return t.Values.iteratorKV(keys, values, prefix, withQuotes...)
}

func (t *ExtraItem) Iterator(item string, prefix string, withQuotes ...bool) interface{} {
	k := t.GenName(item)
	values := t.Values.Values[k]
	return t.Values.iterator(values, prefix, withQuotes...)
}

func (t *ExtraItem) ServerGroup(key string, customHost string, withQuotes ...bool) interface{} {
	val := t.Get(key)
	prefix := t.FullPrefix()
	if len(prefix) > 0 {
		return t.Values.serverGroup(val, customHost, func(name string) string {
			return t.Values.Get(prefix + GenNameSuffix(name))
		}, withQuotes...)
	}
	return t.Values.serverGroup(val, customHost, t.Values.Get, withQuotes...)
}

func (t *ExtraItem) IteratorHeaderKV(item string, plusPrefix string, minusPrefix string, withValueAndQuotes ...bool) interface{} {
	k := t.GenName(item + `_k`)
	keys := t.Values.Values[k]

	k = t.GenName(item + `_v`)
	values := t.Values.Values[k]
	return t.Values.iteratorHeaderKV(keys, values, plusPrefix, minusPrefix, withValueAndQuotes...)
}

func (t *ExtraItem) IteratorNginxProxyHeaderKV() interface{} {
	item := `header_downstream`
	k := t.GenName(item + `_k`)
	keys := t.Values.Values[k]

	k = t.GenName(item + `_v`)
	values := t.Values.Values[k]
	return t.Values.iteratorNginxProxyHeaderKV(keys, values)
}
