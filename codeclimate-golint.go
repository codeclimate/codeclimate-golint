package main

import "github.com/codeclimate/cc-engine-go/engine"
import "github.com/golang/lint"
import "strings"
import "os"
import "io/ioutil"
import "sort"

func main() {
	rootPath := "/code/"
	analysisFiles, err := engine.GoFileWalk(rootPath)
	if err != nil {
		os.Exit(1)
	}

	config, err := engine.LoadConfig()
	if err != nil {
		os.Exit(1)
	}

	excludedFiles := []string{}
	if config["exclude_paths"] != nil {
		for _, file := range config["exclude_paths"].([]interface{}) {
			excludedFiles = append(excludedFiles, file.(string))
		}
		sort.Strings(excludedFiles)
	}

	for _, path := range analysisFiles {
		relativePath := strings.SplitAfter(path, rootPath)[1]
		i := sort.SearchStrings(excludedFiles, relativePath)
		if i < len(excludedFiles) && excludedFiles[i] == relativePath {
			continue
		}

		linter := new(lint.Linter)

		code, err := ioutil.ReadFile(path)
		if err != nil {
			warning := &engine.Warning{
				Description: "Could not read file",
				Location: &engine.Location{
					Path: path,
					Lines: &engine.Position{
						Begin: 1,
						End:   1,
					},
				},
			}
			engine.PrintWarning(warning)
		}

		problems, err := linter.Lint("", code)
		if err != nil {
			warning := &engine.Warning{
				Description: "Could not lint file",
				Location: &engine.Location{
					Path: path,
				},
			}
			engine.PrintWarning(warning)
		}

		for _, problem := range problems {
			issue := &engine.Issue{
				Type:              "issue",
				Check:             codeClimateCheckName(&problem),
				Description:       (&problem).Text,
				RemediationPoints: 500,
				Categories:        []string{"Style"},
				Location:          codeClimateLocation(&problem, path, rootPath),
			}
			engine.PrintIssue(issue)
		}
	}

}

func codeClimateCheckName(l *lint.Problem) string {
	return "GoLint/" + strings.Title(l.Category) + "/" + codeClimateIssueName(l)
}

func codeClimateIssueName(l *lint.Problem) string {
	splitURL := strings.Split(l.Link, "#")
	kebabLink := splitURL[len(splitURL)-1]
	splitLink := strings.Split(kebabLink, "-")
	for i, link := range splitLink {
		splitLink[i] = strings.Title(link)
	}
	camelLink := strings.Join(splitLink, "")
	return camelLink
}

func codeClimateLocation(l *lint.Problem, path string, rootPath string) *engine.Location {
	position := l.Position

	return &engine.Location{
		Path: strings.SplitAfter(path, rootPath)[1],
		Lines: &engine.Position{
			Begin: position.Line,
			End:   position.Line,
		},
	}
}
