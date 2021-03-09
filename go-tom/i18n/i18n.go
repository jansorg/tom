package i18n

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/de_DE"
	"github.com/go-playground/locales/en_US"
	"golang.org/x/text/language"
)

func FindPreferredLanguages() language.Tag {
	supported := []language.Tag{
		language.English,
		language.German,
	}

	matcher := language.NewMatcher(supported)
	tag, _, _ := matcher.Match(userPreference()...)
	return tag
}

func FindLocale(lang language.Tag, isoDates bool) locales.Translator {
	translator := func() locales.Translator {
		s := lang.String()
		switch s {
		case "de_DE":
			return de_DE.New()
		}

		b, _ := lang.Base()
		switch b.String() {
		case "de":
			return de_DE.New()
		}

		// fallback
		return en_US.New()
	}
	if isoDates {
		return &isoDelegate{locateDelegate{delegate: translator()}}
	}
	return translator()
}

func userPreference() []language.Tag {
	if langEnv, err := GetLocale(); err == nil {
		if lang, err := language.Parse(langEnv); err == nil {
			return []language.Tag{lang}
		}
	}

	return []language.Tag{
		language.English,
	}
}

// from https://stackoverflow.com/questions/51829386/golang-get-system-language
func GetLocale() (string, error) {
	// Check the LANG environment variable, common on UNIX.
	// XXX: we can easily override as a nice feature/bug.
	envLang, ok := os.LookupEnv("LANG")
	if ok {
		return strings.Split(envLang, ".")[0], nil
	}

	// Exec powershell Get-Culture on Windows.
	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	output, err := cmd.Output()
	if err == nil {
		return strings.Trim(string(output), "\r\n"), nil
	}

	return "", fmt.Errorf("cannot determine locale")
}
