package utils

const (
	// Common messages
	InvalidJsonProvided = "inavlid json provided"

	//User
	UserSuccessCreate               = "user created successfully"
	UserFailedGetToken              = "failed to get token"
	UserNotFound                    = "user not found"
	UserInvalidCredentials          = "user invalid credentials"
	UserIdRequired                  = "user id is required"
	UserWithoutAccesPermission      = "user without access permission"
	UserCannotChangeTaskAnotherUser = "user cannot change a task of another user"

	//Task
	TaskNotFound        = "task not found"
	TaskTitleRequired   = "task title is required"
	TaskSummaryRequired = "task summary is required"
)
