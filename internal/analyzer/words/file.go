package words

import (
	"bufio"
	"os"
	"strings"
)

// FileProvider читает слова из файла (по одному на строку).
// Пустые строки и строки, начинающиеся с #, пропускаются.
type FileProvider struct{}

// Provide принимает путь к файлу и возвращает слова из него.
func (p FileProvider) Provide(source string) ([]string, error) {
	f, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		result = append(result, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
