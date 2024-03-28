// Copyright (c) HashiCorp, Inc.

package param

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

type QueryParam struct {
	url.Values
}

func NewQueryParam() QueryParam {
	return QueryParam{
		Values: make(url.Values),
	}
}

// 查询条件，查询API的查询条件类似于MySQL数据库。
// 省略该字段将返回所有记录，返回记录数的上限受限于limit字段
func (params *QueryParam) AddQ(q string) *QueryParam {
	if params.Get("q") == "" {
		params.Set("q", q)
	} else {
		params.Add("q", q)
	}
	return params
}

// 最多返回的记录数，类似MySQL的limit，默认值1000
func (params *QueryParam) Limit(limit int) *QueryParam {
	params.Set("limit", fmt.Sprintf("%d", limit))
	return params
}

// 起始查询记录位置，类似MySQL的offset。跟limit配合使用可以实现分页
func (params *QueryParam) Start(start int) *QueryParam {
	params.Set("start", fmt.Sprintf("%d", start))
	return params
}

// 计数查询，相当于MySQL中的count()函数。当设置成true时，API只返回的是满足查询条件的记录数
func (params *QueryParam) Count(count bool) *QueryParam {
	params.Set("count", fmt.Sprintf("%t", count))
	return params
}

// 以字段分组，相当于MySQL中的group by关键字。例如groupBy=type
func (params *QueryParam) GroupBy(groupBy string) *QueryParam {
	params.Set("groupBy", groupBy)
	return params
}

// replyWithCount被设置成true后，查询返回中会包含满足查询条件的记录总数，跟start值比较就可以得知还需几次分页。
func (params *QueryParam) ReplyWithCount(replyWithCount bool) *QueryParam {
	params.Set("replyWithCount", fmt.Sprintf("%t", replyWithCount))
	return params
}

// 未知，来自ZStack Java SDK【sdk-4.4.0.jar】
func (params *QueryParam) FilterName(filterName string) *QueryParam {
	params.Set("filterName", filterName)
	return params
}

// 以字段排序，等同于MySQL中的sort by关键字。必须跟+或者-配合使用，+表示升序，-表示降序，后面跟排序字段名，例如：
// sort=+key，根据key进行升序排序
// sort=-key，根据 key进行降序排序
func (params *QueryParam) Sort(sort string) *QueryParam {
	params.Set("sort", sort)
	return params
}

// 指定返回的字段，等同于MySQL中的select字段功能。例如fields=name,uuid，则只返回满足条件记录的name和uuid字段
func (params *QueryParam) Fields(fields []string) *QueryParam {
	params.Set("fields", strings.Join(fields, ","))
	return params
}

// ConvertStruct2UrlValues param should be
func ConvertStruct2UrlValues(param interface{}) (url.Values, error) {
	if reflect.Ptr != reflect.TypeOf(param).Kind() {
		return nil, errors.New("model should be pointer kind")
	}
	result := url.Values{}
	if param == nil || reflect.ValueOf(param).IsNil() {
		return nil, errors.New("param is nil")
	}

	s := structs.New(param)
	s.TagName = "json"
	mappedOpts := s.Map()
	for k, v := range mappedOpts {
		result.Set(k, fmt.Sprintf("%v", v))
	}
	return result, nil
}
