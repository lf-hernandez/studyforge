package utils

import (
	"regexp"
	"strings"
)

// CleanPDFText removes common PDF extraction artifacts
func CleanPDFText(text string) string {
	// Remove spurious 'i' at word boundaries (common artifact from ledongthuc/pdf)

	// Pattern 1: lowercase + 'i' + uppercase (e.g., "nationsiand" -> "nations and")
	re1 := regexp.MustCompile(`([a-z])i([A-Z])`)
	text = re1.ReplaceAllString(text, "$1 $2")

	// Pattern 2: Specific problematic combinations only
	// "shi" followed by consonants that start new words (e.g., "Spanishiforce")
	re2 := regexp.MustCompile(`([s][h])i([bcdfghjklmnpqrstvwxyz])`)
	text = re2.ReplaceAllString(text, "$1 $2")

	// Pattern 2b: Accented letter + s + i + consonant (e.g., "Cortésihoped")
	re2b := regexp.MustCompile(`([éáíóúàèìòùñ][s])i([bcdfghjklmnpqrstvwxyz])`)
	text = re2b.ReplaceAllString(text, "$1 $2")

	// Pattern 2c: Common word endings + i + consonant at word boundary
	// Look for patterns where a complete word is followed by i + new word
	// "nish" is end of "Spanish", "tés" is end of "Cortés"
	re2c := regexp.MustCompile(`(nish|tés|rés)i([bcdfghjklmnpqrstvwxyz])`)
	text = re2c.ReplaceAllString(text, "$1 $2")

	// Pattern 3: 'i' between lowercase and start of sentence after period
	re3 := regexp.MustCompile(`([a-z])\.i([A-Z])`)
	text = re3.ReplaceAllString(text, "$1. $2")

	// Pattern 4: 'i' after punctuation before capital letter
	re4 := regexp.MustCompile(`([,;:])i([A-Z])`)
	text = re4.ReplaceAllString(text, "$1 $2")

	// Pattern 5: Common word endings with 'i' before space or punctuation
	// Be very specific to avoid breaking words like "their"
	re5 := regexp.MustCompile(`\bthei\s`)
	text = re5.ReplaceAllString(text, "the ")

	re6 := regexp.MustCompile(`\bandi\s`)
	text = re6.ReplaceAllString(text, "and ")

	re7 := regexp.MustCompile(`\btoi\s`)
	text = re7.ReplaceAllString(text, "to ")

	re8 := regexp.MustCompile(`\bofi\s`)
	text = re8.ReplaceAllString(text, "of ")

	re9 := regexp.MustCompile(`\bini\s`)
	text = re9.ReplaceAllString(text, "in ")

	re10 := regexp.MustCompile(`\bfori\s`)
	text = re10.ReplaceAllString(text, "for ")

	re11 := regexp.MustCompile(`\bwithi\s`)
	text = re11.ReplaceAllString(text, "with ")

	re12 := regexp.MustCompile(`\bfromi\s`)
	text = re12.ReplaceAllString(text, "from ")

	// Fix missing spaces after commas (e.g., "timid,malleable" -> "timid, malleable")
	re13 := regexp.MustCompile(`([a-z]),([a-z])`)
	text = re13.ReplaceAllString(text, "$1, $2")

	// Remove URLs and web references
	// Pattern: http:// or https:// URLs
	re14 := regexp.MustCompile(`https?://[^\s\)]+`)
	text = re14.ReplaceAllString(text, "")

	// Pattern: (http://...) in parentheses
	re15 := regexp.MustCompile(`\(https?://[^\)]+\)`)
	text = re15.ReplaceAllString(text, "")

	// Pattern: Explore/Visit/Check references with links
	re16 := regexp.MustCompile(`(?:Explore|Visit|Check out|See)\s+[^.]+?\(http[^\)]+\)`)
	text = re16.ReplaceAllString(text, "")

	// Remove citation markers and figure references that don't add content value
	// Example: "FIGURE2.5This" or "Figure 2.5"
	re17 := regexp.MustCompile(`FIGURE\s*\d+\.\d+`)
	text = re17.ReplaceAllString(text, "")

	// Remove multiple spaces
	re18 := regexp.MustCompile(`\s+`)
	text = re18.ReplaceAllString(text, " ")

	// Clean up spaces around punctuation
	text = strings.ReplaceAll(text, " .", ".")
	text = strings.ReplaceAll(text, " ,", ",")
	text = strings.ReplaceAll(text, " ;", ";")
	text = strings.ReplaceAll(text, " :", ":")

	return strings.TrimSpace(text)
}
