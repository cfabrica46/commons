package apiconfig

import (
	"errors"

	"github.com/spf13/pflag"
)

var (
	// ErrorStruct Error en la estructira del ConfigEntry.
	ErrStruct = errors.New("error en la estructura del ConfigEntry")

	// ErrorConvert corresponde a un error en la conversion de la interface.
	ErrConvert = errors.New("error convirtiendo el tipo de dato")
)

// FlagConfigurator Interfaz que define configuración de flags.
type FlagConfigurator interface {
	ConfigureFlag(ConfigEntry) error
}

type flagConfigurator struct {
	variableTypeResolver VariableTypeResolver
}

func NewFlagConfigurator(typeResolver VariableTypeResolver) FlagConfigurator {
	return &flagConfigurator{
		variableTypeResolver: typeResolver,
	}
}

func (f *flagConfigurator) ConfigureFlag(config ConfigEntry) error {
	configValue := config.DefaultValue

	if len(config.VariableName) == 0 || len(config.Description) == 0 {
		return ErrStruct
	}

	// Se utiliza el metodo ResolveType para determinar el tipo de la interface
	variableType, err := f.variableTypeResolver.ResolveType(configValue)
	if err != nil || variableType == TypeNone {
		return ErrConvert
	}

	// Se configura el flag segun tipo de variable otorgado por ResolveType
	switch variableType {
	case TypeInt:
		if intVal, ok := configValue.(int); ok {
			pflag.IntP(config.VariableName, config.Shortcut, intVal, config.Description)
		}

	case TypeBool:
		if boolVal, ok := configValue.(bool); ok {
			pflag.BoolP(config.VariableName, config.Shortcut, boolVal, config.Description)
		}

	case TypeString:
		if stringVal, ok := configValue.(string); ok {
			pflag.StringP(config.VariableName, config.Shortcut, stringVal, config.Description)
		}
	}

	return nil
}
