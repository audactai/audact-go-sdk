package audact

import "unicode/utf8"

// Hangul syllable constants (Unicode section 3.12).
const (
	hangulSBase  = 0xAC00
	hangulLBase  = 0x1100
	hangulVBase  = 0x1161
	hangulTBase  = 0x11A7
	hangulLCount = 19
	hangulVCount = 21
	hangulTCount = 28
	hangulNCount = hangulVCount * hangulTCount
	hangulSCount = hangulLCount * hangulNCount
)

// nfc returns s in Unicode Normalization Form C using only the standard
// library (the SDK takes no dependency on golang.org/x/text). Invalid UTF-8
// bytes decode to U+FFFD, as with any range over a string.
func nfc(s string) string {
	ascii := true
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			ascii = false
			break
		}
	}
	if ascii {
		return s
	}
	var buf []rune
	for _, r := range s {
		buf = appendDecomposed(buf, r)
	}
	// Canonical ordering: stable sort of each run of non-starters by ccc.
	for i := 1; i < len(buf); i++ {
		c := nfcCCC[buf[i]]
		if c == 0 {
			continue
		}
		for j := i; j > 0; j-- {
			p := nfcCCC[buf[j-1]]
			if p <= c {
				break
			}
			buf[j-1], buf[j] = buf[j], buf[j-1]
		}
	}
	return string(compose(buf))
}

func appendDecomposed(buf []rune, r rune) []rune {
	if si := r - hangulSBase; si >= 0 && si < hangulSCount {
		buf = append(buf, hangulLBase+si/hangulNCount, hangulVBase+(si%hangulNCount)/hangulTCount)
		if t := si % hangulTCount; t != 0 {
			buf = append(buf, hangulTBase+t)
		}
		return buf
	}
	if d, ok := nfcDecomp[r]; ok {
		for _, x := range d {
			buf = appendDecomposed(buf, x)
		}
		return buf
	}
	return append(buf, r)
}

func composePair(a, b rune) (rune, bool) {
	if l := a - hangulLBase; l >= 0 && l < hangulLCount {
		if v := b - hangulVBase; v >= 0 && v < hangulVCount {
			return hangulSBase + (l*hangulVCount+v)*hangulTCount, true
		}
	}
	if si := a - hangulSBase; si >= 0 && si < hangulSCount && si%hangulTCount == 0 {
		if t := b - hangulTBase; t > 0 && t < hangulTCount {
			return a + t, true
		}
	}
	c, ok := nfcComp[[2]rune{a, b}]
	return c, ok
}

func compose(in []rune) []rune {
	if len(in) == 0 {
		return in
	}
	out := make([]rune, 0, len(in))
	starter := -1
	lastCCC := -1
	for _, r := range in {
		c := int(nfcCCC[r])
		if starter >= 0 && (lastCCC < c || lastCCC == 0) {
			if comp, ok := composePair(out[starter], r); ok {
				out[starter] = comp
				continue
			}
		}
		if c == 0 {
			starter = len(out)
		}
		lastCCC = c
		out = append(out, r)
	}
	return out
}
