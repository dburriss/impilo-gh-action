package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/olekukonko/tablewriter"
)

type Pair[T, U any] struct {
	First  T
	Second U
}

func mapSlice[T any, U any](xs []T, f func(T) U) []U {
	var result []U
	for _, x := range xs {
		result = append(result, f(x))
	}
	return result
}

func mapMap[K1 comparable, T1 any, K2 comparable, T2 any](m map[K1]T1, f func(K1, T1) (K2, T2)) map[K2]T2 {
	var result map[K2]T2
	for k1, v1 := range m {
		k2, v2 := f(k1, v1)
		result[k2] = v2
	}
	return result
}

func mapToSlice[K comparable, V any](m map[K]V) []Pair[K, V] {
	var result []Pair[K, V]
	for k1, v1 := range m {
		result = append(result, Pair[K, V]{k1, v1})
	}
	return result
}

func mapItemExists[K comparable, V any](key K, m map[K]V) bool {
	_, exists := m[key]
	return exists
}

func appendMap[K comparable, V any](m1 map[K]V, m2 map[K]V) {
	for k, _ := range m2 {
		if !mapItemExists(k, m1) {
			m1[k] = m2[k]
		}
	}
}

type Report interface {
	Run()
}

type ScanVulnerabilitiesStdOutReport struct {
	data []map[string]string
}

func (report ScanVulnerabilitiesStdOutReport) Run() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Package", "Severity", "Is Fix?", "Found in", "Fixed in"})

	for _, row := range report.data {
		var rowValues []string
		rowValues = append(rowValues, row["packageName"])
		rowValues = append(rowValues, row["severity"])
		rowValues = append(rowValues, row["fixAvailable"])
		rowValues = append(rowValues, row["foundIn"])
		rowValues = append(rowValues, row["fixedIn"])
		table.Append(rowValues)
	}
	table.Render()
}

type Command interface {
	Execute() []Report
}

type ScanVulnerabilitiesCommand struct {
	PackageManager  string
	TargetDirectory string
}

func trim(s string) string {
	return strings.Trim(s, " ")
}

func correctFoundIn(version string) string {
	version = trim(version)
	if strings.HasPrefix(version, ">=") {
		version = strings.Replace(version, ">=", "", 1)
	}
	return version
}

func correctFixedIn(version string) string {
	version = trim(version)
	if strings.HasPrefix(version, "<=") {
		version = strings.Replace(version, "=", "", 1)
	} else if strings.HasPrefix(version, "<") {
		version = strings.Replace(version, "<", "", 1)
	}
	return version
}

func parseNpmRange(rangeS string) map[string]string {
	r := make(map[string]string)
	ranges := strings.Split(rangeS, " - ")
	switch len(ranges) {
	case 1:
		// " " || "*"
		rangeS = trim(rangeS)
		if rangeS == "*" || rangeS == " " || rangeS == "" {
			r["foundIn"] = "0.0.0"
			r["fixedIn"] = ""
			break
		}

		ranges := strings.Split(rangeS, " ")
		switch len(ranges) {
		case 1:
			if strings.HasPrefix(rangeS, "<=") {
				r["foundIn"] = "0.0.0"
				r["fixedIn"] = ">" + trim(strings.Replace(ranges[0], "<=", "", 1))
			} else if strings.HasPrefix(rangeS, "<") {
				r["foundIn"] = "0.0.0"
				r["fixedIn"] = trim(strings.Replace(ranges[0], "<", "", 1))
			} else {
				r["foundIn"] = correctFoundIn(ranges[0])
				r["fixedIn"] = ""
			}
		case 2:
			// >=a.b.c <x.y.z
			r["foundIn"] = correctFoundIn(ranges[0])
			r["fixedIn"] = correctFixedIn(ranges[1])
		}

	case 2:
		// x - y
		r["foundIn"] = trim(ranges[0])
		r["fixedIn"] = ">" + trim(ranges[1])
	}
	return r
}

