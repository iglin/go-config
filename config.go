package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"sigs.k8s.io/yaml"
)

// Config represents storage of properties that were read from file.
type Config struct {
	properties map[string]interface{}
}

const (
	// Yaml specifies config file format. To be used in NewConfig constructor.
	Yaml = iota
	// Json specifies config file format. To be used in NewConfig constructor.
	Json
)

// NewConfig builds Config structure reading the file from path provided.
// Argument format is one of the constants: config.Yaml or config.Json.
func NewConfig(filePath string, format int) *Config {
	configHolder := Config{}
	plane, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panicf("Failed to read json config file: %v", err)
	}
	switch format {
	case Yaml:
		plane, err = yaml.YAMLToJSON(plane)
		if err != nil {
			log.Panicf("Failed to convert yaml config file to json: %v", err)
		}
		break
	case Json:
		break
	default:
		log.Panicf("Unknown config format: %v (allowed values config.Yaml, config.Json)", format)
	}

	var originalConfigMap map[string]interface{}
	if err := json.Unmarshal(plane, &originalConfigMap); err != nil {
		log.Panicf("Failed to unmarshal json config file: %v", err)
	}
	configHolder.properties = originalConfigMap
	return &configHolder
}

// GetSecret returns value read from property and decoded from base64.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
func (c *Config) GetSecret(key string) string {
	prop := c.GetProp(key)
	var strProp string
	if prop == nil {
		strProp = readStringFromEnv(key)
	} else {
		strProp = prop.(string)
	}
	if strProp == "" {
		return ""
	}
	bytes, err := base64.StdEncoding.DecodeString(strProp)
	if err != nil {
		log.Panicf("Failed to decode property %v: %v", key, err)
	}
	return string(bytes)
}

// RequireSecret returns value read from property and decoded from base64.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing this function will panic.
func (c *Config) RequireSecret(key string) string {
	prop := c.GetSecret(key)
	if prop == "" {
		log.Panic("Couldn't resolve required property " + key)
	}
	return prop
}

// GetString returns string value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing it will return the provided
// defaultVal or empty string in case there was no default specified.
func (c *Config) GetString(key string, defaultVal ...string) string {
	prop := c.GetProp(key)
	if prop == nil {
		return readStringFromEnv(key, defaultVal...)
	}
	strProp := fmt.Sprintf("%v", prop)
	return strProp
}

// RequireString returns string value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing this function will panic.
func (c *Config) RequireString(key string) string {
	prop := c.GetProp(key)
	var strVal string
	if prop == nil {
		strVal = readStringFromEnv(key)
	} else {
		strVal = fmt.Sprintf("%v", prop)
	}
	if strVal == "" {
		log.Panic("Couldn't resolve required property " + key)
	}
	return strVal
}

// GetBool returns bool value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing it will return the provided
// defaultVal or false in case there was no default specified.
func (c *Config) GetBool(key string, defaultVal ...bool) bool {
	prop := c.GetProp(key)
	if prop == nil {
		strVal := readStringFromEnv(key)
		if strVal != "" {
			return strings.EqualFold("true", strVal)
		}
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return false
	}
	return prop.(bool)
}

// RequireBool returns bool value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing this function will panic.
func (c *Config) RequireBool(key string) bool {
	prop := c.GetProp(key)
	var strVal string
	if prop == nil {
		strVal = readStringFromEnv(key)
		if strVal == "" {
			log.Panicf("Required boolean property %v is not present", key)
		} else {
			return strings.EqualFold("true", strVal)
		}
	}
	return prop.(bool)
}

