package config

import (
	assertions "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig_Yaml(t *testing.T) {
	config := NewConfig("./test_config.yaml", Yaml)
	testConfigPositiveCases(t, config)
}

func TestConfig_Json(t *testing.T) {
	config := NewConfig("./test_config.json", Json)
	testConfigPositiveCases(t, config)
}

func TestConfig_YamlConversionErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on yaml conversion")
		}
	}()
	_ = NewConfig("./config_test.go", Yaml)
}

func TestConfig_ParsingErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on json parsing error")
		}
	}()
	_ = NewConfig("./config_test.go", Json)
}

func testConfigPositiveCases(t *testing.T, config *Config) {
	assert := assertions.New(t)

	assert.Equal("test11", config.GetString("root.family1.key1"))
	assert.Equal("test11", config.RequireString("root.family1.key1"))

	assert.Equal("test121", config.GetString("root.family1.key2.subkey1"))
	assert.Equal("test121", config.RequireString("root.family1.key2.subkey1"))

	assert.Equal(122, config.GetInt("root.family1.key2.subkey2"))
	assert.Equal(122, config.RequireInt("root.family1.key2.subkey2"))
	assert.Equal(float32(122), config.GetFloat32("root.family1.key2.subkey2"))
	assert.Equal(float32(122), config.RequireFloat32("root.family1.key2.subkey2"))
	assert.Equal(float64(122), config.GetFloat64("root.family1.key2.subkey2"))
	assert.Equal(float64(122), config.RequireFloat64("root.family1.key2.subkey2"))

	assert.Equal("test2", config.GetString("root.family2"))
	assert.Equal("test2", config.RequireString("root.family2"))

	assert.True(config.GetBool("root.family3.key1"))
	assert.True(config.RequireBool("root.family3.key1"))
	assert.False(config.GetBool("root.family3.key2"))
	assert.False(config.RequireBool("root.family3.key2"))

	assert.Equal(211, config.GetInt("subroot.family1.key1"))
	assert.Equal(211, config.RequireInt("subroot.family1.key1"))
	assert.Equal(float32(211), config.GetFloat32("subroot.family1.key1"))
	assert.Equal(float32(211), config.RequireFloat32("subroot.family1.key1"))
	assert.Equal(float64(211), config.GetFloat64("subroot.family1.key1"))
	assert.Equal(float64(211), config.RequireFloat64("subroot.family1.key1"))

	assert.Equal(float32(212.212), config.GetFloat32("subroot.family1.key2"))
	assert.Equal(float32(212.212), config.RequireFloat32("subroot.family1.key2"))
	assert.Equal(212.212, config.GetFloat64("subroot.family1.key2"))
	assert.Equal(212.212, config.RequireFloat64("subroot.family1.key2"))

	assert.Equal("subtest_secret", config.GetSecret("subroot.family1.key3.secret"))
	assert.Equal("subtest_secret", config.RequireSecret("subroot.family1.key3.secret"))

	assert.Equal(3, config.GetInt("simpleprop"))
	assert.Equal(3, config.RequireInt("simpleprop"))
	assert.Equal(float32(3), config.GetFloat32("simpleprop"))
	assert.Equal(float32(3), config.RequireFloat32("simpleprop"))
	assert.Equal(float64(3), config.GetFloat64("simpleprop"))
	assert.Equal(float64(3), config.RequireFloat64("simpleprop"))

	assert.Equal(4, config.GetInt("another.simple.prop"))
	assert.Equal(4, config.RequireInt("another.simple.prop"))
	assert.Equal(float32(4), config.GetFloat32("another.simple.prop"))
	assert.Equal(float32(4), config.RequireFloat32("another.simple.prop"))
	assert.Equal(float64(4), config.GetFloat64("another.simple.prop"))
	assert.Equal(float64(4), config.RequireFloat64("another.simple.prop"))
}

func TestGetMissingProperty(t *testing.T) {
	assert := assertions.New(t)

	config := NewConfig("./test_config.yaml", Yaml)

	assert.Equal("", config.GetSecret("missing_property"))
	assert.Equal("", config.GetString("missing_property"))
	assert.Equal(false, config.GetBool("missing_property"))
	assert.Equal(0, config.GetInt("missing_property"))
	assert.Equal(float64(0), config.GetFloat64("missing_property"))
	assert.Equal(float32(0), config.GetFloat32("missing_property"))
	assert.Nil(config.GetProp("missing_property"))
}

