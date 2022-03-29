package reporter

import (
	"os"

	"github.com/sqweek/dialog"
)

func FatalError(err error) {
	msgBuilder := dialog.Message("Aetherometer has encountered an error and must exit. Please check the configuration and/or logs: %s", err.Error())
	msgBuilder.Error()
	os.Exit(1)
}
