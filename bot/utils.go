/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

var m = minify.New()

func init() {
	m.Add("text/js", &js.Minifier{})
}

// Сжатие JavaScript кода
func CompressJS(code string) (string, error) {
	t, err := m.String("text/js", code)

	return t, err
}
