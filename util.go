package gui

func growSlice[Type any](slice []Type, newLen int) []Type {
	newSlice := make([]Type, newLen)
	copy(newSlice, slice)
	return newSlice[:len(slice)]
}
