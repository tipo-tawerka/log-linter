package rules

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RulesRegistrySuite struct {
	suite.Suite
	originalRules []RuleData
}

func (s *RulesRegistrySuite) SetupTest() {
	s.originalRules = make([]RuleData, len(registeredRules))
	copy(s.originalRules, registeredRules)
}

func (s *RulesRegistrySuite) TeardownTest() {
	registeredRules = s.originalRules
}

func (s *RulesRegistrySuite) TestGetRules_ReturnsRequested() {
	result, err := GetRules([]string{"lowercase", "english"})
	s.Require().NoError(err)
	s.Require().Len(result, 2)
	s.Equal("lowercase", result[0].Name)
	s.Equal("english", result[1].Name)
}

func (s *RulesRegistrySuite) TestGetRules_ErrorOnUnknown() {
	_, err := GetRules([]string{"nonexistent"})
	s.Error(err)
	s.Contains(err.Error(), "nonexistent")
}

func (s *RulesRegistrySuite) TestGetRules_EmptyList() {
	result, err := GetRules([]string{})
	s.Require().NoError(err)
	s.Empty(result)
}

func (s *RulesRegistrySuite) TestGetAllRules_ContainsDefaults() {
	all := GetAllRules()
	names := make([]string, len(all))
	for i, r := range all {
		names[i] = r.Name
	}
	s.Contains(names, "lowercase")
	s.Contains(names, "english")
	s.Contains(names, "specchar")
	s.Contains(names, "sensitive")
}

func (s *RulesRegistrySuite) TestAddRule_PanicOnDuplicate() {
	addRule(RuleLowercase{}, "unique_rule_xyz", "desc")

	s.Panics(func() {
		addRule(RuleLowercase{}, "unique_rule_xyz", "desc")
	})
}

func TestRulesRegistry(t *testing.T) {
	suite.Run(t, new(RulesRegistrySuite))
}
