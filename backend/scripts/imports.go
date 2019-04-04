// +build tools

package scripts

// Import packages necessary for tools but not necessary for
// compilation
import (
	_ "github.com/99designs/gqlgen/cmd"
	_ "github.com/99designs/gqlgen/codegen"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
)
