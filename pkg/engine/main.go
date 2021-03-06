package engine

import (
	"fmt"

	"github.com/olblak/updateCli/pkg/config"

	"path/filepath"
	"strings"
)

var engine Engine

// Engine defined parameters for a specific engine run
type Engine struct {
	conf    config.Config
	Options Options
}

// Apply run the full process one yaml file
func (e *Engine) Apply(cfgFile string) error {

	_, basename := filepath.Split(cfgFile)
	cfgFileName := strings.TrimSuffix(basename, filepath.Ext(basename))

	fmt.Printf("\n\n%s\n", strings.Repeat("#", len(cfgFileName)+4))
	fmt.Printf("# %s #\n", strings.ToTitle(cfgFileName))
	fmt.Printf("%s\n\n", strings.Repeat("#", len(cfgFileName)+4))

	e.conf.ReadFile(cfgFile, e.Options.ValuesFile)

	source, err := e.conf.Source.Execute()

	if err != nil {
		return err
	}

	if source == "" {
		fmt.Printf("\n\u26A0 No value returned from Source, nothing else to do")
		return nil
	}

	if len(e.conf.Conditions) > 0 {
		ok, err := e.conditions(source)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}
	}

	if len(e.conf.Targets) > 0 {
		e.targets(source)
	}

	return nil
}

// conditions iterates on every conditions and test the result
func (e *Engine) conditions(source string) (bool, error) {

	fmt.Printf("\n\n%s:\n", strings.ToTitle("conditions"))
	fmt.Printf("%s\n\n", strings.Repeat("=", len("conditions")+1))

	for _, c := range e.conf.Conditions {
		ok, err := c.Execute(source)
		if err != nil {
			return false, err
		}

		if !ok {
			fmt.Printf("\n\u26A0 skipping: condition not met")
			ok = false
			return false, nil
		}

	}

	return true, nil
}

// targets iterate on every targets and then call target on each of them
func (e *Engine) targets(source string) error {

	fmt.Printf("\n\n%s:\n", strings.ToTitle("Targets"))
	fmt.Printf("%s\n\n", strings.Repeat("=", len("Targets")+1))

	for id, t := range e.conf.Targets {
		err := t.Execute(source, &e.Options.Target)
		if err != nil {
			fmt.Printf("Something went wrong in target \"%v\" :\n", id)
			fmt.Printf("%v\n\n", err)
		}
	}
	return nil
}

// Show displays the configuration that should be apply
func (e *Engine) Show(cfgFile string) error {

	_, basename := filepath.Split(cfgFile)
	cfgFileName := strings.TrimSuffix(basename, filepath.Ext(basename))

	fmt.Printf("\n\n%s\n", strings.Repeat("#", len(cfgFileName)+4))
	fmt.Printf("# %s #\n", strings.ToTitle(cfgFileName))
	fmt.Printf("%s\n\n", strings.Repeat("#", len(cfgFileName)+4))

	e.conf.ReadFile(cfgFile, e.Options.ValuesFile)
	err := e.conf.Display()
	if err != nil {
		return err
	}

	return nil
}
