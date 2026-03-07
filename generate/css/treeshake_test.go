package main

import (
	"regexp"
	"testing"
)

var (
	excessiveNewlinesPattern     = regexp.MustCompile(`\n{3,}`)
	singleTrailingNewlinePattern = regexp.MustCompile(`[^\n]\n$`)
)

func TestTreeShakeTokens_KeepsDirectlyReferencedToken(t *testing.T) {
	tokens := ":root {\n  --ax-color-red: #f00;\n}\n"
	consumer := ".foo { color: var(--ax-color-red); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--ax-color-red: #f00;") {
		t.Errorf("expected result to contain --ax-color-red: #f00; but got:\n%s", result)
	}
}

func TestTreeShakeTokens_RemovesUnreferencedToken(t *testing.T) {
	tokens := ":root {\n  --ax-color-red: #f00;\n  --ax-color-blue: #00f;\n}\n"
	consumer := ".foo { color: var(--ax-color-red); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--ax-color-red") {
		t.Errorf("expected result to contain --ax-color-red but got:\n%s", result)
	}
	if contains(result, "--ax-color-blue") {
		t.Errorf("expected result NOT to contain --ax-color-blue but got:\n%s", result)
	}
}

func TestTreeShakeTokens_KeepsTransitivelyReferencedTokens(t *testing.T) {
	tokens := ":root {\n  --ax-surface: var(--ax-blue-500);\n  --ax-blue-500: #0060c0;\n  --ax-unused: #999;\n}"
	consumer := ".card { background: var(--ax-surface); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--ax-surface") {
		t.Errorf("expected result to contain --ax-surface")
	}
	if !contains(result, "--ax-blue-500") {
		t.Errorf("expected result to contain --ax-blue-500")
	}
	if contains(result, "--ax-unused") {
		t.Errorf("expected result NOT to contain --ax-unused")
	}
}

func TestTreeShakeTokens_KeepsDeeplyTransitiveTokenChains(t *testing.T) {
	tokens := ":root {\n  --a: var(--b);\n  --b: var(--c);\n  --c: #123;\n  --d: #456;\n}"
	consumer := ".x { color: var(--a); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--a") {
		t.Errorf("expected result to contain --a")
	}
	if !contains(result, "--b") {
		t.Errorf("expected result to contain --b")
	}
	if !contains(result, "--c") {
		t.Errorf("expected result to contain --c")
	}
	if contains(result, "--d") {
		t.Errorf("expected result NOT to contain --d")
	}
}

func TestTreeShakeTokens_HandlesMultipleVarRefsInSingleValue(t *testing.T) {
	tokens := ":root {\n  --ax-shadow: 0 2px 4px var(--ax-shadow-color) var(--ax-shadow-alpha);\n  --ax-shadow-color: #000;\n  --ax-shadow-alpha: 0.2;\n  --ax-unused: red;\n}"
	consumer := ".box { box-shadow: var(--ax-shadow); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--ax-shadow-color") {
		t.Errorf("expected result to contain --ax-shadow-color")
	}
	if !contains(result, "--ax-shadow-alpha") {
		t.Errorf("expected result to contain --ax-shadow-alpha")
	}
	if contains(result, "--ax-unused") {
		t.Errorf("expected result NOT to contain --ax-unused")
	}
}

func TestTreeShakeTokens_HandlesMultipleVarRefsInConsumerCSS(t *testing.T) {
	tokens := ":root {\n  --ax-color-red: #f00;\n  --ax-color-blue: #00f;\n  --ax-color-green: #0f0;\n}"
	consumer := ".a { color: var(--ax-color-red); } .b { color: var(--ax-color-blue); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--ax-color-red") {
		t.Errorf("expected result to contain --ax-color-red")
	}
	if !contains(result, "--ax-color-blue") {
		t.Errorf("expected result to contain --ax-color-blue")
	}
	if contains(result, "--ax-color-green") {
		t.Errorf("expected result NOT to contain --ax-color-green")
	}
}

func TestTreeShakeTokens_RemovesEmptyRuleBlocksAfterStripping(t *testing.T) {
	tokens := ":root {\n  --ax-used: #f00;\n}\n.dark {\n  --ax-unused: #0f0;\n}"
	consumer := ".x { color: var(--ax-used); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, ":root {") {
		t.Errorf("expected result to contain :root {")
	}
	if contains(result, ".dark") {
		t.Errorf("expected result NOT to contain .dark")
	}
}

