package main

import (
	"regexp"
	"strings"
)

var (
	declLinePattern       = regexp.MustCompile(`^\s*(--[\w-]+)\s*:`)
	emptyBlockPattern     = regexp.MustCompile(`(?m)^[^\n{}]*\{\s*\}\s*$`)
	excessNewlinesPattern = regexp.MustCompile(`\n{3,}`)
	varRefPattern         = regexp.MustCompile(`var\((--[\w-]+)`)
	declarationPattern    = regexp.MustCompile(`(?m)^\s*(--[\w-]+)\s*:\s*(.+?)\s*;`)
)

// TreeShakeTokens removes unused CSS custom property declarations from tokensCss,
// keeping only those that are directly or transitively referenced by consumerCss.
func TreeShakeTokens(tokensCss, consumerCss string) string {
	needed := extractVarRefs(consumerCss)

	// Build a dependency map: property -> set of var references in its value.
	depsOf := make(map[string]map[string]struct{})
	for _, decl := range parseDeclarations(tokensCss) {
		refs := extractVarRefs(decl.value)
		if len(refs) > 0 {
			existing, ok := depsOf[decl.property]
			if !ok {
				existing = make(map[string]struct{})
				depsOf[decl.property] = existing
			}
			for ref := range refs {
				existing[ref] = struct{}{}
			}
		}
	}

	// Walk the dependency graph to find all transitively needed properties.
	queue := make([]string, 0, len(needed))
	for prop := range needed {
		queue = append(queue, prop)
	}

	for len(queue) > 0 {
		prop := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		for dep := range depsOf[prop] {
			if _, ok := needed[dep]; !ok {
				needed[dep] = struct{}{}
				queue = append(queue, dep)
			}
		}
	}

	// Filter token lines: keep non-declaration lines and declarations for needed properties.
	lines := strings.Split(tokensCss, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		matches := declLinePattern.FindStringSubmatch(line)
		if matches == nil {
			// Not a custom property declaration line — keep it.
			kept = append(kept, line)
		} else {
			property := matches[1]
			if _, ok := needed[property]; ok {
				kept = append(kept, line)
			}
		}
	}

	result := strings.Join(kept, "\n")

	// Iteratively remove empty blocks and collapse excessive newlines.
	for {
		prev := result
		result = emptyBlockPattern.ReplaceAllString(result, "")
		result = excessNewlinesPattern.ReplaceAllString(result, "\n\n")
		if result == prev {
			break
		}
	}

	return strings.TrimSpace(result) + "\n"
}

// extractVarRefs finds all CSS custom property names referenced via var() in the given CSS string.
func extractVarRefs(css string) map[string]struct{} {
	refs := make(map[string]struct{})
	for _, match := range varRefPattern.FindAllStringSubmatch(css, -1) {
		if len(match) > 1 {
			refs[match[1]] = struct{}{}
		}
	}
	return refs
}

type declaration struct {
	property string
	value    string
}

// parseDeclarations extracts all CSS custom property declarations from the given CSS string.
func parseDeclarations(css string) []declaration {
	matches := declarationPattern.FindAllStringSubmatch(css, -1)
	decls := make([]declaration, 0, len(matches))
	for _, m := range matches {
		if len(m) > 2 {
			decls = append(decls, declaration{property: m[1], value: m[2]})
		}
	}
	return decls
}
