package utilities

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func GetDomainsFromFile(filename string) []string {
	content, err := os.ReadFile(filename)
	CheckError(err)

	re := regexp.MustCompile(`\s+`)
	lines := strings.Split(string(content), "\n")

	var domains = []string{}
	for index := range lines {
		line := re.ReplaceAllString(lines[index], "")
		if line != "" {
			domains = append(domains, line)
		}
	}
	return domains
}