func TestTreeShakeTokens_NoTokensReferenced(t *testing.T) {
	tokens := ":root {\n  --ax-color-red: #f00;\n  --ax-color-blue: #00f;\n}"
	consumer := ".foo { color: red; }"

	result := TreeShakeTokens(tokens, consumer)

	if contains(result, "--ax-color-red") {
		t.Errorf("expected result NOT to contain --ax-color-red")
	}
	if contains(result, "--ax-color-blue") {
		t.Errorf("expected result NOT to contain --ax-color-blue")
	}
}

func TestTreeShakeTokens_PreservesNonDeclarationLinesLikeComments(t *testing.T) {
	tokens := "/* Token definitions */\n:root {\n  --ax-used: #f00;\n  --ax-unused: #00f;\n}"
	consumer := ".x { color: var(--ax-used); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "/* Token definitions */") {
		t.Errorf("expected result to contain comment")
	}
	if !contains(result, "--ax-used") {
		t.Errorf("expected result to contain --ax-used")
	}
	if contains(result, "--ax-unused") {
		t.Errorf("expected result NOT to contain --ax-unused")
	}
}

func TestTreeShakeTokens_DoesNotLeaveExcessiveBlankLines(t *testing.T) {
	tokens := ":root {\n  --a: #111;\n  --b: #222;\n  --c: #333;\n  --d: #444;\n  --e: #555;\n}"
	consumer := ".x { color: var(--a); background: var(--e); }"

	result := TreeShakeTokens(tokens, consumer)

	if excessiveNewlinesPattern.MatchString(result) {
		t.Errorf("result contains excessive blank lines:\n%s", result)
	}
}

func TestTreeShakeTokens_ResultEndsWithExactlyOneNewline(t *testing.T) {
	tokens := ":root {\n  --a: #f00;\n}\n"
	consumer := ".x { color: var(--a); }"

	result := TreeShakeTokens(tokens, consumer)

	if !singleTrailingNewlinePattern.MatchString(result) {
		t.Errorf("expected result to end with exactly one newline, got: %q", result[len(result)-5:])
	}
}

func TestTreeShakeTokens_HandlesVarRefsWithFallbackValues(t *testing.T) {
	tokens := ":root {\n  --ax-primary: #f00;\n  --ax-unused: #0f0;\n}"
	consumer := ".x { color: var(--ax-primary, red); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--ax-primary") {
		t.Errorf("expected result to contain --ax-primary")
	}
	if contains(result, "--ax-unused") {
		t.Errorf("expected result NOT to contain --ax-unused")
	}
}

func TestTreeShakeTokens_HandlesPropertyDefinedInMultipleBlocks(t *testing.T) {
	tokens := `:root {
  --ax-border-neutral-subtleA: var(--ax-neutral-400A);
  --ax-neutral-400A: rgba(0, 22, 48, .19);
  --ax-border-meta-lime-subtleA: var(--ax-meta-lime-400A);
  --ax-meta-lime-400A: rgba(172, 191, 0, .576);
}
[data-color="neutral"] {
  --ax-border-subtleA: var(--ax-border-neutral-subtleA);
}
[data-color="meta-lime"] {
  --ax-border-subtleA: var(--ax-border-meta-lime-subtleA);
}`
	consumer := ".card { border: 1px solid var(--ax-border-subtleA); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--ax-border-subtleA") {
		t.Error("expected result to contain --ax-border-subtleA")
	}
	if !contains(result, "--ax-border-neutral-subtleA") {
		t.Error("expected result to contain --ax-border-neutral-subtleA")
	}
	if !contains(result, "--ax-neutral-400A") {
		t.Error("expected result to contain --ax-neutral-400A")
	}
	if !contains(result, "--ax-border-meta-lime-subtleA") {
		t.Error("expected result to contain --ax-border-meta-lime-subtleA")
	}
	if !contains(result, "--ax-meta-lime-400A") {
		t.Error("expected result to contain --ax-meta-lime-400A")
	}
}

func TestTreeShakeTokens_HandlesCircularReferencesWithoutInfiniteLoop(t *testing.T) {
	tokens := ":root {\n  --a: var(--b);\n  --b: var(--a);\n}"
	consumer := ".x { color: var(--a); }"

	result := TreeShakeTokens(tokens, consumer)

	if !contains(result, "--a") {
		t.Errorf("expected result to contain --a")
	}
	if !contains(result, "--b") {
		t.Errorf("expected result to contain --b")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
