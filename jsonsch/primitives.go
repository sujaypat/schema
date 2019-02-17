package jsonsch

import (
	"fmt"
	"reflect"
)

type Type string

const (
	Object  Type = "object"
	Boolean Type = "boolean"
	Array   Type = "array"
	Number  Type = "number"
	Integer Type = "integer"
	String  Type = "string"
	Null    Type = "null"
)

type Primitive struct {
	Type        Type   `json:"type"`
	Description string `json:"description,omitempty"`
}

func NewNull(params *FromExampleParams) Primitive {
	if params.NullAs == "" {
		return Primitive{Type: Null}
	}

	switch params.NullAs {
	case "null", "nil":
		return Primitive{Type: Null}
	case "bool":
		return Primitive{Type: Boolean}
	case "string":
		return Primitive{Type: String}
	case "number", "float":
		return Primitive{Type: Number}
	case "object":
		return Primitive{Type: Object}
	default:
		return Primitive{Type: Null}
	}
}

func NewBoolean() Primitive {
	return Primitive{Type: Boolean}
}

func NewNumber() Primitive {
	return Primitive{Type: Number}
}

func NewString() Primitive {
	return Primitive{Type: String}
}

type ArraySchema struct {
	Type  Type        `json:"type"`
	Items interface{} `json:"items"`
}

func NewArray(items interface{}) ArraySchema {
	return ArraySchema{Type: Array, Items: items}
}

func InferArrayTypeFromElem(data []interface{}, params *FromExampleParams) (ArraySchema, error) {
	var elem interface{}

	if len(data) == 0 {
		if params.EmptyArraysAs == "" {
			return ArraySchema{}, fmt.Errorf("cannot infer type of empty array; consider using --empty-arrays-as")
		}
		switch params.EmptyArraysAs {
		case "null", "nil":
			elem = nil
		case "bool", "boolean":
			elem = false
		case "string", "str":
			elem = ""
		case "number", "float":
			elem = 0.0
		case "object":
			elem = make(map[string]interface{})
		default:
			return ArraySchema{}, fmt.Errorf("invalid --empty-arrays-as value '%v'", params.EmptyArraysAs)
		}
	} else {
		elem = data[0]
		for _, checkElem := range data {
			if reflect.TypeOf(checkElem) != reflect.TypeOf(elem) {
				return ArraySchema{}, fmt.Errorf("mismatched types %T, %T", elem, checkElem)
			}
		}
	}

	a := ArraySchema{Type: Array}

	if err := buildSchema(elem, &a.Items, params); err != nil {
		return ArraySchema{}, err
	}

	return a, nil
}
