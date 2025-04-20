package parser

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type Summarizer interface {
	Replace(line string) (string, bool)
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

func (s *SceneSummarizer) Replace(line string) (string, bool) {
	if s.Start.MatchString(line) {
		s.isActive = true
		s.startLine = strings.Replace(line, "\n", "", 1)
	}
	if s.isActive {
		groups := s.Finish.FindStringSubmatch(line)
		if groups != nil && len(groups) > 0 {
			totalOpTime := groups[1]
			s.isActive = false
			return fmt.Sprintf("%s [%sms]", s.startLine, totalOpTime), false
		}
		return "", false
	}
	return line, false
}

func AssetCountSummarizer() *CountSummarizer {
	return &CountSummarizer{
		countRegexes: []*regexp.Regexp{
			regexp.MustCompile(`\s*Start importing.*`),
			regexp.MustCompile(`\s*(\[Worker\s?\w+\])\s*Start importing.*`),
		},
		finishRegex:         regexp.MustCompile(`^Asset Pipeline Refresh: Total: .+ seconds - Initiated by .+$`),
		count:               0,
		multiplicativePrint: true,
		multiplicativeFrom:  10,
		multiplicativeBase:  100,
		msgTemplate:         "Imported %d Assets",
	}
}

type CountSummarizer struct {
	countRegexes        []*regexp.Regexp
	finishRegex         *regexp.Regexp
	count               int
	multiplicativePrint bool
	multiplicativeFrom  int
	multiplicativeBase  int
	msgTemplate         string
}

func (c *CountSummarizer) Replace(line string) (string, bool) {
	match := false
	for _, r := range c.countRegexes {
		if r.MatchString(line) {
			c.count = c.count + 1
			match = true
			break
		}
	}
	if c.finishRegex.MatchString(line) {
		finalCountMsg := c.message()
		c.count = 0
		return finalCountMsg + "\n" + line, false
	}

	if !match {
		return line, false
	}
	if match && !c.multiplicativePrint {
		return c.message(), false
	}
	if match && c.multiplicativePrint && c.count < c.multiplicativeFrom {
		return c.message(), false
	}
	if match && c.multiplicativePrint && c.count == c.multiplicativeFrom {
		return c.message(), false
	}
	if match && c.multiplicativePrint && math.Remainder(float64(c.count), float64(c.multiplicativeBase)) == 0 {
		return c.message(), false
	}
	//match == true here, skip this line
	return "", true
}

func (c *CountSummarizer) message() string {
	return fmt.Sprintf(c.msgTemplate, c.count)
}

//^\[.+\] Csc.+
//type
