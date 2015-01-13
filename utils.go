package youtube

import (
	"io/ioutil"
	"net/http"
)

func reverseString(s string) string {
	n := len(s)
	runes := make([]rune, n)

	for _, r := range s {
		n--
		runes[n] = r
	}

	return string(runes[n:])
}

func getURLData(address string) ([]byte, error) {
	data, err := http.Get(address)
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
