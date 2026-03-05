package words

import "strings"

// CLIProvider парсит слова из строки, разделённой запятыми.
type CLIProvider struct{}

// Provide разбивает source по запятым, тримит пробелы, пропускает пустые элементы.
func (p CLIProvider) Provide(source string) ([]string, error) {
	parts := strings.Split(source, ",")
	result := make([]string, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result, nil
}
