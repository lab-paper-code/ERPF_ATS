package types

import (
	"regexp"
	"strconv"
	"strings"
)

func SizeStringToNum(s string) uint64 {
	regExStr := "^([0-9]+)[ ]*([tTgGmMkK]*[bB]?)$"

	reg := regexp.MustCompile(regExStr)
	matches := reg.FindStringSubmatch(s)

	np := uint64(0)
	if len(matches) >= 2 {
		size, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			return np
		}
		np = size
	}

	if len(matches) >= 3 {
		str := strings.ToLower(matches[2])
		str = strings.TrimSuffix(str, "b")
		switch str {
		case "t":
			np *= 1024
			fallthrough
		case "g":
			np *= 1024
			fallthrough
		case "m":
			np *= 1024
			fallthrough
		case "k":
			np *= 1024
		}
	}

	return np
}
