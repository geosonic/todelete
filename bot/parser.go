/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"fmt"
	"regexp"
	"strings"
)

func parse(triggerWord interface{}) (*regexp.Regexp, error) {
	switch tr := triggerWord.(type) {
	case string:
		// Если указана строка, то сразу компилируем
		// регулярку и возвращаем объект, а ещё
		// это работает как обратная совместимость

		return regexp.MustCompile(fmt.Sprintf("^%v(-)?([0-9]+)?", strings.ToLower(tr))), nil
	case []interface{}:
		// Слайсы (списки) читается как слайс
		// с interface{}, поэтому придётся
		// проверять каждый объект

		var keyWords = make([]string, 0, len(tr))
		for _, v := range tr {
			switch k := v.(type) {
			case string:
				keyWords = append(keyWords, strings.ToLower(k))
			case rune:
				keyWords = append(keyWords, strings.ToLower(string(k)))
			default:
				return nil, nil
			}
		}
		return regexp.MustCompile(fmt.Sprintf("^(?:%v)(-)?([0-9]+)?", strings.Join(keyWords, "|"))), nil
	default:
		return nil, nil
	}
}
