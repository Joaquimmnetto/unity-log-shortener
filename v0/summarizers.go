package v0

import (
	"fmt"
	"regexp"
	"strings"
)

type Summarizer interface {
	Replace(line string) string
}

type SceneSummarizer struct {
	Start     *regexp.Regexp
	Finish    *regexp.Regexp
	startLine string
	isActive  bool
}

func CreateSceneSummarizer() *SceneSummarizer {
	return &SceneSummarizer{
		Start:     regexp.MustCompile(`^\s*Loaded scene\s+'.*'\s*$`),
		Finish:    regexp.MustCompile(`^\s+Total Operation Time:\s+([\d.]+).+$`),
		startLine: "",
		isActive:  false,
	}
}

func (s *SceneSummarizer) Replace(line string) string {
	if s.Start.MatchString(line) {
		s.isActive = true
		s.startLine = strings.Replace(line, "\n", "", 1)
	}
	if s.isActive {
		groups := s.Finish.FindStringSubmatch(line)
		if groups != nil && len(groups) > 0 {
			totalOpTime := groups[1]
			s.isActive = false
			return fmt.Sprintf("%s [%sms]", s.startLine, totalOpTime)
		}
		return ""
	}
	return line
}
