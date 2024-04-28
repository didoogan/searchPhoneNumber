package searchphonenumber

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const SimpleNumberPatternStr = `^\d{10}`
const NumberWithParenthesesPatternStr = `^\(\d{3}\)[\s]?\d{3}-\d{4}$`
const DelimitedNumberPatternStr = `^\d{3}[\s\.-]{1}\d{3}[\s\.-]{1}\d{4}$`

type SearchEngine struct {
	FilePath string
	Tokens   []string
	Result   []string
}

func (e *SearchEngine) ShowTokens() {
	for _, w := range e.Tokens {
		fmt.Println(w)
	}
}

func (e *SearchEngine) clearTokens() {
	e.Tokens = make([]string, 0)
}

func (e *SearchEngine) clearResult() {
	e.Result = make([]string, 0)
}

func (e *SearchEngine) extractTokensFromFile(callBack func(e *SearchEngine, fileLine string)) error {
	e.clearTokens()

	f, err := os.Open(e.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scn := bufio.NewScanner(f)

	for scn.Scan() {
		fileLine := scn.Text()

		callBack(e, fileLine)
	}

	if err := scn.Err(); err != nil {
		return err
	}

	return nil
}

func (e *SearchEngine) ExtractLinesFromFile() error {
	cb := func(e *SearchEngine, fileLine string) {
		e.Tokens = append(e.Tokens, fileLine)
	}

	return e.extractTokensFromFile(cb)
}

func (e *SearchEngine) ExtractWordsFromFile() error {
	redundantSuffixes := []string{".", ",", ";", ":", "!", "?"}

	cb := func(e *SearchEngine, fileLine string) {

		for _, w := range strings.Fields(fileLine) {
			for _, s := range redundantSuffixes {
				if strings.HasSuffix(w, s) {
					w = strings.TrimSuffix(w, s)
				}
			}
			e.Tokens = append(e.Tokens, w)
		}
	}

	return e.extractTokensFromFile(cb)
}

func (e *SearchEngine) SearchByPattern(re *regexp.Regexp) {
	e.clearResult()

	for _, w := range e.Tokens {
		isValid := re.MatchString(w)

		if isValid {
			e.Result = append(e.Result, w)
		}
	}
}
