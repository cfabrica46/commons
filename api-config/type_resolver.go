package apiconfig

import (
	"errors"
	"reflect"
)

const (
	// TypeInt Tipo de dato int
	TypeInt VariableType = iota
	// TypeBool Tipo de dato bool
	TypeBool
	// TypeString Tipo de dato bool
	TypeString
	// TypeInt Tipo de dato string
	TypeNone
)

var (
	// ErrorNil error cuando la interface de entrada es nil
	ErrorNil = errors.New("the input data is nil")
	// ErrorTypeNone error con el tipo de  dato
	ErrorTypeNone = errors.New("data type not supported")
)

// VariableType tipo de dato para tipo de variable
type VariableType int

// VariableTypeResolver interface
type VariableTypeResolver interface {
	ResolveType(any) (VariableType, error)
}

type variableTypeResolver struct{}

func NewVariableTypeResolver() VariableTypeResolver {
	return &variableTypeResolver{}
}

func (v *variableTypeResolver) ResolveType(data any) (VariableType, error) {
	var typeData VariableType

	if data == nil {
		return TypeNone, ErrorNil
	}

	stringType := reflect.TypeOf(data).String()
	switch stringType {
	case "int":
		typeData = TypeInt
	case "bool":
		typeData = TypeBool
	case "string":
		typeData = TypeString
	default:
		return TypeNone, ErrorTypeNone
	}

	return typeData, nil
}
