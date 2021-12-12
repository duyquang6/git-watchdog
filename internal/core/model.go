package core

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
)

type Rule struct {
	ID          string         `json:"id"`
	Description string         `json:"description"`
	Severity    string         `json:"severity"`
	Type        string         `json:"type"`
	RegexStr    string         `json:"regex"`
	Regexp      *regexp.Regexp `json:"-"`
}

type Location struct {
	Path     string   `json:"path"`
	Position Position `json:"position"`
}

type Position struct {
	Begin Begin `json:"begin"`
}

type Begin struct {
	Line uint `json:"line"`
}

type FindingMetadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type Finding struct {
	Type     string          `json:"type"`
	RuleID   string          `json:"ruleId"`
	Location Location        `json:"location"`
	Metadata FindingMetadata `json:"metadata"`
}

func NewRulesFromFile(rulesFilePath string) []Rule {
	var rules []Rule
	rulesBytes, err := ioutil.ReadFile(rulesFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(rulesBytes, &rules)
	if err != nil {
		panic(err)
	}
	for i := range rules {
		rules[i].Regexp = regexp.MustCompile(rules[i].RegexStr)
	}
	return rules
}
