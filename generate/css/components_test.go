package main

import (
	"testing"
)

func TestExtractClassesFromHTML_SingleClass(t *testing.T) {
	classes := extractClassesFromHTML(`<div class="foo">content</div>`)

	if _, ok := classes["foo"]; !ok {
		t.Error("expected classes to contain 'foo'")
	}
}

func TestExtractClassesFromHTML_MultipleClasses(t *testing.T) {
	classes := extractClassesFromHTML(`<div class="foo bar baz">content</div>`)

	for _, expected := range []string{"foo", "bar", "baz"} {
		if _, ok := classes[expected]; !ok {
			t.Errorf("expected classes to contain %q", expected)
		}
	}
}

func TestExtractClassesFromHTML_MultipleElements(t *testing.T) {
	html := `<div class="a"><span class="b c">text</span></div>`
	classes := extractClassesFromHTML(html)

	for _, expected := range []string{"a", "b", "c"} {
		if _, ok := classes[expected]; !ok {
			t.Errorf("expected classes to contain %q", expected)
		}
	}
}

func TestExtractClassesFromHTML_NoClasses(t *testing.T) {
	classes := extractClassesFromHTML(`<div><span>no classes here</span></div>`)

	if len(classes) != 0 {
		t.Errorf("expected no classes, got %d", len(classes))
	}
}

func TestExtractClassesFromHTML_EmptyString(t *testing.T) {
	classes := extractClassesFromHTML("")

	if len(classes) != 0 {
		t.Errorf("expected no classes, got %d", len(classes))
	}
}

func TestExtractClassesFromHTML_EmptyClassAttribute(t *testing.T) {
	classes := extractClassesFromHTML(`<div class="">content</div>`)

	if len(classes) != 0 {
		t.Errorf("expected no classes, got %d", len(classes))
	}
}

func TestExtractClassesFromCSS_SimpleSelectors(t *testing.T) {
	classes := extractClassesFromCSS(`.foo { color: red; } .bar { color: blue; }`)

	for _, expected := range []string{"foo", "bar"} {
		if _, ok := classes[expected]; !ok {
			t.Errorf("expected classes to contain %q", expected)
		}
	}
}

func TestExtractClassesFromCSS_BEMSelectors(t *testing.T) {
	css := `.aksel-link-card { } .aksel-link-card--small { } .aksel-link-card__icon { }`
	classes := extractClassesFromCSS(css)

	for _, expected := range []string{"aksel-link-card", "aksel-link-card--small", "aksel-link-card__icon"} {
		if _, ok := classes[expected]; !ok {
			t.Errorf("expected classes to contain %q", expected)
		}
	}
}

func TestExtractClassesFromCSS_EmptyString(t *testing.T) {
	classes := extractClassesFromCSS("")

	if len(classes) != 0 {
		t.Errorf("expected no classes, got %d", len(classes))
	}
}

func TestFindComponentCss_MatchesByClassOverlap(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "linkcard.min.css",
			Content: ".aksel-link-card { display: block; }",
			Classes: map[string]struct{}{"aksel-link-card": {}},
		},
		{
			Name:    "button.min.css",
			Content: ".aksel-button { display: inline; }",
			Classes: map[string]struct{}{"aksel-button": {}},
		},
	}

	html := `<div class="aksel-link-card">content</div>`
	files, css, _ := FindComponentCss(componentFiles, html)

	if len(files) != 1 {
		t.Fatalf("expected 1 matched file, got %d", len(files))
	}
	if files[0] != "linkcard.min.css" {
		t.Errorf("expected linkcard.min.css, got %s", files[0])
	}
	if css == "" {
		t.Error("expected non-empty CSS")
	}
	if !containsStr(css, ".aksel-link-card") {
		t.Error("expected CSS to contain .aksel-link-card")
	}
}

