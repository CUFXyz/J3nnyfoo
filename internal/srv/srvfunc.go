package srv

func removebyField(arr []jsonPlaceholder, rmID int) []jsonPlaceholder {
	var result []jsonPlaceholder
	for _, item := range arr {
		if item.Id != rmID {
			result = append(result, item)
		}
	}
	return result
}