func parseNpmRanges(rangeS string) map[string]string {
	r := make(map[string]string)
	ranges := strings.Split(rangeS, "||")

	switch len(ranges) {
	case 0:
		r["foundIn"] = "0.0.0"
		r["fixedIn"] = ""
	case 1:
		rangeValues := parseNpmRange(rangeS)
		appendMap(r, rangeValues)
	default:
		regressions := mapSlice(ranges, parseNpmRanges)
		lastIndex := len(regressions) - 1
		first := regressions[0]
		last := regressions[lastIndex]
		foundIn, exists := first["foundIn"]
		if exists {
			r["foundIn"] = foundIn
		} else {
			r["foundIn"] = ""
		}
		fixedIn, exists := last["fixedIn"]
		if exists {
			r["fixedIn"] = fixedIn
		} else {
			r["fixedIn"] = ""
		}
	}

	return r
}

func parseNpmVulnerability(key string, jObj *gabs.Container) []map[string]string {
	result := []map[string]string{}
	var item = make(map[string]string)
	item["packageName"] = key
	severity, ok := jObj.Search("severity").Data().(string)
	if ok {
		item["severity"] = severity
	}
	fixAvailable, ok := jObj.Search("fixAvailable").Data().(bool)
	if ok {
		item["fixAvailable"] = strconv.FormatBool(fixAvailable)
	} else {
		fixObj := jObj.Search("fixAvailable").Data()
		if fixObj != nil {
			item["fixAvailable"] = "true"
		}
	}
	vRange, ok := jObj.Search("range").Data().(string)
	if ok {
		rangeM := parseNpmRanges(vRange)
		appendMap(item, rangeM)
	}
	// todo: duplicate if "regressedIn"
	result = append(result, item)
	return result
}

func parseNpmVulnerabilities(rawJson string) []map[string]string {
	var vulnerabilities []map[string]string
	jsonParsed, err := gabs.ParseJSON([]byte(rawJson))
	if err == nil {
		for key, jObj := range jsonParsed.Path("vulnerabilities").ChildrenMap() {
			items := parseNpmVulnerability(key, jObj)
			vulnerabilities = append(vulnerabilities, items...)
		}
	} else {
		fmt.Errorf(err.Error())
	}
	return vulnerabilities
}

func severityValue(severity string) int {
	switch severity {
	case "critical":
		return 4
	case "high":
		return 3
	case "moderate":
		return 2
	case "low":
		return 1
	default:
		return 1
	}
}

func (data ScanVulnerabilitiesCommand) Execute() []Report {
	fmt.Println("Package manager", data.PackageManager)
	fmt.Println("Target directory", data.TargetDirectory)
	cmdName := "npm"
	cmdArgs := []string{"audit", "--json", "--silent"}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = data.TargetDirectory
	var outbuf, errbuf strings.Builder // or bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	println(cmd.String())

	err := cmd.Run()
	stderr := errbuf.String()
	if len(stderr) > 0 && err != nil {
		log.Fatalf("Command %s failed.\n", cmdName)
		log.Fatalln(err)
		panic("Scan Vulnerabilities failed.")
	}

	var reports []Report
	vulnerabilities := parseNpmVulnerabilities(outbuf.String())
	sort.Slice(vulnerabilities, func(i, j int) bool {
		si := severityValue(vulnerabilities[i]["severity"])
		sj := severityValue(vulnerabilities[j]["severity"])
		if si == sj {
			return vulnerabilities[i]["packageName"] < vulnerabilities[j]["packageName"]
		}
		return severityValue(vulnerabilities[i]["severity"]) > severityValue(vulnerabilities[j]["severity"])
	})
	var stdoutReport Report = ScanVulnerabilitiesStdOutReport{
		data: vulnerabilities,
	}
	reports = append(reports, stdoutReport)
	return reports
}

func makeScanVulnerabilitiesCommand(config Config) []Command {
	commands := []Command{}
	for _, project := range config.Projects {
		var cmd Command = ScanVulnerabilitiesCommand{
			TargetDirectory: project.TargetDirectory,
			PackageManager:  project.PackageManager,
		}
		commands = append(commands, cmd)
	}

	return commands
}

func BuildCommands(config Config) []Command {
	commands := []Command{}

	if config.ScanVulnerabilities {
		commands = append(commands, makeScanVulnerabilitiesCommand(config)...)
	}

	if config.ScanLicenses {

	}

	return commands
}
