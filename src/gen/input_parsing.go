package gen

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// for variable name and filed name
func Normalize(str string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "")
}

// for filed name
func CamelCase(str string) string {
	var titled []string

	words := strings.Split(str, "-")

	for _, value := range words {
		titled = append(titled, cases.Title(language.English).String(Normalize(value)))
	}
	return strings.Join(titled, "")
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

// parses the string to a boolean returning (value, isBoolean)
func AsBool(defaultValue string, typeValue string) (bool, bool) {
	if typeValue == "boolean" {
		switch strings.ToLower(defaultValue) {
		case "true":
			return true, true
		case "false":
			return false, true
		default:
			return false, true
		}
	}

	switch strings.ToLower(defaultValue) {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return false, false
	}
}
