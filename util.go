package main

import "errors"

func containsValueClientlist(slice map[*Client]Client, value interface{}) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func containsValue(slice []interface{}, value interface{}) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
func RemoveValue(slice []interface{}, indexRemove int) []interface{} {
	// Swap the element to remove with the last element
	slice[indexRemove] = slice[len(slice)-1]

	// Remove the last element
	slice = slice[:len(slice)-1]
	return slice
}
func GetValue(slice []interface{}, target interface{}) int {
	for i, value := range slice {
		if value == target {
			return i
		}
	}
	return -1 // Indicates that the value was not found
}
func appendToSlice(slice []interface{}, elements ...interface{}) ([]interface{}, error) {
	// Use the ellipsis (...) to pass multiple elements as arguments to append
	newSlice := append(slice, elements...)

	// Check if the length of the new slice is as expected
	if len(newSlice) != len(slice)+len(elements) {
		return nil, errors.New("failed to append elements to slice")
	}

	return newSlice, nil
}
