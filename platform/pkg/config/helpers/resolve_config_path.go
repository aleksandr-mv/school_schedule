package helpers

import (
	"flag"
	"os"
	"strings"
	"sync"
)

// configPathState инкапсулирует состояние для thread-safety
type configPathState struct {
	flagValue string
	once      sync.Once
	mu        sync.RWMutex
}

// globalConfigPathState - единственная глобальная переменная
var globalConfigPathState = &configPathState{}

// initFlags инициализирует флаги один раз
func (s *configPathState) initFlags() {
	flag.StringVar(&s.flagValue, "config", "", "path to config file")
	if !flag.Parsed() {
		flag.Parse()
	}
}

// ResolveConfigPath возвращает путь к конфигу из флага --config, ENV CONFIG_PATH или дефолтный.
//
// Приоритет (от высшего к низшему):
//  1. Флаг командной строки --config
//  2. Переменная окружения CONFIG_PATH
//  3. Переданное дефолтное значение
//
// Thread-safe: использует мьютекс для защиты от race conditions.
// Флаги инициализируются только один раз при первом вызове.
func ResolveConfigPath(defaultPath string) string {
	// Инициализируем флаги только один раз
	globalConfigPathState.once.Do(func() {
		globalConfigPathState.initFlags()
	})

	// Thread-safe чтение флага
	globalConfigPathState.mu.RLock()
	flagValue := strings.TrimSpace(globalConfigPathState.flagValue)
	globalConfigPathState.mu.RUnlock()

	// Проверяем флаг командной строки (высший приоритет)
	if flagValue != "" {
		return flagValue
	}

	// Проверяем переменную окружения
	if v := strings.TrimSpace(os.Getenv("CONFIG_PATH")); v != "" {
		return v
	}

	// Возвращаем дефолтное значение
	return defaultPath
}
