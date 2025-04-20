package command

import (
	"fmt"
	"github.com/joaquimmnetto/unity-log-processor/parser"
	"io"
)

func PrintConfig(config parser.Config, writer io.Writer) error {
	yamlData, err := config.AsYaml()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, string(yamlData[:]))
	return err
}
