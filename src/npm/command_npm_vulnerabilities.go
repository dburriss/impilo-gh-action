package npm

import (
	//"fmt"
	"dburriss/impilo_gh/domain"
	"dburriss/impilo_gh/maps"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/sethvargo/go-githubactions"
)

// ====================
// Report
// ====================
type NpmScanVulnerabilitiesStdOutReport struct {
	title           string
	vulnerabilities []map[string]string
	summary         map[string]int64
}

type NpmScanVulnerabilitiesGithubReport struct {
	title           string
	vulnerabilities []map[string]string
	summary         map[string]int64
}

func (report NpmScanVulnerabilitiesStdOutReport) Title() string {
	if report.title == "" {
		return "Npm Scan Vulnerabilities StdOut Report"
	}
	return report.title
}

func (report NpmScanVulnerabilitiesGithubReport) Title() string {
	if report.title == "" {
		return "Npm Scan Vulnerabilities Github Report"
	}
	return report.title
}

func (report NpmScanVulnerabilitiesStdOutReport) Run() {
	vulTable := tablewriter.NewWriter(os.Stdout)
	vulTable.SetHeader([]string{"Package", "Severity", "Is Fix?", "Found in", "Fixed in"})

	for _, row := range report.vulnerabilities {
		var rowValues []string
		rowValues = append(rowValues, row["packageName"])
		rowValues = append(rowValues, row["severity"])
		rowValues = append(rowValues, row["fixAvailable"])
		rowValues = append(rowValues, row["foundIn"])
		rowValues = append(rowValues, row["fixedIn"])
		vulTable.Append(rowValues)
	}
	vulTable.Render()

	summaryTable := tablewriter.NewWriter(os.Stdout)
	summaryTable.SetHeader([]string{"SUMMARY"})
	total := "0"
	for k, v := range report.summary {
		if k == "total" {
			total = strconv.FormatInt(v, 10)
			continue
		}
		var rowValues []string
		rowValues = append(rowValues, strings.ToUpper(k))
		rowValues = append(rowValues, strconv.FormatInt(v, 10))
		summaryTable.Append(rowValues)
	}
	summaryTable.SetFooter([]string{"TOTAL", total})
	summaryTable.Render()
}

func (report NpmScanVulnerabilitiesGithubReport) Run() {
	a := githubactions.New()
	c, err := a.Context()
	if err == nil && c.RunID > 0 {
		a.Group(report.Title())

		writer := new(strings.Builder)
		table := tablewriter.NewWriter(writer)
		table.SetHeader([]string{"SEVERITY"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		total := "0"
		for k, v := range report.summary {
			if k == "total" {
				total = strconv.FormatInt(v, 10)
				continue
			}
			var rowValues []string
			rowValues = append(rowValues, k)
			rowValues = append(rowValues, strconv.FormatInt(v, 10))
			table.Append(rowValues)
		}
		table.SetFooter([]string{"Total", total})
		table.Render()
		githubactions.AddStepSummary(writer.String())
		a.EndGroup()
	}
}

// ====================
// Command
// ====================
type ScanNpmVulnerabilitiesCommand struct {
	PackageManager  string
	TargetDirectory string
	Load            func() string
	title           string
}

// Title implements Command
func (cmd ScanNpmVulnerabilitiesCommand) Title() string {
	if cmd.title == "" {
		return "Scan NPM Vulnerabilities Command"
	}
	return cmd.title
}

func (cmd ScanNpmVulnerabilitiesCommand) Execute() []domain.Report {
	// create and run command
	if cmd.Load == nil {
		cmd.Load = func() string { return runNpmAudit(cmd.TargetDirectory) }
	}
	result := cmd.Load()
	reports := handleNpmAudit(result)
	return reports
}

func runNpmAudit(targetDirectory string) string {
	cliCmdName := "npm"
	cliCmdArgs := []string{"audit", "--json", "--silent"}

	cliCmd := exec.Command(cliCmdName, cliCmdArgs...)
	cliCmd.Dir = targetDirectory
	var outbuf, errbuf strings.Builder // or bytes.Buffer
	cliCmd.Stdout = &outbuf
	cliCmd.Stderr = &errbuf
	println(cliCmd.String())

	// actually run the command and wait till done
	err := cliCmd.Run()

	// get output and sort it
	stderr := errbuf.String()
	stdout := outbuf.String()
	if len(stdout) == 0 || (len(stderr) > 0 && err != nil) {
		log.Fatalf("Command %s failed.\n", cliCmdName)
		log.Fatalln(err)
		panic("Scan NPM Vulnerabilities failed.")
	}
	return stdout
}

func handleNpmAudit(auditJson string) []domain.Report {
	// parse JSON output
	jsonParsed, err := gabs.ParseJSON([]byte(auditJson))
	var vulnerabilities []map[string]string
	var summary map[string]int64

	if err == nil {
		// vulnerabilities
		vulnerabilities = parseNpmVulnerabilities(jsonParsed)
		sort.Slice(vulnerabilities, func(i, j int) bool {
			si := severityValue(vulnerabilities[i]["severity"])
			sj := severityValue(vulnerabilities[j]["severity"])
			if si == sj {
				return vulnerabilities[i]["packageName"] < vulnerabilities[j]["packageName"]
			}
			return severityValue(vulnerabilities[i]["severity"]) > severityValue(vulnerabilities[j]["severity"])
		})
		// vulnerabilities summary
		summary = parseNpmVulnerabilitiesSummary(jsonParsed)

	} else {
		fmt.Printf("Parsing audit JSON failed: %s. ", err.Error())
	}
	// create and return reports
	var reports []domain.Report
	var stdoutReport domain.Report = NpmScanVulnerabilitiesStdOutReport{
		vulnerabilities: vulnerabilities,
		summary:         summary,
	}
	var githubReport domain.Report = NpmScanVulnerabilitiesGithubReport{
		vulnerabilities: vulnerabilities,
		summary:         summary,
	}
	reports = append(reports, stdoutReport, githubReport)
	return reports
}

func MakeScanVulnerabilitiesCommand(config domain.Config) []domain.Command {
	commands := []domain.Command{}
	for _, project := range config.Projects {
		if project.PackageManager == "npm" {
			var cmd domain.Command = ScanNpmVulnerabilitiesCommand{
				TargetDirectory: project.TargetDirectory,
				PackageManager:  project.PackageManager,
			}
			commands = append(commands, cmd)
		}
	}

	return commands
}

// ====================
// Helpers
// ====================
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

func parseNpmVulnerabilities(jsonParsed *gabs.Container) []map[string]string {
	var vulnerabilities []map[string]string

	for key, jObj := range jsonParsed.Path("vulnerabilities").ChildrenMap() {
		items := parseNpmVulnerability(key, jObj)
		vulnerabilities = append(vulnerabilities, items...)
	}

	return vulnerabilities
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
		maps.Append(item, rangeM)
	}
	// todo: duplicate if "regressedIn"
	result = append(result, item)
	return result
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
		maps.Append(r, rangeValues)
	default:
		regressions := maps.MapSlice(ranges, parseNpmRanges)
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

func parseNpmVulnerabilitiesSummary(jsonParsed *gabs.Container) map[string]int64 {
	summary := make(map[string]int64)

	for key, v := range jsonParsed.Search("metadata", "vulnerabilities").ChildrenMap() {
		count := v.Data()
		n, ok := count.(float64)
		summary[key] = 0
		if ok {
			summary[key] = int64(n)
		}
	}

	return summary
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
	case "info":
		return 0
	default:
		return 0
	}
}
