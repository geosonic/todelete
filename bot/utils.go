/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"sync"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

var m = minify.New()

var once sync.Once

// Сжатие JavaScript кода
func CompressJS(code string) (string, error) {
	once.Do(func() {
		m.Add("text/js", &js.Minifier{})
	})

	return m.String("text/js", code)
}
