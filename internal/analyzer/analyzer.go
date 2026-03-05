package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/tipo-tawerka/log-linter/internal/analyzer/rules"
	"github.com/tipo-tawerka/log-linter/internal/analyzer/words"
)

// Logger — интерфейс, который должен реализовать каждый поддерживаемый логгер.
type Logger interface {
	// MethodMsgIndex возвращает индекс аргумента-сообщения для данного метода.
	// Если метод не является логгерным, возвращает (0, false).
	MethodMsgIndex(methodName string) (int, bool)

	// IsLogCall проверяет, принадлежит ли selector-выражение данному логгеру.
	IsLogCall(pass *analysis.Pass, sel *ast.SelectorExpr) bool
}

// LoggerData хранит логгер и его метаданные.
type LoggerData struct {
	Logger Logger // Реализация интерфейса Logger
	Name   string // Уникальное имя логгера
}

var registeredLoggers []LoggerData

// addLogger регистрирует новый логгер в глобальном списке.
// Вызывает панику, если логгер с таким именем уже существует.
// Внимание: функция не потокобезопасная.
func addLogger(logger Logger, name string) {
	for _, l := range registeredLoggers {
		if l.Name == name {
			panic(fmt.Sprintf("logger with name %s already exists", name))
		}
	}
	registeredLoggers = append(registeredLoggers, LoggerData{
		Logger: logger, Name: name,
	})
}

// generateDoc формирует документацию анализатора на основе всех зарегистрированных правил.
func generateDoc() string {
	var b strings.Builder
	b.WriteString("Checks log messages with the following rules:\n\n")
	for _, r := range rules.GetAllRules() {
		b.WriteString(fmt.Sprintf("  - %s: %s\n", r.Name, r.Description))
	}
	return b.String()
}

// Analyzer — основной анализатор линтера лог-сообщений.
var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      generateDoc(),
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// Флаги анализатора.
var (
	flagRules          string // список правил через запятую
	flagSkipTests      bool   // пропускать *_test.go файлы
	flagSensitiveWords string // слова через запятую
	flagSensitiveFile  string // путь к файлу со словами
)

func init() {
	allRules := rules.GetAllRules()
	names := make([]string, len(allRules))
	for i, r := range allRules {
		names[i] = r.Name
	}
	Analyzer.Flags.StringVar(&flagRules, "rules", strings.Join(names, ","),
		"comma-separated list of enabled rules (optional, default is all)")
	Analyzer.Flags.BoolVar(&flagSkipTests, "skip-tests", false,
		"skip *_test.go files (optional, default is false)")
	Analyzer.Flags.StringVar(&flagSensitiveWords, words.FlagSensitiveWords, "",
		"comma-separated list of additional sensitive words (optional)")
	Analyzer.Flags.StringVar(&flagSensitiveFile, words.FlagSensitiveFile, "",
		"path to file with additional sensitive words (optional)")
}
