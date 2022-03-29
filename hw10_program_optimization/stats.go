package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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

func parseDataAndComputeStat(r io.Reader, searchDomain string) (result DomainStat, err error) {
	result = make(DomainStat, 100)
	var p fastjson.Parser

	re, err := regexp.Compile("\\." + searchDomain)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		val, err := p.ParseBytes(scanner.Bytes())
		if err == nil {
			email := val.GetStringBytes("Email")
			if re.Match(email) {
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
