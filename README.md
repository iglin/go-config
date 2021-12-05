# go-config

[![codecov](https://codecov.io/gh/iglin/go-config/branch/main/graph/badge.svg?token=VI8IH1PPKS)](https://codecov.io/gh/iglin/go-config)

Golang library for reading properties from configuration files in JSON and YAML format or from environment variables. 

# Usage
Create config instance and read properties from it. Supported file formats are JSON and YAML. 

If property is missing in config file, library will try to look up it in envirnment variables: in this case property name will be fomatted to upper case and all dots will be replaced with `_`, e.g. property 'my.test.property1' will be translated to `MY_TEST_PROPERTY1` envirnoment variable name. 

```go

import goconfig "github.com/iglin/go-config"

func main() {
	config := goconfig.NewConfig("./test_config.yaml", goconfig.Yaml)

	// reading strings
	
	strVal := config.GetString("root.family1.key1")
	strValOrDefault := config.GetString("root.family1.key1", "my-default-val")
	// panics if both property and env variable ROOT_FAMILY1_KEY1 are missing
	requiredStrVal := config.RequireString("root.family1.key1")
	
	// reading ints
	
	intVal := config.GetInt("root.family1.key1")
	intValOrDefault := config.GetInt("root.family1.key1", 1)
	// panics if both property and env variable ROOT_FAMILY1_KEY1 are missing
	requiredIntVal := config.RequireInt("root.family1.key1")
	
	// reading floats
	
	floatVal := config.GetFloat64("root.family1.key1")
	floatValOrDefault := config.GetFloat64("root.family1.key1", 1.1)
	// panics if both property and env variable ROOT_FAMILY1_KEY1 are missing
	requiredFloatVal := config.RequireFloat32("root.family1.key1")
	
	// reading bools
	
	boolVal := config.GetBool("root.family1.key1")
	boolValOrDefault := config.GetBool("root.family1.key1", true)
	// panics if both property and env variable ROOT_FAMILY1_KEY1 are missing
	requiredBoolVal := config.RequireBool("root.family1.key1")
	
	// reading secrets (decodes base64 value before returning the result)
	
	secretStringVal := config.GetSecret("root.family1.key1")
	secretStringValOrDefault := config.GetBool("root.family1.key1", "default-val")
	// panics if both property and env variable ROOT_FAMILY1_KEY1 are missing
	requiredSecretVal := config.RequireBool("root.family1.key1")
}
```

For more examples see test config file [test_config.yaml](./test_config.yaml) and [./config_test.go](./config_test.go)