func TestConfigEnvs(t *testing.T) {
	assert := assertions.New(t)

	config := NewConfig("./test_config.yaml", Yaml)

	err := os.Setenv("SECRET_ENV", "c2VjcmV0X3ZhbA==")
	assert.Nil(err)
	assert.Equal("secret_val", config.GetSecret("SECRET_ENV"))
	assert.Equal("secret_val", config.RequireSecret("SECRET_ENV"))
	err = os.Unsetenv("SECRET_ENV")
	assert.Nil(err)

	assert.Equal("default_val", config.GetString("root.family1.key19", "default_val"))
	err = os.Setenv("ROOT_FAMILY1_KEY19", "test_val")
	assert.Nil(err)
	assert.Equal("test_val", config.GetString("root.family1.key19", "default_val"))
	assert.Equal("test_val", config.RequireString("root.family1.key19"))
	err = os.Setenv("ROOT_FAMILY1_KEY1", "test_val")
	assert.Nil(err)
	assert.Equal("test11", config.GetString("root.family1.key1", "default_val"))

	assert.Equal(99, config.GetInt("root.family1.key2.subkey29", 99))
	err = os.Setenv("ROOT_FAMILY1_KEY2_SUBKEY29", "8")
	assert.Nil(err)
	assert.Equal(8, config.GetInt("root.family1.key2.subkey29", 99))
	assert.Equal(8, config.RequireInt("root.family1.key2.subkey29"))
	err = os.Setenv("ROOT_FAMILY1_KEY2_SUBKEY2", "8")
	assert.Nil(err)
	assert.Equal(122, config.GetInt("root.family1.key2.subkey2"))

	assert.Equal(float32(99), config.GetFloat32("subroot.family1.key29", 99))
	err = os.Setenv("SUBROOT_FAMILY1_KEY29", "2.2")
	assert.Nil(err)
	assert.Equal(float32(2.2), config.GetFloat32("subroot.family1.key29", 99))
	assert.Equal(float32(2.2), config.RequireFloat32("subroot.family1.key29"))
	err = os.Setenv("SUBROOT_FAMILY1_KEY2", "2.2")
	assert.Nil(err)
	assert.Equal(float32(212.212), config.GetFloat32("subroot.family1.key2"))

	assert.Equal(float64(99), config.GetFloat64("subroot.family1.key2.subkey19", 99))
	err = os.Setenv("SUBROOT_FAMILY1_KEY2_SUBKEY19", "2.2")
	assert.Nil(err)
	assert.Equal(2.2, config.GetFloat64("subroot.family1.key2.subkey19", 99))
	assert.Equal(2.2, config.RequireFloat64("subroot.family1.key2.subkey19"))
	err = os.Setenv("SUBROOT_FAMILY1_KEY2_SUBKEY1", "2.2")
	assert.Nil(err)
	assert.Equal(2121.2121, config.GetFloat64("subroot.family1.key2.subkey1"))

	assert.False(config.GetBool("root.family3.key19", false))
	err = os.Setenv("ROOT_FAMILY3_KEY19", "true")
	assert.Nil(err)
	assert.True(config.GetBool("root.family3.key19", false))
	assert.True(config.RequireBool("root.family3.key19"))
	err = os.Setenv("ROOT_FAMILY3_KEY1", "false")
	assert.Nil(err)
	assert.True(config.GetBool("root.family3.key1"))

	assert.True(config.GetBool("root.family3.key29", true))
	err = os.Setenv("ROOT_FAMILY3_KEY29", "false")
	assert.Nil(err)
	assert.False(config.GetBool("root.family3.key29", true))
	err = os.Setenv("ROOT_FAMILY3_KEY2", "true")
	assert.Nil(err)
	assert.False(config.GetBool("root.family3.key2"))

	clearEnvs(assert)
}

