package routego

import "regexp"

var (
	singleBracePattern = regexp.MustCompile(`{.*?}`)
	doubleBracePattern = regexp.MustCompile(`{{.*?}}`)
)