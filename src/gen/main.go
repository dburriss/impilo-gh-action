//go:build ignore

package main

import (
	"dburriss/impilo_gh/gen"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Input struct {
	Description string
	Required    bool
	Default     string
	Type        string
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
	ArgName      string
	FieldName    string
	VarName      string
	DefaultValue string
	EnvKey       string
	OptsType     string
	Type         string
	Description  string
}

type TemplateData struct {
	Items []TemplateItem
	Args  []string
}

func descAndDefault(input Input, t string) string {
	if input.Default == "" || t == "bool" {
		return input.Description
	}
	return input.Description + " Default: " + input.Default
}

func main() {
	// read yaml
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
		_, isBool := gen.AsBool(v.Default, v.Type)
		if isBool {
			t = "bool"
		}
		item := TemplateItem{
			Index:        i,
			ArgName:      k,
			FieldName:    gen.CamelCase(k),
			VarName:      varName,
			DefaultValue: gen.FormatT(t, v.Default),
			EnvKey:       "INPUT_" + strings.ToUpper(k),
			Type:         t,
			OptsType:     t,
			Description:  descAndDefault(v, t),
		}

		if t == "string" && strings.HasSuffix(k, "file") {
			item.OptsType = "flags.Filename"
		}

		items = append(items, item)
		i++
	}

	// create template
	template, err := template.ParseFiles("./gen/action_template.go.tpl")

	if err != nil {
		panic(err)
	}

	// create the src file to write template to
	srcFile, err := os.Create("./domain/input_types_generated.go")
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
