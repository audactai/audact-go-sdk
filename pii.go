package audact

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// PII identifier types accepted by NormalizeIdentifier and ClientHash.
const (
	// PIIIdentifierName is a person's name.
	PIIIdentifierName = "name"
	// PIIIdentifierEmail is an e-mail address.
	PIIIdentifierEmail = "email"
	// PIIIdentifierPhone is an E.164 phone number.
	PIIIdentifierPhone = "phone"
)

var e164Pattern = regexp.MustCompile(`^\+[1-9][0-9]{6,14}$`)

// NormalizeIdentifier canonicalises a subject identifier. typ is one of
// "name", "email" or "phone". Phones must already be E.164 once separators are
// stripped; no country code is ever guessed. Names are not diacritic-stripped.
func NormalizeIdentifier(typ, value string) (string, error) {
	switch typ {
	case PIIIdentifierPhone:
		stripped := strings.Map(func(r rune) rune {
			switch r {
			case ' ', '-', '.', '(', ')':
				return -1
			}
			return r
		}, value)
		if !e164Pattern.MatchString(stripped) {
			return "", errors.New("normalize phone: not E.164")
		}
		return stripped, nil
	case PIIIdentifierEmail:
		out := nfc(strings.ToLower(strings.TrimSpace(value)))
		if out == "" {
			return "", errors.New("normalize email: empty value")
		}
		return out, nil
	case PIIIdentifierName:
		out := strings.ToLower(nfc(strings.TrimSpace(value)))
		out = strings.Join(strings.FieldsFunc(out, unicode.IsSpace), " ")
		if out == "" {
			return "", errors.New("normalize name: empty value")
		}
		return out, nil
	default:
		return "", fmt.Errorf("normalize identifier: unsupported type %q", typ)
	}
}

// ClientHash returns the lowercase hex SHA-256 of NormalizeIdentifier(typ,
// value). Send this hash, never the plaintext identifier.
func ClientHash(typ, value string) (string, error) {
	n, err := NormalizeIdentifier(typ, value)
	if err != nil {
		return "", fmt.Errorf("client hash: %w", err)
	}
	sum := sha256.Sum256([]byte(n))
	return hex.EncodeToString(sum[:]), nil
}
