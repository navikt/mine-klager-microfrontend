package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	classAttrPattern        = regexp.MustCompile(`class="([^"]*)"`)
	cssClassSelectorPattern = regexp.MustCompile(`\.([a-zA-Z_][\w-]*)`)
)

// ComponentCssFile represents a pre-read and indexed component CSS file.
type ComponentCssFile struct {
	Name    string
	Content string
	Classes map[string]struct{}
}

// ComponentMatch records which CSS file was included and which classes caused it to match.
type ComponentMatch struct {
	File           string
	MatchedClasses []string
}

// extractClassesFromHTML extracts all CSS class names from an HTML string by parsing class="..." attributes.
func extractClassesFromHTML(html string) map[string]struct{} {
	classes := make(map[string]struct{})
	for _, match := range classAttrPattern.FindAllStringSubmatch(html, -1) {
		if len(match) > 1 {
			for _, className := range strings.Fields(match[1]) {
				if className != "" {
					classes[className] = struct{}{}
				}
			}
		}
	}
	return classes
}

// extractClassesFromCSS extracts all class names referenced as selectors in a CSS string.
func extractClassesFromCSS(css string) map[string]struct{} {
	classes := make(map[string]struct{})
	for _, match := range cssClassSelectorPattern.FindAllStringSubmatch(css, -1) {
		if len(match) > 1 {
			classes[match[1]] = struct{}{}
		}
	}
	return classes
}

// ReadComponentCssFiles reads all *.min.css files from the given directory and returns them indexed.
func ReadComponentCssFiles(dir string) ([]ComponentCssFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".min.css") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	files := make([]ComponentCssFile, 0, len(names))
	for _, name := range names {
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}

		files = append(files, ComponentCssFile{
			Name:    name,
			Content: string(content),
			Classes: extractClassesFromCSS(string(content)),
		})
	}

	return files, nil
}

// FindComponentCss finds all component CSS files whose selectors reference at least one class
// present in the given HTML strings. Returns the matched file names, their concatenated CSS content,
// and detailed match information showing which classes caused each file to be included.
func FindComponentCss(componentFiles []ComponentCssFile, htmlStrings ...string) (files []string, css string, matches []ComponentMatch) {
	htmlClasses := make(map[string]struct{})
	for _, html := range htmlStrings {
		for cls := range extractClassesFromHTML(html) {
			htmlClasses[cls] = struct{}{}
		}
	}

	var matchedFiles []string
	var matchedCSS []string
	var matchedDetails []ComponentMatch

	for _, file := range componentFiles {
		var overlapping []string
		for cls := range file.Classes {
			if _, ok := htmlClasses[cls]; ok {
				overlapping = append(overlapping, cls)
			}
		}

		if len(overlapping) > 0 {
			sort.Strings(overlapping)
			matchedFiles = append(matchedFiles, file.Name)
			matchedCSS = append(matchedCSS, file.Content)
			matchedDetails = append(matchedDetails, ComponentMatch{
				File:           file.Name,
				MatchedClasses: overlapping,
			})
		}
	}

	return matchedFiles, strings.Join(matchedCSS, "\n"), matchedDetails
}
