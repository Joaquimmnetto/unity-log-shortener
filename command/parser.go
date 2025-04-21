package command

import (
	"bufio"
	"fmt"
	"github.com/joaquimmnetto/unity-log-processor/parser"
	"io"
	"os"
	"regexp"
	"strings"
)

func NewParserWriter(printUnparsed bool) *ParserWriters {
	return &ParserWriters{
		mainWriter:      os.Stdout,
		secondaryWriter: nil,
		printUnparsed:   printUnparsed,
	}
}

func ParserWriterWithSecondaryLogFile(secondaryLogFile string, printUnparsed bool) *ParserWriters {
	secondaryWriter, err := os.Create(secondaryLogFile)
	if err != nil {
		panic(fmt.Errorf("cant open specified secondary file (%s): %w", secondaryLogFile, err))
	}
	return &ParserWriters{
		mainWriter:      os.Stdout,
		secondaryWriter: secondaryWriter,
		printUnparsed:   printUnparsed,
	}
}

type ParserWriters struct {
	mainWriter      io.StringWriter
	secondaryWriter io.StringWriter
	printUnparsed   bool
}

func (p *ParserWriters) ParseWholeInput(input *bufio.Scanner, config *parser.Config) error {
	for input.Scan() {
		line := input.Text()
		if p.printUnparsed {
			mustWriteLine(p.mainWriter, line)
		} else {
			mustWriteLine(p.secondaryWriter, line)
		}
		parsedLine, skip := parseLine(config, line)
		if skip {
			continue
		}
		if p.printUnparsed {
			mustWriteLine(p.secondaryWriter, parsedLine)
		} else {
			mustWriteLine(p.mainWriter, parsedLine)
		}
	}
	if err := input.Err(); err != nil {
		return err
	}
	return nil
}

func parseLine(config *parser.Config, line string) (string, bool) {
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
		return "", true
	}
	for _, summarizer := range config.Summarizers.AllSummarizers() {
		line, skipLine = summarizer.Replace(line)
		if skipLine {
			return "", true
		}
	}
	return line, false
}

func matchAnyRegex(line string, matchers []*regexp.Regexp) bool {
	for _, matcher := range matchers {
		if matcher.MatchString(line) {
			return true
		}
	}
	return false
}

func matchAnyMatcher(line string, matchers []parser.Matcher) bool {
	for _, matcher := range matchers {
		if matcher.Match(line) {
			return true
		}
	}
	return false
}

func mustWriteLine(writer io.StringWriter, line string) {
	if writer == nil {
		return
	}
	_, err := writer.WriteString(line + "\n")
	if err != nil {
		panic(err)
	}
}