// GetInt returns int value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing it will return the provided
// defaultVal or 0 in case there was no default specified.
func (c *Config) GetInt(key string, defaultVal ...int) int {
	prop := c.GetProp(key)
	if prop == nil {
		strVal := readStringFromEnv(key)
		if strVal != "" {
			if res, err := strconv.Atoi(strVal); err != nil {
				log.Panicf("Failed to convert env value for key %s to int: %s", key, err)
			} else {
				return res
			}
		}
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	return int(prop.(float64))
}

// RequireInt returns int value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing this function will panic.
func (c *Config) RequireInt(key string) int {
	prop := c.GetProp(key)
	if prop == nil {
		strVal := readStringFromEnv(key)
		if strVal != "" {
			if res, err := strconv.Atoi(strVal); err != nil {
				log.Panicf("Failed to convert env value for key %s to int: %s", key, err)
			} else {
				return res
			}
		}
		log.Panicf("Failed to find required property for key %s", key)
	}
	return int(prop.(float64))
}

// GetFloat64 returns float64 value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing it will return the provided
// defaultVal or 0 in case there was no default specified.
func (c *Config) GetFloat64(key string, defaultVal ...float64) float64 {
	prop := c.GetProp(key)
	if prop == nil {
		strVal := readStringFromEnv(key)
		if strVal != "" {
			if res, err := strconv.ParseFloat(strVal, 64); err != nil {
				log.Panicf("Failed to convert env value for key %s to float64: %s", key, err)
			} else {
				return res
			}
		}
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	return prop.(float64)
}

// RequireFloat64 returns float64 value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing this function will panic.
func (c *Config) RequireFloat64(key string) float64 {
	prop := c.GetProp(key)
	var strVal string
	if prop == nil {
		strVal = readStringFromEnv(key)
		if strVal == "" {
			log.Panicf("Required float64 property %v is not present", key)
		} else {
			if res, err := strconv.ParseFloat(strVal, 64); err != nil {
				log.Panicf("Failed to convert env value for key %s to float64: %s", key, err)
			} else {
				return res
			}
		}
	}
	return prop.(float64)
}

// GetFloat32 returns float32 value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing it will return the provided
// defaultVal or 0 in case there was no default specified.
func (c *Config) GetFloat32(key string, defaultVal ...float32) float32 {
	prop := c.GetProp(key)
	if prop == nil {
		strVal := readStringFromEnv(key)
		if strVal != "" {
			if res, err := strconv.ParseFloat(strVal, 32); err != nil {
				log.Panicf("Failed to convert env value for key %s to float64: %s", key, err)
			} else {
				return float32(res)
			}
		}
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	return float32(prop.(float64))
}

// RequireFloat32 returns float32 value read from property.
//
// If property for the specified key is missing, it will try to read value from
// the environment variable. Environment variable name will be constructed by
// formatting property key to upper case and replacing dots with underscore,
// e.g. property 'my.test.property1' will be translated to 'MY_TEST_PROPERTY1'.
//
// If both property and env variable are missing this function will panic.
func (c *Config) RequireFloat32(key string) float32 {
	prop := c.GetProp(key)
	var strVal string
	if prop == nil {
		strVal = readStringFromEnv(key)
		if strVal == "" {
			log.Panicf("Required float32 property %v is not present", key)
		} else {
			if res, err := strconv.ParseFloat(strVal, 32); err != nil {
				log.Panicf("Failed to convert env value for key %s to float64: %s", key, err)
			} else {
				return float32(res)
			}
		}
	}
	return float32(prop.(float64))
}

// GetProp returns value read from property as interface{}.
// The function will not try to lookup environment variable if property is missing.
// If no property found for the key the function returns nil.
func (c *Config) GetProp(key string) interface{} {
	return findPropInMap(key, c.properties)
}

func findPropInMap(key string, props map[string]interface{}) interface{} {
	if val, ok := props[key]; ok {
		return resolveProp(val)
	}
	dotIdx := strings.Index(key, ".")
	if dotIdx != -1 {
		prefix := key[:dotIdx]
		suffix := key[dotIdx+1:]
		if val, ok := props[prefix]; ok {
			var foundProp interface{}
			if nextSubmap, ok := val.(map[string]interface{}); ok {
				foundProp = findPropInMap(suffix, nextSubmap)
			}
			if foundProp != nil {
				return foundProp
			}
		}
	}
	return nil
}

func resolveProp(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	if _, ok := val.(map[string]interface{}); ok {
		return nil
	}
	return val
}

func readStringFromEnv(propertyKey string, defaultVal ...string) string {
	envName := strings.ToUpper(propertyKey)
	envName = strings.ReplaceAll(envName, ".", "_")
	env := os.Getenv(envName)
	if env == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return env
}
