package main

import (
	"os"
	"web-hosting/internal/package/script"
	"web-hosting/internal/providers"

	"github.com/samber/do/v2"
)

func args(injector do.Injector) bool {
	if len(os.Args) > 1 {
		flag := script.Commands(injector)
		return flag
	}

	return true
}

func main() {
	var (
		injector = do.New()
	)

	providers.RegisterProviders(injector)

	if !args(injector) {
		return
	}

}
