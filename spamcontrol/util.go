package spamcontrol

// unique returns a copy of the given string array
// cleaned from duplicate entries.
func unique(list []string) []string {
	var cleanedList []string

	index := make(map[string]int)
	for _, entry := range list {
		if _, exists := index[entry]; exists {
			continue
		}

		index[entry] = 1

		cleanedList = append(cleanedList, entry)
	}

	return cleanedList
}
