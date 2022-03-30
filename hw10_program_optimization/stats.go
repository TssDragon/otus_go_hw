package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat, err := parseDataAndComputeStat(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get emails error: %w", err)
	}
	return domainStat, nil
}

func parseDataAndComputeStat(r io.Reader, firstLevelDomain string) (result DomainStat, err error) {
	result = make(DomainStat, 100)

	var jsonParser fastjson.Parser
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		val, err := jsonParser.ParseBytes(scanner.Bytes())
		if err == nil {
			email := string(val.GetStringBytes("Email"))
			if strings.Contains(email, firstLevelDomain) {
				fullDomain := strings.ToLower(strings.SplitN(string(email), "@", 2)[1])

				_, keyExists := result[fullDomain]
				if !keyExists {
					result[fullDomain] = 0
				}
				result[fullDomain]++
			}
		}
	}
	return
}