func TestFindComponentCss_NoMatchReturnsEmpty(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "button.min.css",
			Content: ".aksel-button { display: inline; }",
			Classes: map[string]struct{}{"aksel-button": {}},
		},
	}

	html := `<div class="my-custom-class">content</div>`
	files, css, _ := FindComponentCss(componentFiles, html)

	if len(files) != 0 {
		t.Errorf("expected 0 matched files, got %d", len(files))
	}
	if css != "" {
		t.Errorf("expected empty CSS, got %q", css)
	}
}

func TestFindComponentCss_MatchesMultipleComponents(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "button.min.css",
			Content: ".aksel-button { display: inline; }",
			Classes: map[string]struct{}{"aksel-button": {}},
		},
		{
			Name:    "linkcard.min.css",
			Content: ".aksel-link-card { display: block; }",
			Classes: map[string]struct{}{"aksel-link-card": {}, "aksel-link-card--small": {}},
		},
		{
			Name:    "typography.min.css",
			Content: ".aksel-heading { font-weight: bold; }",
			Classes: map[string]struct{}{"aksel-heading": {}, "aksel-heading--small": {}},
		},
	}

	html := `<div class="aksel-link-card"><h2 class="aksel-heading">Title</h2></div>`
	files, css, _ := FindComponentCss(componentFiles, html)

	if len(files) != 2 {
		t.Fatalf("expected 2 matched files, got %d: %v", len(files), files)
	}

	hasLinkcard := false
	hasTypography := false
	for _, f := range files {
		if f == "linkcard.min.css" {
			hasLinkcard = true
		}
		if f == "typography.min.css" {
			hasTypography = true
		}
	}

	if !hasLinkcard {
		t.Error("expected linkcard.min.css to be matched")
	}
	if !hasTypography {
		t.Error("expected typography.min.css to be matched")
	}
	if !containsStr(css, ".aksel-link-card") {
		t.Error("expected CSS to contain .aksel-link-card")
	}
	if !containsStr(css, ".aksel-heading") {
		t.Error("expected CSS to contain .aksel-heading")
	}
}

func TestFindComponentCss_MultipleHTMLStrings(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "button.min.css",
			Content: ".aksel-button { display: inline; }",
			Classes: map[string]struct{}{"aksel-button": {}},
		},
		{
			Name:    "linkcard.min.css",
			Content: ".aksel-link-card { display: block; }",
			Classes: map[string]struct{}{"aksel-link-card": {}},
		},
	}

	html1 := `<div class="aksel-link-card">card</div>`
	html2 := `<button class="aksel-button">click</button>`
	files, _, _ := FindComponentCss(componentFiles, html1, html2)

	if len(files) != 2 {
		t.Fatalf("expected 2 matched files, got %d: %v", len(files), files)
	}
}

func TestFindComponentCss_DeduplicatesFiles(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "linkcard.min.css",
			Content: ".aksel-link-card { display: block; }",
			Classes: map[string]struct{}{"aksel-link-card": {}},
		},
	}

	html1 := `<div class="aksel-link-card">card 1</div>`
	html2 := `<div class="aksel-link-card">card 2</div>`
	files, _, _ := FindComponentCss(componentFiles, html1, html2)

	count := 0
	for _, f := range files {
		if f == "linkcard.min.css" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected linkcard.min.css to appear exactly once, got %d", count)
	}
}

func TestFindComponentCss_EmptyHTMLReturnsEmpty(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "button.min.css",
			Content: ".aksel-button { display: inline; }",
			Classes: map[string]struct{}{"aksel-button": {}},
		},
	}

	files, css, _ := FindComponentCss(componentFiles, "")

	if len(files) != 0 {
		t.Errorf("expected 0 matched files, got %d", len(files))
	}
	if css != "" {
		t.Errorf("expected empty CSS, got %q", css)
	}
}

func TestFindComponentCss_BEMElementClassMatchesComponentFile(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "linkcard.min.css",
			Content: ".aksel-link-card { } .aksel-link-card__icon { }",
			Classes: map[string]struct{}{"aksel-link-card": {}, "aksel-link-card__icon": {}},
		},
	}

	html := `<div class="aksel-link-card__icon">icon</div>`
	files, _, _ := FindComponentCss(componentFiles, html)

	if len(files) != 1 || files[0] != "linkcard.min.css" {
		t.Errorf("expected [linkcard.min.css], got %v", files)
	}
}

