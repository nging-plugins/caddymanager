package form

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsNginxLogFormat(t *testing.T) {
	value := `{combined}`
	logFormat := AsNginxLogFormat(value)
	assert.Equal(t, `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`, logFormat)
}

/**
Nginx中location匹配：
= 表示精确匹配
^~ 表示uri以某个常规字符串开头,大多情况下用来匹配url路径，nginx不对url做编码，因此请求为/static/20%/aa，可以被规则 “^~ /static/ /aa” 匹配到（注意是空格，即所见即所得）。
~ 正则匹配(区分大小写)
~* 正则匹配(不区分大小写)
/ 任何请求都会匹配

Nginx 通配符优先级：
首先匹配 =，其次匹配 ^~, 其次是按文件中顺序的正则匹配，最后是交给 / 通用匹配。当有匹配成功时候，停止匹配，按当前匹配规则处理请求。
*/
