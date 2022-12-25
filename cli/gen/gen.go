//go:build ignore

package main

import (
	"dburriss/impilo_gh/gen"
	"io/ioutil"
	"log"
	"os"
	"text/template"

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
	// fmt.Println("YAML: ", actionYaml)
	// map action yaml to template data

	// get args from runs.args
	var args []string
	for _, v := range actionYaml.Runs.Args {
		arg, matched := gen.StripArg(v)
		if matched {
			args = append(args, arg)
		}
	}

	var items []TemplateItem
	for i, k := range args {
		varName := gen.Normalize(k)
		v := actionYaml.Inputs[k]
		t := "string"
		_, isBool := gen.AsBool(v.Default)
		if isBool {
			t = "bool"
		}
		item := TemplateItem{
			Index:        i,
			FieldName:    gen.CamelCase(k),
			VarName:      varName,
			DefaultValue: gen.FormatT(t, v.Default),
			Type:         t,
		}
		items = append(items, item)
		i++
	}

	// create template
	template, err := template.ParseFiles("./gen/action.template.txt")

	if err != nil {
		panic(err)
	}

	// create the src file to write template to
	srcFile, err := os.Create("./input.generated.go")
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
