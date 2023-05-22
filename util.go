package main

func containsValue(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
func RemoveValue(slice []string, indexRemove int) []string {
	// Swap the element to remove with the last element
	slice[indexRemove] = slice[len(slice)-1]

	// Remove the last element
	slice = slice[:len(slice)-1]
	return slice
}
func GetValue(slice []string, target string) int {
	for i, value := range slice {
		if value == target {
			return i
		}
	}
	return -1 // Indicates that the value was not found
}
