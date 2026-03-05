package words

import "fmt"

// Имена флагов анализатора, привязанные к провайдерам.
const (
	FlagSensitiveWords = "sensitive-words"      // CSV через запятую
	FlagSensitiveFile  = "sensitive-words-file" // путь к файлу
)

// Provider — интерфейс для получения списка чувтсвительных слов.
// Принимает строку-источник и возвращает распознанные слова.
type Provider interface {
	// Provide возвращает список слов из переданного источника.
	Provide(source string) ([]string, error)
}

// NewProvider возвращает провайдер, соответствующий имени флага.
func NewProvider(flagName string) (Provider, error) {
	switch flagName {
	case FlagSensitiveWords:
		return CLIProvider{}, nil
	case FlagSensitiveFile:
		return FileProvider{}, nil
	default:
		return nil, fmt.Errorf("unknown provider flag: %s", flagName)
	}
}
