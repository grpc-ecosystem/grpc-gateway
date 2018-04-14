// +build !go1.8

// Copyright (c) 2015-2018 Jeevanandam M (jeeva@myjeeva.com)
// 2016 Andrew Grigorev (https://github.com/ei-grad)
// All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import "strings"

func errIsContextCanceled(err error) bool {
	return strings.Contains(err.Error(), "request canceled")
}
