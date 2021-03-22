/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

type Any = interface{}

type Filler = func(Any) error

type FillerError struct {
	FillerError error
}

func (ferr FillerError) Error() string {
	return ferr.FillerError.Error()
}

type ErrorResult struct {
	StatusCode       int
	StatusPhrase     string
	DeveloperMessage string
	ErrorString      string
	TraceId          string
	ParentSpanId     string
	SpanId           string
	Cause            map[string]string
	Args             []Any
	Details          []ErrorDetail
}

type ErrorDetail struct {
	error string
	Args  []string
	Path  string
}

type RequestContext struct {
}
