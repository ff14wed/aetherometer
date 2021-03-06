package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Config stores configuration values for the Aetherometer core
type Config struct {
	// APIPort provides the port on which the core API is served.
	APIPort uint16 `toml:"api_port" validate:"nonempty"`

	// DataPath provides the path to the folder with raw EXD files (in CSV format)
	// containing game data.
	DataPath string `toml:"data_path" validate:"directory"`

	// AdminOTP provides a one time password that the admin can use to create an
	// admin token for the API.
	AdminOTP string `toml:"admin_otp" validate:"nonempty"`

	// DisableAuth allows starting the API server without requiring an auth
	// token for queries. CORS validation will still be enforced.
	DisableAuth bool `toml:"disable_auth"`

	// AllowOrigins allows the listed sites to bypass CORS validation without
	// having to register them. Note that scheme and subdomain must be provided.
	// Example: allow_sites = ["https://plugins.foo.com"]
	AllowOrigins []string `toml:"allow_origins"`

	// Maps provides the configuration for the Map endpoint of the API.
	Maps MapConfig `toml:"maps"`

	// Adapters contains the configuration for all the adapters enabled for
	// the core API.
	Adapters Adapters `toml:"adapters"`
}

// Maps sets the configuration for the Map endpoint of the API.
type MapConfig struct {
	// Cache provides the path of the maps on the local disk.
	Cache string `toml:"cache" validate:"directory"`

	// APIPath provides the URL of an xivapi environment serving the maps if the
	// map could not be found on the local disk. Defaults to https://xivapi.com.
	APIPath string `toml:"api_path"`
}

// Adapters stores configuration structs for adapters
type Adapters struct {
	// Hook provides the configuration for the Hook adapter.
	Hook HookConfig `toml:"hook"`

	// Test is for testing purposes only. Do not use.
	Test struct{}
}

// IsEnabled returns whether or not the provided adapter name is enabled
func (a Adapters) IsEnabled(adapterName string) bool {
	rs := reflect.ValueOf(a)
	adapterConfig := rs.FieldByName(adapterName)
	if !adapterConfig.IsValid() {
		panic(fmt.Sprintf("ERROR: Adapter config for %s does not exist", adapterName))
	}
	if f := adapterConfig.FieldByName("Enabled"); f.IsValid() {
		return f.Bool()
	}
	return true
}

// SourceDirs is a table of directories that provide data used to interpret
// indexes sent over the network
type SourceDirs struct {
	MapsDir string `toml:"maps_dir" validate:"directory"`
}

func buildError(ctx []string, msg string) error {
	if len(ctx) > 0 {
		return fmt.Errorf(`config error in [%s]: %s`, strings.Join(ctx, "."), msg)
	}
	return fmt.Errorf(`config error: %s`, msg)
}

func validateFile(name string, ctx []string, val reflect.Value) error {
	err := validateNonEmpty(name, ctx, val)
	if err != nil {
		return err
	}
	filename := val.String()
	info, err := os.Stat(filename)
	if err != nil {
		return buildError(ctx, fmt.Sprintf(`%s file ("%s") does not exist`, name, filename))
	}
	if info.IsDir() {
		return buildError(ctx, fmt.Sprintf(`%s ("%s") must be a file`, name, filename))
	}
	return nil
}

func validateDir(name string, ctx []string, val reflect.Value) error {
	err := validateNonEmpty(name, ctx, val)
	if err != nil {
		return err
	}
	pathname := val.String()
	info, err := os.Stat(pathname)
	if err != nil {
		return buildError(ctx, fmt.Sprintf(`%s directory ("%s") does not exist`, name, pathname))
	}
	if !info.IsDir() {
		return buildError(ctx, fmt.Sprintf(`%s ("%s") must be a directory`, name, pathname))
	}
	return nil
}

func validateNonEmpty(name string, ctx []string, val reflect.Value) error {
	if reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface()) {
		return buildError(ctx, fmt.Sprintf("%s must be provided", name))
	}
	return nil
}

func validateField(name, validation string, ctx []string, val reflect.Value) error {
	switch validation {
	case "nonempty":
		return validateNonEmpty(name, ctx, val)
	case "file":
		return validateFile(name, ctx, val)
	case "directory":
		return validateDir(name, ctx, val)
	}
	return nil
}

func validateStruct(rs reflect.Value, ctx []string) error {
	if rs.Kind() != reflect.Struct {
		panic("BUG: improper type passed into validateStruct")
	}

	enabledField := rs.FieldByName("Enabled")
	if enabledField.IsValid() && !enabledField.Bool() {
		return nil
	}
	numFields := rs.Type().NumField()

	for i := 0; i < numFields; i++ {
		field := rs.Type().Field(i)
		validation := field.Tag.Get("validate")
		tomlTag := field.Tag.Get("toml")
		val := rs.Field(i)
		if val.Kind() == reflect.Struct {
			if err := validateStruct(val, append(ctx, tomlTag)); err != nil {
				return err
			}
		} else if validation != "" {
			if err := validateField(tomlTag, validation, ctx, val); err != nil {
				return err
			}
		}
	}
	return nil
}

// Validate returns an error when the provided configuration values do not
// pass validation
func (c *Config) Validate() error {
	rs := reflect.ValueOf(c).Elem()
	return validateStruct(rs, nil)
}
