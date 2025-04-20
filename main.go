package main

import (
	"bufio"
	"flag"
	"github.com/joaquimmnetto/unity-log-processor/command"
	"github.com/joaquimmnetto/unity-log-processor/parser"
	"os"
)

func main() {

	var configFilePath string
	flag.StringVar(&configFilePath, "parseConfig", "log_parser.yaml", "path to the yaml parsing configuration file. Defaults log_parser.yaml or internal configuration if no such file is present")
	var secondaryLogFile string
	flag.StringVar(&secondaryLogFile, "secondaryLogFile", "", "path to secondary log file. Secondary log file is the non-parsed input by default, or parsed input if -printUnparsed is set")
	var printUnparsed bool
	flag.BoolVar(&printUnparsed, "printUnparsed", false, "prints unparsed input to stdout. Parsed input will be sent to -secondaryLogFile, if present")
	var printDefaultConfig bool
	flag.BoolVar(&printDefaultConfig, "printDefaultConfig", false, "prints the default configuration as yaml to stdout and then finishes the application")
	var help bool
	flag.BoolVar(&help, "help", false, "show help message")

	flag.Parse()

	config := loadConfig(configFilePath)
	scanner := bufio.NewScanner(os.Stdin)

	if printDefaultConfig {
		err := command.PrintConfig(parser.DefaultConfig(), os.Stdout)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	if help {
		flag.Usage()
		os.Exit(0)
	}

	parserWriter := command.NewParserWriter(printUnparsed)
	if secondaryLogFile != "" {
		parserWriter = command.ParserWriterWithSecondaryLogFile(secondaryLogFile, printUnparsed)
	}
	err := parserWriter.ParseWholeInput(scanner, &config)
	if err = scanner.Err(); err != nil {
		panic(err)
	}
}

func loadConfig(configFilePath string) parser.Config {
	if _, err := os.Stat(configFilePath); err == nil {
		config, err := parser.LoadConfigFromYamlFile(configFilePath)
		if err != nil {
			panic(err)
		}
		return config
	} else {
		return parser.DefaultConfig()
	}
}
