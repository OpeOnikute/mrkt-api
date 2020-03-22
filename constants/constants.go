package constants

import "fmt"

// NotAdmin ...
var NotAdmin = "This user is not an admin"
var IncorrectCredentials = "The details you entered seem to be incorrect."
var AccessDenied = "You do not have permission to access this resource."
var InvalidParams = "The data you provided is incorrect."
var UserExists = "This user already exists"

// ResourceNotFound ...
func ResourceNotFound(resource string) string {
	return fmt.Sprintf("This %s was not found.", resource)
}
