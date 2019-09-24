package util

import (
	"bufio"
	"bugbounty/logger"
	"os"
	"strings"
)

func GetListOfDomainsFromFile(path string) []string {
	lines, err := readLines(path)
	if err != nil {
		logger.DefaultLogger.Error("%+v", err)
	}

	// assumes that it's a CSV with host,ip (ignoring ip)
	var hosts []string
	for _, line := range lines {
		hosts = append(hosts, strings.Split(line, ",")[0])
	}

	return hosts
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}