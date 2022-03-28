package oapiutil

import (
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"strings"
)

func ValidationError(err error) error {

	var sb strings.Builder
	if reqErr, ok := err.(*openapi3filter.RequestError); ok {

		if reqErr.Reason != "" {
			sb.WriteString(reqErr.Reason)
			sb.WriteString(" - ")
			if schErr, ok := reqErr.Err.(*openapi3.SchemaError); ok {
				sb.WriteString(schErr.Reason)
			}

			return errors.New(sb.String())
		}
	}

	return err
}
