package finder

import "strings"

//Matcher represents the interface which checks for matching search
type Matcher interface {
	Match(line string, search string) []string
}

//FullMatcher matches the full word with search text and is case sensitive
type FullMatcher struct {
}

//FullMatcherIgnoreCase matches the full word with search text and is case in-sensitive
type FullMatcherIgnoreCase struct {
}

//PartialMatcher matches any word containing the search text and is case sensitive
type PartialMatcher struct {
}

//PartialMatcherIgnoreCase matches any word containing the search text and is case in-sensitive
type PartialMatcherIgnoreCase struct {
}

//Match matches the full word with search text and is case sensitive
func (m FullMatcher) Match(line string, search string) []string {

	if len(strings.Fields(search)) > 1 {
		return m.matchFullLine(line, search)
	}

	return m.matchSingleWords(line, search)

}

func (m FullMatcher) matchFullLine(line string, search string) []string {
	matches := make([]string, 0)

	if strings.Index(line, search) != -1 {
		matches = append(matches, search)
	}
	return matches

}

func (m FullMatcher) matchSingleWords(line string, search string) []string {
	matches := make([]string, 0)
	for _, word := range strings.Fields(line) {
		if word == search {
			matches = append(matches, word)
		}
	}
	return matches

}

//Match matches the full word with search text and is case in-sensitive
func (m FullMatcherIgnoreCase) Match(line string, search string) []string {

	return FullMatcher{}.Match(strings.ToLower(line), strings.ToLower(search))
}

//Match matches any word containing the search text and is case sensitive
func (m PartialMatcher) Match(line string, search string) []string {

	matches := make([]string, 0)
	words := strings.Fields(line)
	searches := strings.Fields(search)

	for _, w := range words {
		for _, s := range searches {
			if strings.Contains(w, s) {
				matches = append(matches, w)
			}
		}
	}

	return matches
}

//Match matches any word containing the search text and is case in-sensitive
func (m PartialMatcherIgnoreCase) Match(line string, search string) []string {

	return PartialMatcher{}.Match(strings.ToLower(line), strings.ToLower(search))
}
