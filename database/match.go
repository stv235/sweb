package database

import "strings"

func matchOr(search interface{}, texts ...string) bool {
	//tokens := strings.Split(search, " ")

	tokens := search.([]string)

	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}

	for i, text := range texts {
		texts[i] = strings.ToLower(text)
	}

	containsToken := func(token string) bool {
		for _, text := range texts {
			if strings.Contains(text, token) {
				return true
			}
		}

		return false
	}

	for _, token := range tokens {
		if !containsToken(token) {
			return false
		}
	}

	return true
}
