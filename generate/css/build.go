package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BuildResult contains the final CSS and statistics about the build.
type BuildResult struct {
	CSS                    string
	TokensBefore           int
	TokensAfter            int
	TokensReduction        int
	TokensReductionPercent int
	ComponentFiles         []string
	GlobalFiles            []string
	TemplateClasses        []string
	ComponentMatches       []ComponentMatch
}

// Build constructs the final CSS for the microfrontend by:
//  1. Extracting CSS classes directly from the template HTML.
//  2. Reading global CSS files (excluding tokens) from the global directory.
//  3. Dynamically discovering which component CSS files are needed based on the classes in the template.
//  4. Combining global CSS with the matched component CSS.
//  5. Tree-shaking tokens to only keep custom properties referenced by the combined CSS.
//  6. Rewriting :root to :host for Shadow DOM scoping.
func Build(dsCssDir string, templateHTML string) (*BuildResult, error) {
	globalDir := filepath.Join(dsCssDir, "global")
	componentDir := filepath.Join(dsCssDir, "component")

	// 0. Extract classes from the template HTML for reporting.
	templateClasses := extractClassesFromHTML(templateHTML)
	sortedTemplateClasses := make([]string, 0, len(templateClasses))
	for cls := range templateClasses {
		sortedTemplateClasses = append(sortedTemplateClasses, cls)
	}
	sort.Strings(sortedTemplateClasses)

	// 1. Read global CSS files (*.min.css, excluding tokens.*).
	globalFiles, globalCss, err := readGlobalCss(globalDir)

	if err != nil {
		return nil, fmt.Errorf("reading global CSS: %w", err)
	}

	// 2. Read tokens.css (non-minified, for tree-shaking).
	tokens, err := os.ReadFile(filepath.Join(globalDir, "tokens.css"))

	if err != nil {
		return nil, fmt.Errorf("reading tokens.css: %w", err)
	}

	// 3. Discover component CSS files needed based on classes in the template HTML.
	componentFiles, err := ReadComponentCssFiles(componentDir)

	if err != nil {
		return nil, fmt.Errorf("reading component CSS files: %w", err)
	}

	matchedFiles, componentCss, componentMatches := FindComponentCss(componentFiles, templateHTML)

	// 4. Combine global + component CSS, then tree-shake tokens.
	allConsumerCss := globalCss + "\n" + componentCss

	tokensCss := string(tokens)
	treeShakenTokens := TreeShakeTokens(tokensCss, allConsumerCss)

	tokensBefore := len(tokensCss)
	tokensAfter := len(treeShakenTokens)
	tokensReduction := tokensBefore - tokensAfter
	tokensReductionPercent := 0

	if tokensBefore > 0 {
		tokensReductionPercent = int(float64(tokensReduction) / float64(tokensBefore) * 100)
	}

	// 5. Assemble final CSS: tree-shaken tokens + global + component.
	rawCss := treeShakenTokens + "\n" + allConsumerCss

	// Rewrite :root to :host for Shadow DOM scoping, then deduplicate :host, :host pairs.
	css := strings.ReplaceAll(rawCss, ":root", ":host")

	for strings.Contains(css, ":host, :host") {
		css = strings.ReplaceAll(css, ":host, :host", ":host")
	}

	return &BuildResult{
		CSS:                    css,
		TokensBefore:           tokensBefore,
		TokensAfter:            tokensAfter,
		TokensReduction:        tokensReduction,
		TokensReductionPercent: tokensReductionPercent,
		ComponentFiles:         matchedFiles,
		GlobalFiles:            globalFiles,
		TemplateClasses:        sortedTemplateClasses,
		ComponentMatches:       componentMatches,
	}, nil
}

// readGlobalCss reads all *.min.css files from the global directory, excluding tokens.* files.
// Returns the sorted file names and the concatenated CSS content.
func readGlobalCss(dir string) ([]string, string, error) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, "", err
	}

	var names []string

	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(name, ".min.css") && !strings.HasPrefix(name, "tokens.") {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	parts := make([]string, 0, len(names))

	for _, name := range names {
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, "", err
		}
		parts = append(parts, string(content))
	}

	return names, strings.Join(parts, "\n"), nil
}
