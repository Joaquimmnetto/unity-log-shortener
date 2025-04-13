package main

import (
	"bufio"
	"fmt"
	v0 "github.com/joaquimmnetto/unity-log-processor/v0"
	"os"
	"regexp"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := v0.DefaultConfig()

	for scanner.Scan() {
		line := scanner.Text()
		for _, matcher := range config.Preprocessors.FirstMatchInlineRegexes() {
			found := matcher.FindString(line)
			if found != "" {
				line = strings.Replace(line, found, "", 1)
			}
		}

		for _, matcher := range config.Preprocessors.AllMatchInLineRegexes() {
			line = matcher.ReplaceAllString(line, "")
		}

		skipLine := false
		skipLine = matchAnyRegex(line, config.Matchers.WholeLineRegexes())
		if !skipLine {
			skipLine = matchAnyMatcher(line, config.Matchers.AllMatchers())
		}
		if skipLine {
			continue
		}
		for _, summarizer := range config.Summarizers.AllSummarizers() {
			line = summarizer.Replace(line)
		}
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func matchAnyRegex(line string, matchers []*regexp.Regexp) bool {
	for _, matcher := range matchers {
		if matcher.MatchString(line) {
			return true
		}
	}
	return false
}

func matchAnyMatcher(line string, matchers []v0.Matcher) bool {
	for _, matcher := range matchers {
		if matcher.Match(line) {
			return true
		}
	}
	return false
}
