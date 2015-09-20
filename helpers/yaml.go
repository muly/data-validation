package helpers

// yaml related functions and types

import (
	//"errors"
	"io/ioutil"
	//"strings"

	"gopkg.in/yaml.v2"
)

type configRaw struct {
	ValidationRules []map[string]string //Note: fields must be exportable for the yaml unmarshal to work (so, starting with capital letter)
}

// the function parseYaml() reads the yaml file
//		and unmarshals into a struct with slice of map of string to string (config:value)
func parseYaml() (c configRaw, err error) {
	yamlData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return configRaw{}, err
	}

	err = yaml.Unmarshal(yamlData, &c)
	if err != nil {
		return configRaw{}, err
	}
	return c, nil

}
