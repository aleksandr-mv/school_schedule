package helpers

import (
	"flag"
	"os"
)

// ResolveConfigPath возвращает путь к конфигу из флага --config, ENV CONFIG_PATH или дефолтный
func ResolveConfigPath(defaultPath string) string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path != "" {
		return path
	}

	if v := os.Getenv("CONFIG_PATH"); v != "" {
		return v
	}

	return defaultPath
}
