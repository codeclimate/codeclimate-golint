package main

import (
	"fmt"
	"github.com/codeclimate/cc-engine-go/engine"
	"github.com/golang/lint"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

const defaultMinConfidence = 0.8

func main() {
	rootPath := "/code/"
	analysisFiles, err := engine.GoFileWalk(rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing: %s", err)
		os.Exit(1)
	}

	config, err := engine.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
		os.Exit(1)
	}

	excludedFiles := getExcludedFiles(config)
	minConfidence := getMinConfidence(config)

	for _, path := range analysisFiles {
		relativePath := strings.SplitAfter(path, rootPath)[1]
		if isFileExcluded(relativePath, excludedFiles) {
			continue
		}

		lintFile(path, relativePath, minConfidence)
	}
}

func getExcludedFiles(config engine.Config) []string {
	excludedFiles := []string{}
	if config["exclude_paths"] != nil {
		for _, file := range config["exclude_paths"].([]interface{}) {
			excludedFiles = append(excludedFiles, file.(string))
		}
		sort.Strings(excludedFiles)
	}
	return excludedFiles
}

func isFileExcluded(filePath string, excludedFiles []string) bool {
	i := sort.SearchStrings(excludedFiles, filePath)
	return i < len(excludedFiles) && excludedFiles[i] == filePath
}

func lintFile(fullPath string, relativePath string, minConfidence float64) {
	linter := new(lint.Linter)

	code, err := ioutil.ReadFile(fullPath)
	if err != nil {
		warning := &engine.Warning{
			Description: "Could not read file",
			Location: &engine.Location{
				Path: fullPath,
				Lines: &engine.LinesOnlyPosition{
					Begin: 1,
					End:   1,
				},
			},
		}
		engine.PrintWarning(warning)
	}

	problems, err := linter.Lint("", code)
	if err != nil {
		warningDesc := fmt.Sprintf("Could not lint file (%s)", err)
		warning := &engine.Warning{
			Description: warningDesc,
			Location: &engine.Location{
				Path: fullPath,
			},
		}
		engine.PrintWarning(warning)
	}

	for _, problem := range problems {
		if problem.Confidence < minConfidence {
			continue
		}

		issue := &engine.Issue{
			Type:              "issue",
			Check:             codeClimateCheckName(&problem),
			Description:       (&problem).Text,
			RemediationPoints: 50000,
			Categories:        []string{"Style"},
			Location:          codeClimateLocation(&problem, relativePath),
		}
		engine.PrintIssue(issue)
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

func codeClimateLocation(l *lint.Problem, relativePath string) *engine.Location {
	position := l.Position

	return &engine.Location{
		Path: relativePath,
		Lines: &engine.LinesOnlyPosition{
			Begin: position.Line,
			End:   position.Line,
		},
	}
}

func getMinConfidence(config engine.Config) float64 {
	if subConfig, ok := config["config"].(map[string]interface{}); ok {
		if minConfidence, ok := subConfig["min_confidence"].(string); ok {
			val, err := strconv.ParseFloat(minConfidence, 64)
			if err == nil {
				return val
			}
		}
	}
	return defaultMinConfidence
}
