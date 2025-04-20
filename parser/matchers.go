package parser

import (
	"regexp"
)

type Matcher interface {
	Match(line string) bool
}

type TabulatedLineMatcher struct {
	Start      *regexp.Regexp
	MatchFirst bool
	tabbedLine *regexp.Regexp
	isActive   bool
}

func CreateTabulatedLineMatcher(startRegex string, matchFirst bool) *TabulatedLineMatcher {
	return &TabulatedLineMatcher{
		Start:      regexp.MustCompile(startRegex),
		MatchFirst: matchFirst,
		tabbedLine: regexp.MustCompile(`^\s+.+$`),
		isActive:   false,
	}
}

func (m *TabulatedLineMatcher) Match(line string) bool {
	if m.Start.MatchString(line) {
		m.isActive = true
		return m.MatchFirst
	}
	if m.isActive && m.tabbedLine.MatchString(line) {
		return true
	}
	if m.isActive {
		m.isActive = false
		return false
	}
	return false
}

type StartEndBlockMatcher struct {
	Start      *regexp.Regexp
	End        *regexp.Regexp
	MatchFirst bool
	MatchLast  bool
	isActive   bool
}

func CreateStartEndBlockMatcher(startRegex string, endRegex string, matchFirst bool, matchLast bool) *StartEndBlockMatcher {
	return &StartEndBlockMatcher{
		Start:      regexp.MustCompile(startRegex),
		End:        regexp.MustCompile(endRegex),
		MatchFirst: matchFirst,
		MatchLast:  matchLast,
		isActive:   false,
	}
}

func (m *StartEndBlockMatcher) Match(line string) bool {
	if m.Start.MatchString(line) {
		m.isActive = true
		return m.MatchFirst
	}
	if m.isActive && m.End.MatchString(line) {
		m.isActive = false
		return m.MatchLast
	}
	if m.isActive {
		return true
	}
	return false
}
