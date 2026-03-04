package rules

import (
	"fmt"
)

// Rule — это интерфейс для проверки соответствия текста лога правилам.
type Rule interface {
	// IsCorrect проверяет, соответствует ли текст лога правилу.
	// Возвращает true, если текст корректен, и false в противном случае.
	IsCorrect(text string) bool
}

// CorrectionHint — это интерфейс для предоставления подсказки
// по исправлению текста лога, если он не прошел проверку Rule.
type CorrectionHint interface {
	// Fix принимает ошибочный текст лога и возвращает исправленный вариант.
	Fix(text string) string
}

// RuleData хранит само правило и его метаданные.
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
func addRule(rule Rule, name string, description string) {
	for _, r := range registeredRules {
		if r.Name == name {
			panic(fmt.Sprintf("rule with name %s already exist", name))
		}
	}
	registeredRules = append(registeredRules, RuleData{
		Rule: rule, Name: name, Description: description,
	})
}

// GetRules возвращает срез правил, используемых для проверки текстов логов.
func GetRules(names []string) ([]RuleData, error) {
	rules := make([]RuleData, 0, len(names))
	for _, name := range names {
		found := false
		for _, r := range registeredRules {
			if r.Name == name {
				rules = append(rules, r)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("rule with name %s not found", name)
		}
	}
	return rules, nil
}

func GetAllRules() []RuleData {
	return registeredRules
}