func TestFindComponentCss_MatchedFilesAreSorted(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "linkanchor.min.css",
			Content: ".aksel-link-anchor { }",
			Classes: map[string]struct{}{"aksel-link-anchor": {}},
		},
		{
			Name:    "linkcard.min.css",
			Content: ".aksel-link-card { }",
			Classes: map[string]struct{}{"aksel-link-card": {}},
		},
		{
			Name:    "typography.min.css",
			Content: ".aksel-heading { }",
			Classes: map[string]struct{}{"aksel-heading": {}},
		},
	}

	html := `<div class="aksel-heading aksel-link-card aksel-link-anchor">content</div>`
	files, _, _ := FindComponentCss(componentFiles, html)

	for i := 1; i < len(files); i++ {
		if files[i] < files[i-1] {
			t.Errorf("expected files to be sorted, but got %v", files)
			break
		}
	}
}

func TestFindComponentCss_ActualMicrofrontendMarkup(t *testing.T) {
	componentFiles := []ComponentCssFile{
		{
			Name:    "button.min.css",
			Content: ".aksel-button { }",
			Classes: map[string]struct{}{"aksel-button": {}},
		},
		{
			Name:    "linkanchor.min.css",
			Content: ".aksel-link-anchor { } .aksel-link-anchor__overlay { } .aksel-link-anchor__arrow { }",
			Classes: map[string]struct{}{"aksel-link-anchor": {}, "aksel-link-anchor__overlay": {}, "aksel-link-anchor__arrow": {}},
		},
		{
			Name:    "linkcard.min.css",
			Content: ".aksel-link-card { } .aksel-link-card--medium { } .aksel-link-card__icon { } .aksel-link-card__title { } .aksel-link-card__description { } .aksel-link-card__arrow { }",
			Classes: map[string]struct{}{
				"aksel-link-card": {}, "aksel-link-card--medium": {},
				"aksel-link-card__icon": {}, "aksel-link-card__title": {},
				"aksel-link-card__description": {}, "aksel-link-card__arrow": {},
			},
		},
		{
			Name:    "typography.min.css",
			Content: ".aksel-heading { } .aksel-heading--small { } .aksel-body-long { } .aksel-body-long--medium { }",
			Classes: map[string]struct{}{
				"aksel-heading": {}, "aksel-heading--small": {},
				"aksel-body-long": {}, "aksel-body-long--medium": {},
			},
		},
		{
			Name:    "table.min.css",
			Content: ".aksel-table { }",
			Classes: map[string]struct{}{"aksel-table": {}},
		},
	}

	html := `<div data-color="neutral" data-align-arrow="baseline" class="aksel-link-anchor__overlay aksel-link-card aksel-link-card--medium aksel-body-long aksel-body-long--medium">` +
		`<div aria-hidden="true" class="aksel-link-card__icon"><svg></svg></div>` +
		`<h2 class="aksel-link-card__title aksel-heading aksel-heading--small">` +
		`<a href="#" class="aksel-link-anchor">Title</a></h2>` +
		`<div class="aksel-link-card__description">Description</div>` +
		`<svg class="aksel-link-anchor__arrow aksel-link-card__arrow"></svg>` +
		`</div>`

	files, _, _ := FindComponentCss(componentFiles, html)

	if len(files) != 3 {
		t.Fatalf("expected exactly 3 matched files, got %d: %v", len(files), files)
	}

	expected := map[string]bool{
		"linkanchor.min.css": false,
		"linkcard.min.css":   false,
		"typography.min.css": false,
	}

	for _, f := range files {
		if _, ok := expected[f]; ok {
			expected[f] = true
		} else {
			t.Errorf("unexpected file matched: %s", f)
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("expected %s to be matched but it wasn't", name)
		}
	}
}
