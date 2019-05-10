// Package structs contains the base definition of structs shared by multiple utilities
package structs

import (
	"time"
)

type Task struct {
	Task_id     int
	Timestamp   time.Time
	Deleted     bool
	Description string
	Ip_address  string
	Stamp       string
}

// It is nessesary to encapsulate the list of tasks as a struct so that
// it can be handed to the templating engine
type TaskList struct {
	List []Task
}
