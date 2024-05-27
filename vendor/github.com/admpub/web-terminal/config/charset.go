package config

import (
	"io"
	"sort"
	"strings"

	"github.com/admpub/log"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var CharsetList = map[string]encoding.Encoding{
	`GB18030`:       simplifiedchinese.GB18030,
	`GB2312`:        simplifiedchinese.HZGB2312,
	`HZ-GB2312`:     simplifiedchinese.HZGB2312,
	`GBK`:           simplifiedchinese.GBK,
	`BIG5`:          traditionalchinese.Big5,
	`EUC-JP`:        japanese.EUCJP,
	`ISO2022JP`:     japanese.ISO2022JP,
	`SHIFTJIS`:      japanese.ShiftJIS,
	`EUC-KR`:        korean.EUCKR,
	`UTF-8`:         encoding.Nop,
	`UTF-16`:        unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	`UTF-16-BE`:     unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	`UTF-16-LE`:     unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
	`UTF-16-BOM`:    unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
	`UTF-16-BE-BOM`: unicode.UTF16(unicode.BigEndian, unicode.UseBOM),
	`UTF-16-LE-BOM`: unicode.UTF16(unicode.LittleEndian, unicode.UseBOM),
}

func SupportedCharsets() []string {
	r := []string{}
	for k := range CharsetList {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

func CharsetEncoding(charset string) encoding.Encoding {
	charset = strings.ToUpper(charset)
	switch charset {
	case "UTF8":
		charset = "UTF-8"
	case "UTF16-BOM":
		charset = "UTF-16-BOM"
	case "UTF16-BE-BOM":
		charset = "UTF-16-BE-BOM"
	case "UTF16-LE-BOM":
		charset = "UTF-16-LE-BOM"
	case "UTF16":
		charset = "UTF-16"
	case "UTF16-BE":
		charset = "UTF-16-BE"
	case "UTF16-LE":
		charset = "UTF-16-LE"
	case "UTF32":
		charset = "UTF-32"
	}
	if enc, ok := CharsetList[charset]; ok {
		return enc
	}
	return nil
}

func DecodeBy(charset string, dst io.Writer) io.Writer {
	cs := CharsetEncoding(charset)
	if nil == cs {
		log.Warnf("[web-terminal]charset %q is not exists.", charset)
		return dst
	}

	return transform.NewWriter(dst, cs.NewDecoder())
}
