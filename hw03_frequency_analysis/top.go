package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var replaceRegExp = regexp.MustCompile(`([.,;]+)|([^a-zа-я]+-[^a-zа-я]+)`)

func Top10(inStr string) []string {
	inStr = strings.ToLower(inStr)
	inStr = replaceRegExp.ReplaceAllString(inStr, " ")
	splitWords := strings.Fields(inStr)

	wordsMap := make(map[string]int)
	for _, word := range splitWords {
		if word == "" || word == "-" {
			continue
		}

		_, keyExists := wordsMap[word]
		if !keyExists {
			wordsMap[word] = 0
		}
		wordsMap[word]++
	}

	type userMap struct {
		Key   string
		Value int
	}

	var wordsSlice []userMap //nolint:prealloc
	for k, v := range wordsMap {
		wordsSlice = append(wordsSlice, userMap{k, v})
	}

	sort.Slice(wordsSlice, func(i, j int) bool {
		if wordsSlice[i].Value == wordsSlice[j].Value {
			return wordsSlice[i].Key < wordsSlice[j].Key
		}
		return wordsSlice[i].Value > wordsSlice[j].Value
	})

	var result []string //nolint:prealloc
	for i := 0; i < len(wordsSlice) && i < 10; i++ {
		result = append(result, wordsSlice[i].Key)
	}

	return result
}
