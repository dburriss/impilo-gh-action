package input

import (
	"fmt"
	"regexp"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// for variable name and filed name
func Normalize(str string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "")
}

// for filed name
func Title(str string) string {
	return cases.Title(language.English).String(str)
}

func FormatT(t string, v string) string {
	switch t {
	case "bool":
		return v
	default:
		return fmt.Sprintf("\"%s\"", v)
	}
}

func StripArg(arg string) (string, bool) {
	//b := []byte(arg)
	ms := regexp.MustCompile(`\${{\sinputs\.([\w-]+)\s}}`).FindStringSubmatch(arg)
	if len(ms) > 1 {
		return ms[1], true
	}
	return "", false
}