func TestConfigDefaults(t *testing.T) {
	assert := assertions.New(t)

	config := NewConfig("./test_config.yaml", Yaml)

	assert.Equal("test11", config.GetString("root.family1.key1"))
	assert.Equal("default_val", config.GetString("root.family1.key19", "default_val"))

	assert.Equal("test121", config.GetString("root.family1.key2.subkey1"))
	assert.Equal("default_val", config.GetString("root.family1.key2.subkey19", "default_val"))

	assert.Equal(122, config.GetInt("root.family1.key2.subkey2"))
	assert.Equal(99, config.GetInt("root.family1.key2.subkey29", 99))
	assert.Equal(float32(122), config.GetFloat32("root.family1.key2.subkey2"))
	assert.Equal(float32(99), config.GetFloat32("root.family1.key2.subkey29", 99))
	assert.Equal(float64(122), config.GetFloat64("root.family1.key2.subkey2"))
	assert.Equal(float64(99), config.GetFloat64("root.family1.key2.subkey29", 99))

	assert.Equal("test2", config.GetString("root.family2"))
	assert.Equal("default_val", config.GetString("root.family29", "default_val"))

	assert.True(config.GetBool("root.family3.key1"))
	assert.False(config.GetBool("root.family3.key19", false))
	assert.False(config.GetBool("root.family3.key2"))
	assert.True(config.GetBool("root.family3.key29", true))

	assert.Equal(float32(212.212), config.GetFloat32("subroot.family1.key2"))
	assert.Equal(float32(99.99), config.GetFloat32("subroot.family1.key29", float32(99.99)))
	assert.Equal(212.212, config.GetFloat64("subroot.family1.key2"))
	assert.Equal(99.99, config.GetFloat64("subroot.family1.key29", 99.99))

	assert.Equal(99, config.GetInt("simpleprop9", 99))

	assert.Equal(99, config.GetInt("another.simple.prop99", 99))
}

func TestConfig_RequireBool(t *testing.T) {
	config := NewConfig("./test_config.yaml", Yaml)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on missing property")
		}
	}()
	config.RequireBool("missing.property")
}

func TestConfig_RequireString(t *testing.T) {
	config := NewConfig("./test_config.yaml", Yaml)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on missing property")
		}
	}()
	config.RequireString("missing.property")
}

func TestConfig_RequireInt(t *testing.T) {
	config := NewConfig("./test_config.yaml", Yaml)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on missing property")
		}
	}()
	config.RequireInt("missing.property")
}

func TestConfig_RequireFloat32(t *testing.T) {
	config := NewConfig("./test_config.yaml", Yaml)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on missing property")
		}
	}()
	config.RequireFloat32("missing.property")
}

func TestConfig_RequireFloat64(t *testing.T) {
	config := NewConfig("./test_config.yaml", Yaml)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on missing property")
		}
	}()
	config.RequireFloat64("missing.property")
}

func TestConfig_RequireSecret(t *testing.T) {
	config := NewConfig("./test_config.yaml", Yaml)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on missing property")
		}
	}()
	config.RequireSecret("missing.property")
}

func TestConfig_GetSecret_Base64DecodeErr(t *testing.T) {
	assert := assertions.New(t)

	config := NewConfig("./test_config.yaml", Yaml)

	err := os.Setenv("TEST_ENV_VAR", "val")
	assert.Nil(err)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic in base64 decode")
		}
	}()
	config.GetSecret("test.env.var")

	err = os.Unsetenv("TEST_ENV_VAR")
	assert.Nil(err)
}

func TestConfig_RequireSecret_Base64DecodeErr(t *testing.T) {
	assert := assertions.New(t)

	config := NewConfig("./test_config.yaml", Yaml)
	err := os.Setenv("TEST_ENV_VAR", "val")
	assert.Nil(err)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic in base64 decode")
		}
	}()
	config.RequireSecret("test.env.var")

	err = os.Unsetenv("TEST_ENV_VAR")
	assert.Nil(err)
}

func clearEnvs(assert *assertions.Assertions) {
	err := os.Unsetenv("ROOT_FAMILY1_KEY19")
	assert.Nil(err)
	err = os.Unsetenv("ROOT_FAMILY1_KEY1")
	assert.Nil(err)

	err = os.Unsetenv("ROOT_FAMILY1_KEY2_SUBKEY29")
	assert.Nil(err)
	err = os.Unsetenv("ROOT_FAMILY1_KEY2_SUBKEY2")
	assert.Nil(err)

	err = os.Unsetenv("SUBROOT_FAMILY1_KEY29")
	assert.Nil(err)
	err = os.Unsetenv("SUBROOT_FAMILY1_KEY2")
	assert.Nil(err)

	err = os.Unsetenv("SUBROOT_FAMILY1_KEY2_SUBKEY19")
	assert.Nil(err)
	err = os.Unsetenv("SUBROOT_FAMILY1_KEY2_SUBKEY1")
	assert.Nil(err)

	err = os.Unsetenv("ROOT_FAMILY3_KEY19")
	assert.Nil(err)
	err = os.Unsetenv("ROOT_FAMILY3_KEY1")
	assert.Nil(err)

	err = os.Unsetenv("ROOT_FAMILY3_KEY29")
	assert.Nil(err)
	err = os.Unsetenv("ROOT_FAMILY3_KEY2")
	assert.Nil(err)
}
