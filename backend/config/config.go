package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Config stores configuration values for the Sibyl backend
type Config struct {
	HookDLL      string `toml:"hook_dll" validate:"file"`
	FFXIVProcess string `toml:"ffxiv_process" validate:"nonempty"`
	APIPort      uint16 `toml:"api_port" validate:"nonempty"`

	Sources SourceDirs `toml:"sources"`
}

// SourceDirs is a table of directories that provide data used to interpret
// indexes sent over the network
type SourceDirs struct {
	MapsDir string `toml:"maps_dir" validate:"directory"`
	DataDir string `toml:"data_dir" validate:"directory"`
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
