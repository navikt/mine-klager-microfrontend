#!/bin/sh

# Build script to generate HTML files from template and CSS

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

TEMPLATE_FILE="$PROJECT_ROOT/src/template.hbs"
CSS_FILE="$PROJECT_ROOT/src/style.css"
OUTPUT_DIR="$PROJECT_ROOT/html"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Function to generate HTML for a specific language
generate_html() {
    lang="$1"
    title="$2"
    heading="$3"
    description="$4"
    url="$5"

    echo "Generating $lang.html..."

    awk -v title="$title" \
        -v heading="$heading" \
        -v description="$description" \
        -v url="$url" \
        -v css_file="$CSS_FILE" '
    BEGIN {
        css = "\n"
        while ((getline line < css_file) > 0) {
            css = css "    " line "\n"
        }
        close(css_file)
        # Escape & characters to prevent gsub from treating them as backreferences
        gsub(/&/, "\\\\&", css)
    }
    {
        gsub(/\{\{CSS\}\}/, css)
        gsub(/\{\{TITLE\}\}/, title)
        gsub(/\{\{HEADING\}\}/, heading)
        gsub(/\{\{DESCRIPTION\}\}/, description)
        gsub(/\{\{URL\}\}/, url)
        print
    }
    ' "$TEMPLATE_FILE" > "$OUTPUT_DIR/$lang.html"

    echo "Generated $OUTPUT_DIR/$lang.html"
}

# Norwegian Bokmål (nb)
generate_html "nb" \
    "Mine saker hos Klageinstans" \
    "Mine saker hos Klageinstans" \
    "Her kan du se status på dine saker hos Klageinstans." \
    "https://mine-klager.nav.no"

# Norwegian Nynorsk (nn)
generate_html "nn" \
    "Mine saker hjå Klageinstans" \
    "Mine saker hjå Klageinstans" \
    "Her kan du sjå status på dine saker hjå Klageinstans." \
    "https://mine-klager.nav.no/nn"

# English (en)
generate_html "en" \
    "My cases with Nav Complaints Unit" \
    "My cases with Nav Complaints Unit" \
    "Here you can see the status of your cases with Nav Complaints Unit (Klageinstans)." \
    "https://mine-klager.nav.no/en"

echo "HTML files generated!"
