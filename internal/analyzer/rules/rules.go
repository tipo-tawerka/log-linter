package rules

// Rule — это интерфейс для проверки соответствия текста лога правилам.
type Rule interface {
	// IsCorrect проверяет, соответствует ли текст лога правилу.
	// Возвращает true, если текст корректен, и false в противном случае.
	IsCorrect(text string) bool
}

// CorrectionHint — это интерфейс для предоставления подсказки
// по исправлению текста лога, если он не прошел проверку Rule.
type CorrectionHint interface {
	// ToCorrect принимает ошибочный текст лога и возвращает исправленный вариант.
	ToCorrect(text string) string
}

// RuleData хранит само правило и его метаданные (имя и описание).
type RuleData struct {
	Rule        Rule   // Правило для проверки текста лога
	Name        string // Уникальное имя правила
	Description string // Текстовое описание того, что делает правило
}

var registeredRules []RuleData

// AddRule регистрирует новое правило в глобальном списке.
// Вызывает панику, если правило с таким именем уже существует.
// Внимание, функция не потокобезопасная, предполагается, что она будет вызываться только в init() функции пакетов
// которые используют правила.
func AddRule(rule Rule, name string, description string) {
	for _, r := range registeredRules {
		if r.Name == name {
			panic("rule " + name + " already exists")
		}
	}
	registeredRules = append(registeredRules, RuleData{
		Rule: rule, Name: name, Description: description,
	})
}

// GetRules возвращает срез всех зарегистрированных правил.
func GetRules() []RuleData {
	return registeredRules
}
