package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(inStr string) []string {
	replaceRegExp := regexp.MustCompile(`([.,;]+)|([^a-zA-Zа-яА-Я]+-[^a-zA-Zа-яА-Я]+)`)
	inStr = replaceRegExp.ReplaceAllString(inStr, " ")

	inStr = strings.ToLower(inStr)

	splitRegExp := regexp.MustCompile(`\s+`)
	splitWords := splitRegExp.Split(inStr, -1)

	wordsMap := make(map[string]int)
	for _, word := range splitWords {
		if word == "" || word == "-" {
			continue
		}

		_, keyExists := wordsMap[word]
		if !keyExists {
			wordsMap[word] = 1
		} else {
			wordsMap[word] += 1
		}
	}

	type userMap struct {
		Key   string
		Value int
	}

	wordsSlice := make([]userMap, 0)
	for k, v := range wordsMap {
		wordsSlice = append(wordsSlice, userMap{k, v})
	}

	sort.Slice(wordsSlice, func(i, j int) bool {
		if wordsSlice[i].Value == wordsSlice[j].Value {
			return wordsSlice[i].Key < wordsSlice[j].Key
		}
		return wordsSlice[i].Value > wordsSlice[j].Value
	})

	result := make([]string, 0)
	for _, val := range wordsSlice {
		result = append(result, val.Key)

		if len(result) == 10 {
			break
		}
	}

	return result
}
