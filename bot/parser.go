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

		// language=regexp
		return regexp.MustCompile(fmt.Sprintf("(?i)^%s(-)?([0-9]+)?", tr)), nil
	case []interface{}:
		// Слайсы (списки) читается как слайс
		// с interface{}, поэтому придётся
		// проверять каждый объект

		var keyWords = make([]string, 0, len(tr))
		for _, v := range tr {
			switch k := v.(type) {
			case string:
				keyWords = append(keyWords, k)
			case rune:
				keyWords = append(keyWords, string(k))
			default:
				return nil, nil
			}
		}
		// language=regexp
		return regexp.MustCompile(fmt.Sprintf("(?i)^(?:%s)(-)?([0-9]+)?", strings.Join(keyWords, "|"))), nil
	default:
		return nil, nil
	}
}
