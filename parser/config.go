package parser

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

func LoadConfigFromYamlFile(filepath string) (Config, error) {
	var config Config
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("error reading file %s: %w", filepath, err)
	}
	err = yaml.Unmarshal(fileData, &config)
	if err != nil {
		err = fmt.Errorf("error loading configuration from yaml file %s: %w", filepath, err)
	}
	return config, err
}

type Config struct {
	Preprocessors Preprocessors `json:"preprocessors" yaml:"preprocessors"`
	Matchers      Matchers      `json:"matchers" yaml:"matchers"`
	Summarizers   Summarizers   `json:"summarizers" yaml:"summarizers"`
}

type Preprocessors struct {
	RemoveFirstMatchingFromLine []string `json:"removeFirstMatchingFromLine" yaml:"removeFirstMatchingFromLine"`
	RemoveAllMatchingFromLine   []string `json:"removeAllMatchingFromLine" yaml:"removeAllMatchingFromLine"`
}

func (m Preprocessors) FirstMatchInlineRegexes() []*regexp.Regexp {
	return stringsToRegexes(m.RemoveFirstMatchingFromLine)
}

func (m Preprocessors) AllMatchInLineRegexes() []*regexp.Regexp {
	return stringsToRegexes(m.RemoveAllMatchingFromLine)
}

type Matchers struct {
	RemoveLine            []string                  `json:"removeLine" yaml:"removeLine"`
	RemoveTabulatedBlocks map[string]TabulatedBlock `json:"removeTabulatedBlocks" yaml:"removeTabulatedBlocks"`
	RemoveStartEndBlocks  map[string]StartEndBlock  `json:"removeStartEndBlocks" yaml:"removeStartEndBlocks"`

	matchers          []Matcher
	removeLineRegexes []*regexp.Regexp
}

func (m *Matchers) WholeLineRegexes() []*regexp.Regexp {
	if m.removeLineRegexes == nil {
		m.removeLineRegexes = stringsToRegexes(m.RemoveLine)
	}
	return m.removeLineRegexes
}

func stringsToRegexes(strs []string) []*regexp.Regexp {
	result := make([]*regexp.Regexp, 0, len(strs))
	for _, regex := range strs {
		result = append(result, regexp.MustCompile(regex))
	}
	return result
}

func (m *Matchers) AllMatchers() []Matcher {
	if m.matchers == nil {
		m.matchers = make([]Matcher, 0, len(m.RemoveStartEndBlocks))
		for _, tabBlock := range m.RemoveTabulatedBlocks {
			m.matchers = append(m.matchers, tabBlock.AsMatcher())
		}
		for _, seBlock := range m.RemoveStartEndBlocks {
			m.matchers = append(m.matchers, seBlock.AsMatcher())
		}
	}
	return m.matchers
}

type TabulatedBlock struct {
	Start      string `json:"start" yaml:"start"`
	MatchStart bool   `json:"matchStart"  yaml:"matchStart" default:"true"`
}

func (t TabulatedBlock) AsMatcher() *TabulatedLineMatcher {
	return CreateTabulatedLineMatcher(t.Start, t.MatchStart)
}

type StartEndBlock struct {
	Start      string `json:"start" yaml:"start"`
	End        string `json:"end" yaml:"end"`
	MatchStart bool   `json:"matchStart"  yaml:"matchStart" default:"true"`
	MatchEnd   bool   `json:"matchEnd"  yaml:"matchEnd" default:"false"`
}

func (t StartEndBlock) AsMatcher() *StartEndBlockMatcher {
	return CreateStartEndBlockMatcher(t.Start, t.End, t.MatchStart, t.MatchEnd)
}

type Summarizers struct {
	EnableSceneSummarizer      bool `json:"enableSceneSummarizer" yaml:"enableSceneSummarizer"`
	EnableAssetsSummarizer     bool `json:"enableAssetsSummarizer" yaml:"enableAssetsSummarizer"`
	EnableCscWarningsSumarizer bool `json:"enableCscWarningsSumarizer" yaml:"enableCscWarningsSumarizer"`

	summarizers []Summarizer
}

func (s *Summarizers) AllSummarizers() []Summarizer {
	if s.summarizers == nil {
		s.summarizers = make([]Summarizer, 0, 1)
		if s.EnableSceneSummarizer {
			s.summarizers = append(s.summarizers, CreateSceneSummarizer())
		}
		if s.EnableAssetsSummarizer {
			s.summarizers = append(s.summarizers, AssetCountSummarizer())
		}
		if s.EnableCscWarningsSumarizer {
			s.summarizers = append(s.summarizers, CscWarningsCountSummarizer())
		}
	}
	return s.summarizers
}

func (c Config) AsJson() ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error parsing configuration into json: %w", err)
	}
	return jsonBytes, err
}

func (c Config) ToJsonFile(filepath string) error {
	jsonBytes, err := c.AsJson()
	if err != nil {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error opening configuration json file for writting: %w", err)
	}
	_, err = file.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("error writting configuration json file: %w", err)
	}
	return nil
}

func (c Config) AsYaml() ([]byte, error) {
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("error parsing configuration into yaml: %w", err)
	}
	return yamlBytes, err
}

func (c Config) ToYamlFile(filepath string) error {
	yamlBytes, err := c.AsYaml()
	if err != nil {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error opening configuration yaml file for writting: %w", err)
	}
	_, err = file.Write(yamlBytes)
	if err != nil {
		return fmt.Errorf("error writting configuration yaml file: %w", err)
	}
	return nil
}
