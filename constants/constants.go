package constants

import "fmt"

// NotAdmin ...
var NotAdmin = "This user is not an admin"
var IncorrectCredentials = "The details you entered seem to be incorrect."
var AccessDenied = "You do not have permission to access this resource."

// ResourceNotFound ...
func ResourceNotFound(resource string) string {
	return fmt.Sprintf("This %s was not found.", resource)
}
