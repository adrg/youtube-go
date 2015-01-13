package youtube

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	decodeFuncNameReg     = regexp.MustCompile(`\.sig\|\|([a-zA-Z0-9]+)\(`)
	decodeFuncBodyPattern = `var\s+..={(.+?)};function\s+%s\([^)]*\){(.+?)}`

	decodeFuncTransReg = regexp.MustCompile(`(..):function\([^)]*\){([^}]*)}`)
	decodeFuncRulesReg = regexp.MustCompile(`..\.(..)\([^,]+,(\d+)\)`)
)

var decodeTransformationList [][]string = [][]string{
	[]string{".reverse(", "reverse"},
	[]string{".splice(", "slice"},
	[]string{"var c=", "swap"},
}

func parseFormats(streamMap string, jsData []byte) ([]*VideoFormat, error) {
	formats := []*VideoFormat{}

	fmtStreams := strings.Split(streamMap, ",")
	for _, fmtStream := range fmtStreams {
		query, err := url.ParseQuery(fmtStream)
		if err != nil {
			continue
		}

		fmtItag, _ := strconv.Atoi(query.Get("itag"))

		fmtType := query.Get("type")
		if fmtType == "" {
			continue
		}

		fmtQuality := query.Get("quality")
		if fmtQuality == "" {
			continue
		}

		fmtURL := query.Get("url")
		if fmtURL == "" {
			continue
		}

		signature := query.Get("sig")
		if signature == "" {
			if jsData == nil {
				continue
			}

			if signature = query.Get("s"); signature != "" {
				signature, err = decodeSignature(jsData, signature)
				if err != nil {
					continue
				}
			}
		}

		format := &VideoFormat{
			Itag:    fmtItag,
			Type:    fmtType,
			Quality: fmtQuality,
			URL:     fmt.Sprintf("%s&signature=%s", fmtURL, signature),
		}

		formats = append(formats, format)
	}

	if len(formats) == 0 {
		return nil, errors.New("Could not retrieve video formats")
	}

	return formats, nil
}

func decodeSignature(playerData []byte, signature string) (string, error) {
	if playerData == nil || len(playerData) == 0 {
		return "", errors.New("Could not decode signature: player data empty")
	}

	// Find the name of the signature decoding function
	lines := []string{}
	decodeFuncName := ""

	scanner := bufio.NewScanner(strings.NewReader(string(playerData)))
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		matches := decodeFuncNameReg.FindStringSubmatch(line)
		if matches != nil && len(matches) > 1 {
			decodeFuncName = matches[1]
			break
		}
	}

	if decodeFuncName == "" {
		return "", errors.New("Could not find decode signature function name")
	}

	// Find transformations and rules in the signature decoding function body
	decodeFuncPattern := fmt.Sprintf(decodeFuncBodyPattern, decodeFuncName)
	decodeFuncBodyReg := regexp.MustCompile(decodeFuncPattern)

	transformations, rules := "", ""
	for _, line := range lines {
		matches := decodeFuncBodyReg.FindStringSubmatch(line)
		if matches != nil && len(matches) > 2 {
			transformations, rules = matches[1], matches[2]
			break
		}
	}

	if transformations == "" || rules == "" {
		return "", errors.New("Could not read signature transformations")
	}

	// Parse transformations
	matches := decodeFuncTransReg.FindAllStringSubmatch(transformations, -1)
	if matches == nil {
		return "", errors.New("Could not find signature transformations")
	}

	trans := map[string]string{}
	for _, match := range matches {
		if len(match) <= 2 {
			return "", errors.New("Could not parse signature transformations")
		}

		method, operation := match[1], match[2]

		transformationFound := false
		for _, transformationPair := range decodeTransformationList {
			if strings.Contains(operation, transformationPair[0]) {
				trans[method] = transformationPair[1]
				transformationFound = true
				break
			}
		}

		if !transformationFound {
			return "", errors.New("Unknown tranformation found")
		}
	}

	if len(trans) == 0 {
		return "", errors.New("Could not parse any signature transformations")
	}

	// Parse rules and apply transformations
	matches = decodeFuncRulesReg.FindAllStringSubmatch(rules, -1)
	if matches == nil {
		return "", errors.New("Could not find signature rules")
	}

	for _, match := range matches {
		if len(match) <= 2 {
			return "", errors.New("Could not parse signature rules")
		}

		index, err := strconv.Atoi(match[2])
		if err != nil {
			return "", errors.New("Could not parse transformation index")
		}

		transformation, ok := trans[match[1]]
		if !ok {
			return "", errors.New("Invalid transformation method")
		}

		switch transformation {
		case "reverse":
			signature = reverseString(signature)

		case "slice":
			signature = signature[index:]

		case "swap":
			runes := []rune(signature)
			c := runes[0]
			runes[0] = runes[index%len(runes)]
			runes[index] = c
			signature = string(runes)
		}
	}

	return signature, nil
}
