/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

import "github.com/pvillela/go-foa-realworld/internal/arch"

type ErrorResult struct {
	StatusCode       int
	StatusPhrase     string
	DeveloperMessage string
	ErrorString      string
	TraceId          string
	ParentSpanId     string
	SpanId           string
	Cause            map[string]string
	Args             []arch.Any
	Details          []ErrorDetail
}

type ErrorDetail struct {
	Err  string
	Args []string
	Path string
}

type RequestContext struct {
	Username string
}
