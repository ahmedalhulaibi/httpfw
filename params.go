package httpfw

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type httpRequestContext interface {
	context.Context
}

func GetQueryParamInt(params url.Values, name string) (int, error) {
	var param int
	if paramStr := params.Get(name); paramStr != "" {
		paramInt, err := strconv.Atoi(paramStr)
		if err != nil {

			return param, fmt.Errorf("%s=%v not valid", name, paramStr)
		}

		param = paramInt
	}

	return param, nil
}

func GetQueryParam(params url.Values, name string) string {
	return params.Get(name)
}
