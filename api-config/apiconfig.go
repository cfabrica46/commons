package apiconfig

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	ErrTypeResolver = errors.New("error on ResolveType")
	ErrTypeMap      = errors.New("error assign value on map")
)

// ConfigEntry entrada de configuración/
type ConfigEntry struct {
	DefaultValue any
	VariableName string
	Description  string
	Shortcut     string
}

// CfgBase configuración basica de la API.
type CfgBase struct {
	Port      string
	URIPrefix string
	Timeout   time.Duration
}

// Configurator Interfaz configurador de apiconfig.
type Configurator interface {
	Configure(entries []ConfigEntry) (map[string]any, error)
}

type configurator struct {
	flagConfigurator FlagConfigurator
	typeResolver     VariableTypeResolver
}

const (
	noEntriesError = "no entries given"
)

// NewConfigurator constructor.
func NewConfigurator(flagConfigurator FlagConfigurator, typeResolver VariableTypeResolver) Configurator {
	return &configurator{
		flagConfigurator: flagConfigurator,
		typeResolver:     typeResolver,
	}
}

func (c *configurator) Configure(entries []ConfigEntry) (map[string]any, error) {
	if len(entries) == 0 {
		return nil, errors.New(noEntriesError)
	}

	// Configuration for environmental variables on the system
	viper.AutomaticEnv()

	// Configuration for flags
	if err := flagConfiguration(entries, c.flagConfigurator); err != nil {
		return nil, err
	}

	values := make(map[string]any)

	for i := range entries {
		var ok bool

		valType, err := validateEnty(entries[i], c.typeResolver)
		if err != nil {
			return nil, err
		}

		name := entries[i].VariableName
		val := viperConfiguration(name, entries[i])

		switch valType {
		case TypeInt:
			values[name], ok = val.(int)
		case TypeBool:
			values[name], ok = val.(bool)
		case TypeString:
			values[name], ok = val.(string)
		}

		if !ok {
			return nil, fmt.Errorf("%w: type: %v -  val: %v", ErrTypeMap, valType, val)
		}
	}

	return values, nil
}

func validateEnty(entry ConfigEntry, typeResolver VariableTypeResolver) (VariableType, error) {
	val, err := typeResolver.ResolveType(entry.DefaultValue)
	if err != nil {
		return TypeNone, ErrTypeResolver
	}

	return val, nil
}

// Sets the value by getting the viper value If there's no viper value, sets the default value.
func viperConfiguration(name string, entry ConfigEntry) any {
	val := viper.Get(name)
	if val == nil {
		val = entry.DefaultValue
	}

	return val
}

// Sets flags using flagConfigurator
func flagConfiguration(entries []ConfigEntry, flagConfigurator FlagConfigurator) error {
	for i := range entries {
		if err := flagConfigurator.ConfigureFlag(entries[i]); err != nil {
			return err
		}
	}

	pflag.Parse()

	return viper.BindPFlags(pflag.CommandLine)
}
