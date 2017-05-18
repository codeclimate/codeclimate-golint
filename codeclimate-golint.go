package main

import (
	"fmt"
	"github.com/codeclimate/cc-engine-go/engine"
	"github.com/golang/lint"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const defaultMinConfidence = 0.8

func main() {
	rootPath := "/code/"

	config, err := engine.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
		os.Exit(1)
	}

	analysisFiles, err := engine.GoFileWalk(rootPath, engine.IncludePaths(rootPath, config))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing: %s", err)
		os.Exit(1)
	}

	minConfidence := getMinConfidence(config)

	for _, path := range analysisFiles {
		relativePath := strings.SplitAfter(path, rootPath)[1]

		lintFile(path, relativePath, minConfidence)
	}
}

func lintFile(fullPath string, relativePath string, minConfidence float64) {
	linter := new(lint.Linter)

	code, err := ioutil.ReadFile(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %v\n", fullPath)
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}

	problems, err := linter.Lint(fullPath, code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not lint file: %v\n", fullPath)
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
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
