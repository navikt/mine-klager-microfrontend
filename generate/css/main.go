package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	dsCssDir := "node_modules/@navikt/ds-css/dist"
	outFile := "templates/style.css"

	templateHTML, err := os.ReadFile("templates/template.html")
	if err != nil {
		log.Fatalf("Failed to read template: %v", err)
	}

	result, err := Build(dsCssDir, string(templateHTML))
	if err != nil {
		log.Fatalf("Failed to build CSS: %v", err)
	}

	fmt.Printf("Template classes: %s\n", join(result.TemplateClasses))
	fmt.Printf("Included global CSS files: %s\n", join(result.GlobalFiles))
	fmt.Printf("Matched component CSS files: %s\n", join(result.ComponentFiles))
	for _, m := range result.ComponentMatches {
		fmt.Printf("  %s ← %s\n", m.File, join(m.MatchedClasses))
	}
	fmt.Printf("Tree-shaken tokens.css: %s → %s (%s / %d%% reduction)\n",
		formatBytes(result.TokensBefore),
		formatBytes(result.TokensAfter),
		formatBytes(result.TokensReduction),
		result.TokensReductionPercent,
	)

	if err := os.WriteFile(outFile, []byte(result.CSS), 0o644); err != nil {
		log.Fatalf("Failed to write %s: %v", outFile, err)
	}

	fmt.Printf("Written %s (%s)\n", outFile, formatBytes(len(result.CSS)))

	summaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
	if summaryFile != "" {
		var matchRows strings.Builder
		for _, m := range result.ComponentMatches {
			fmt.Fprintf(&matchRows, "| `%s` | %s |\n", m.File, formatClassList(m.MatchedClasses))
		}

		summary := fmt.Sprintf(`### 📦 Microfrontend CSS Build

- **HTML:** %s
- **CSS:** %s

🌳 **Token tree-shaking:** %s → %s (%s / %d%% reduction)

**Template classes (%d):** %s

| File | Matched classes |
| --- | --- |
%s
`,
			formatBytes(len(templateHTML)),
			formatBytes(len(result.CSS)),
			formatBytes(result.TokensBefore),
			formatBytes(result.TokensAfter),
			formatBytes(result.TokensReduction),
			result.TokensReductionPercent,
			len(result.TemplateClasses),
			formatClassList(result.TemplateClasses),
			matchRows.String(),
		)

		f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			log.Printf("Warning: failed to open GITHUB_STEP_SUMMARY: %v", err)
		} else {
			defer f.Close()
			if _, err := f.WriteString(summary); err != nil {
				log.Printf("Warning: failed to write GITHUB_STEP_SUMMARY: %v", err)
			}
		}
	}
}

func join(s []string) string {
	if len(s) == 0 {
		return "(none)"
	}

	var result strings.Builder
	result.WriteString(s[0])

	for _, v := range s[1:] {
		result.WriteString(", " + v)
	}

	return result.String()
}

func formatClassList(classes []string) string {
	if len(classes) == 0 {
		return "(none)"
	}

	parts := make([]string, len(classes))
	for i, cls := range classes {
		parts[i] = "`" + cls + "`"
	}

	return strings.Join(parts, " ")
}
