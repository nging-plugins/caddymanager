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

package logcounter

type Counter struct {
	Views       map[string]uint64  //每个网址的浏览器
	Elapsed     map[string]float64 //每个网址的总耗时(秒)
	MaxElapsed  map[string]float64 //每个网址的最大耗时(秒)
	BodyBytes   map[string]uint64  //每个网址的总字节数
	IPs         map[string]uint64  //每个网址的每分钟的IP数量
	UserAgents  map[string][]int64
	StatusCodes map[string][]int64

	//200/301/302/400/403/404/499/500/502/503/504
	StatusCode map[int]uint64    //状态码次数
	Browsers   map[string]uint64 //浏览器访问次数
}

type URLCounter struct {
	URL      string  // 网址
	Views    uint64  // 浏览量
	Elapsed  float64 // 耗时(秒)
	BodySize uint64  // 响应Body尺寸(字节)
	IPCount  uint64  // IP数量
	Code2xx  uint64  // 状态码 2xx 数量
	Code3xx  uint64  // 状态码 3xx 数量
	Code4xx  uint64  // 状态码 4xx 数量
	Code5xx  uint64  // 状态码 5xx 数量
}

type CodeCounter struct {
	Code  uint   // 状态码
	Count uint64 // 数量
}

type IPCounter struct {
	IP     string // IP
	Region string // 地区
	Views  uint64 // 浏览量
}

type RegionCounter struct {
	Region  string // 地区
	IPCount uint64 // IP数量
}
