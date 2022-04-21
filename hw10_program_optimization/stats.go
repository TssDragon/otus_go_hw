package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

const splitEmailPartsCount = 2

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return parseDataAndComputeStat(r, domain), nil
}

func parseDataAndComputeStat(r io.Reader, firstLevelDomain string) (result DomainStat) {
	result = make(DomainStat, 100)

	var jsonParser fastjson.Parser
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		val, err := jsonParser.ParseBytes(scanner.Bytes())
		if err != nil {
			fmt.Printf("err appears while json parsing: %s", err)
			continue
		}

		email := string(val.GetStringBytes("Email"))
		if strings.Contains(email, firstLevelDomain) {
			fullDomain := strings.ToLower(strings.SplitN(email, "@", splitEmailPartsCount)[1])

			_, keyExists := result[fullDomain]
			if !keyExists {
				result[fullDomain] = 0
			}
			result[fullDomain]++
		}
	}
	return
}
