package convert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
)

// Parameter converts parameter's value(s) according to parameter's type
// and format. Type and format MUST match OAS 2.0.
// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#parameterObject
func Parameter(vals []string, param *spec.Parameter) (value interface{}, err error) {
	if param.Type == "array" {
		if param.Items == nil {
			return nil, fmt.Errorf("type array has no `items` field")
		}

		// Array can be represented in different ways based on
		// collectionFormat property
		if len(vals) != 0 {
			switch param.SimpleSchema.CollectionFormat {
			case "ssv":
				// Space-separated values
				vals = strings.Split(vals[0], " ")
			case "tsv":
				// Tab-separated values
				vals = strings.Split(vals[0], "\t")
			case "pipes":
				// Pipe-separated values
				vals = strings.Split(vals[0], "|")
			case "multi":
				// Multiple parameter instances rather than multiple values
				// Do nothing, values are already represented as an array in vals
			default: // "csv"
				// Comma-separated values
				vals = strings.Split(vals[0], ",")
			}
		}

		return Array(vals, param.Items.Type, param.Items.Format)
	}

	if param.Type == "file" {
		// TODO
		return nil, fmt.Errorf("type %s: NOT IMPLEMENTED", param.Type)
	}

	if len(vals) != 1 {
		return nil, fmt.Errorf(
			"values count is %d, want 1",
			len(vals),
		)
	}

	return Primitive(vals[0], param.Type, param.Format)
}

// Primitive converts string values according to type and format described
// in OAS 2.0.
// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#parameterObject
func Primitive(val string, typ, format string) (value interface{}, err error) {
	switch typ {
	case "string":
		return convertString(val, format)
	case "number":
		return convertNumber(val, format)
	case "integer":
		return convertInteger(val, format)
	case "boolean":
		return convertBoolean(val)
	default:
		return nil, fmt.Errorf(
			"unknown type: %s",
			typ,
		)
	}
}

var evaluatesAsTrue = map[string]struct{}{
	"true":     struct{}{},
	"1":        struct{}{},
	"yes":      struct{}{},
	"ok":       struct{}{},
	"y":        struct{}{},
	"on":       struct{}{},
	"selected": struct{}{},
	"checked":  struct{}{},
	"t":        struct{}{},
	"enabled":  struct{}{},
}

var evaluatesAsFalse = map[string]struct{}{
	"false":    struct{}{},
	"0":        struct{}{},
	"no":       struct{}{},
	"n":        struct{}{},
	"off":      struct{}{},
	"f":        struct{}{},
	"disabled": struct{}{},
}

func convertString(val, format string) (interface{}, error) {
	switch format {
	case "":
		return val, nil
	case "partial-time", "uuid":
		// For now, return as-is.
		return val, nil
	default:
		// TODO: parse formats byte, binary, date, date-time
		return nil, fmt.Errorf(
			"unknown format %s for type string",
			format,
		)
	}
}

func convertInteger(val, format string) (interface{}, error) {
	switch format {
	case "int32":
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %v to int32", val)
		}
		return int32(i), nil
	case "int64":
		fallthrough
	case "":
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %v to int64", val)
		}
		return i, nil
	default:
		return nil, fmt.Errorf(
			"unknown format %s for type integer",
			format,
		)
	}
}

func convertNumber(val, format string) (interface{}, error) {
	switch format {
	case "float":
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %v to float", val)
		}
		return float32(f), nil
	case "double":
		fallthrough
	case "":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %v to double", val)
		}
		return f, nil
	default:
		return nil, fmt.Errorf(
			"unknown format %s for type integer",
			format,
		)
	}
}

func convertBoolean(val string) (interface{}, error) {
	v := strings.ToLower(val)
	if _, ok := evaluatesAsTrue[v]; ok {
		return true, nil
	}
	if _, ok := evaluatesAsFalse[v]; ok {
		return false, nil
	}
	return false, fmt.Errorf("unknown format %s for type boolean", val)
}
