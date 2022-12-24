//go:build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"dburriss/impilo_gh/input"

	"gopkg.in/yaml.v2"
)

type Input struct {
	Description string
	Required    bool
	Default     string
}

type Runs struct {
	Args []string
}

// represents inputs needed from action.yml
type Action struct {
	Inputs map[string]Input
	Runs
}

// template data
type TemplateItem struct {
	Index        int
	FieldName    string
	VarName      string
	DefaultValue string
	Type         string
}

type TemplateData struct {
	Items []TemplateItem
	Args  []string
}

func main() {
	// read json
	content, err := ioutil.ReadFile("../action.yml")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `actionYaml`
	var actionYaml Action
	err = yaml.Unmarshal(content, &actionYaml)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	fmt.Println("YAML: ", actionYaml)
	// map action yaml to template data

	// get args from runs.args
	var args []string
	for _, v := range actionYaml.Runs.Args {
		arg, matched := input.StripArg(v)
		if matched {
			args = append(args, arg)
		}
	}

	var items []TemplateItem
	i := 0
	for k, v := range actionYaml.Inputs {
		varName := input.Normalize(k)
		item := TemplateItem{
			Index:        i,
			FieldName:    input.Title(varName),
			VarName:      varName,
			DefaultValue: input.FormatT("string", v.Default),
			Type:         "string",
		}
		items = append(items, item)
		i++
	}

	// create template
	template, err := template.ParseFiles("./input/action.template.txt")

	if err != nil {
		panic(err)
	}

	// create the src file to write template to
	srcFile, err := os.Create("./input/action.generated.go")
	if err != nil {
		log.Println("file creation ERROR: ", err)
		return
	}

	// write template to standard out
	data := TemplateData{Items: items}
	err = template.Execute(os.Stdout, data)
	// write template to src file
	err = template.Execute(srcFile, data)
	if err != nil {
		log.Print("Template ERROR: ", err)
		return
	}

	// close the src file
	srcFile.Close()
}
