package configuration

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const separatorPattern string = "_"

type Binder struct {
	expandableValuePattern *regexp.Regexp
}

func New() *Binder {
	return &Binder{expandableValuePattern: regexp.MustCompile(`\${(?P<variable>.*?)((:-?)(?P<default>.*?))?}`)}
}

// ExpandMap expands the map: ENV_PROPERTY_VALUE => {"ENV" : { "PROPERTY" : "VALUE"}}
func (binder *Binder) ExpandMap(mappings map[string]string) map[string]interface{} {
	expandedMap := make(map[string]interface{})
	for key, value := range mappings {
		binder.expand(expandedMap, key, value)
	}
	return expandedMap
}

func (binder *Binder) expand(mapping map[string]interface{}, key string, value interface{}) {
	if strings.Contains(key, separatorPattern) {
		currentKey := strings.Split(key, separatorPattern)[0]
		if innerMap, ok := (mapping[currentKey]).(map[string]interface{}); ok {
			binder.expand(innerMap, substring(key, separatorPattern), value)
		} else {
			innerMap := make(map[string]interface{})
			mapping[currentKey] = innerMap
			binder.expand(innerMap, substring(key, separatorPattern), value)
		}
	} else {
		mapping[key] = value
	}
}

func substring(value string, delimiter string) string {
	index := strings.Index(value, delimiter)
	return value[index+1:]
}

func getEnvironment() map[string]string {
	settings := make(map[string]string)
	entries := os.Environ()
	for _, entry := range entries {
		split := strings.Split(entry, "=")
		settings[strings.ToLower(split[0])] = split[1]
	}
	return settings
}

// Bind the mapping to the configInterface
func (binder *Binder) Bind(mapping map[string]string, configInterface interface{}) error {
	expandedMapping := binder.ExpandMap(mapping)
	data, marshalErr := json.Marshal(expandedMapping)
	if marshalErr != nil {
		return marshalErr
	}
	unmarshalErr := json.Unmarshal(data, configInterface)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	return nil
}

// BindEnvironment the configInterface to the environment variables
func (binder *Binder) BindEnvironment(configInterface interface{}) (err error) {
	err = binder.Bind(getEnvironment(), configInterface)
	return
}

func (binder *Binder) InitializeConfigFromYaml(yamlData []byte, config interface{}) error {
	dataMap := make(map[interface{}]interface{})
	err := yaml.Unmarshal(yamlData, &dataMap)
	if err != nil {
		return err
	}
	bindErr := binder.bindConfig(dataMap, config)
	if bindErr != nil {
		return bindErr
	}
	return nil
}

// BindConfig with config tag
func (binder *Binder) bindConfig(data map[interface{}]interface{}, a interface{}) error {
	elementValue := reflect.ValueOf(a).Elem()
	elementType := reflect.TypeOf(a).Elem()

	for j := 0; j < elementValue.NumField(); j++ {
		field := elementValue.Field(j)
		if !field.IsValid() && !field.CanSet() {
			continue
		}
		fieldType := elementType.Field(j)
		configTag := elementValue.Type().Field(j).Tag.Get("config")
		if configTag == "" {
			return fmt.Errorf("no config tag")
		}
		value, ok := data[configTag]
		if !ok {
			return fmt.Errorf("no value found for config tag %s", configTag)
		}
		switch fieldType.Type.Kind() {
		case reflect.Bool:
			newValue, err := binder.expandValue(value)
			if err != nil {
				return err
			}
			value, err := strconv.ParseBool(fmt.Sprint(newValue))
			if err != nil {
				return err
			}
			field.SetBool(value)
		case reflect.Struct:
			entry, ok := value.(map[interface{}]interface{})
			if !ok {
				return fmt.Errorf("value: %+v for kind struct must be map[interface{}]interface{}", value)
			}

			ptr := reflect.PtrTo(elementValue.Type().Field(j).Type)
			structure := reflect.New(ptr.Elem())
			field.Set(structure.Elem())
			err := binder.bindConfig(entry, structure.Interface())
			if err != nil {
				return err
			}
			field.Set(structure.Elem())
		case reflect.Float64:
			newValue, err := binder.expandValue(value)
			if err != nil {
				return err
			}
			value, err := strconv.ParseFloat(fmt.Sprint(newValue), 64)
			if err != nil {
				return err
			}
			field.SetFloat(value)
		case reflect.Int:
			newValue, err := binder.expandValue(value)
			if err != nil {
				return err
			}
			value, err := strconv.ParseInt(fmt.Sprint(newValue), 0, 0)
			if err != nil {
				return err
			}
			field.SetInt(value)
		case reflect.String:
			newValue, err := binder.expandValue(value)
			if err != nil {
				return err
			}
			field.SetString(fmt.Sprint(newValue))
		default:
			return fmt.Errorf("value: %+v for kind not supported yet", value)
		}
	}
	return nil
}

func (binder *Binder) expandValue(value interface{}) (interface{}, error) {
	valueStr := fmt.Sprint(value)
	match := binder.expandableValuePattern.FindStringSubmatch(valueStr)
	if len(match) > 0 {
		result := make(map[string]string)
		matches := binder.expandableValuePattern.FindStringSubmatch(valueStr)
		names := binder.expandableValuePattern.SubexpNames()
		for i, match := range matches {
			if i != 0 {
				result[names[i]] = match
			}
		}
		variable := ""
		defaultValue := ""
		if resultValue, found := result["variable"]; found {
			variable = resultValue
		}
		if resultValue, found := result["default"]; found {
			defaultValue = resultValue
		}
		newValue := binder.getEnvironmentValue(variable, defaultValue)
		if newValue == "" {
			return nil, fmt.Errorf("enviornment value is missing for environment variable %s", variable)
		}
		return newValue, nil
	} else {
		return value, nil
	}
}

func (binder *Binder) getEnvironmentValue(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
