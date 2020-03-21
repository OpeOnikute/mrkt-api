package constants

import "fmt"

// ResourceNotFound ...
func ResourceNotFound(resource string) string {
	return fmt.Sprintf("This %s was not found.", resource)
}
