package cloudlink

import (
	// "errors"
	"fmt"
)

/*
func containsValueClientlist(slice map[snowflake.ID]*Client, value interface{}) bool {
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
*/

// Generates a value for client identification.
func (client *Client) GenerateUserObject() *UserObject {
	client.RLock()
	defer client.RUnlock()
	if client.username != nil {
		return &UserObject{
			Id:       fmt.Sprint(client.id),
			Username: client.username,
			Uuid:     fmt.Sprint(client.uuid),
		}
	} else {
		return &UserObject{
			Id:   fmt.Sprint(client.id),
			Uuid: fmt.Sprint(client.uuid),
		}
	}
}

func (room *Room) GenerateUserList() []*UserObject {
	var output []*UserObject
	for _, client := range room.clients {

		// Read attributes
		client.RLock()
		usernameset := (client.username != nil)
		protocol := client.protocol
		client.RUnlock()

		// Require a set username and a compatible protocol
		if !usernameset || (protocol != 1) {
			continue
		}

		// Add to userlist if valid
		output = append(output, client.GenerateUserObject())
	}
	return output
}

// Creates a temporary deep copy of a client's rooms map attribute.
func TempCopyRooms(origin map[interface{}]*Room) map[interface{}]*Room {
	clone := make(map[interface{}]*Room, len(origin))
	for x, y := range origin {
		clone[x] = y
	}
	return clone
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

/*
func appendToSlice(slice []interface{}, elements ...interface{}) ([]interface{}, error) {
	// Use the ellipsis (...) to pass multiple elements as arguments to append
	newSlice := append(slice, elements...)

	// Check if the length of the new slice is as expected
	if len(newSlice) != len(slice)+len(elements) {
		return nil, errors.New("failed to append elements to slice")
	}

	return newSlice, nil
}
*/
