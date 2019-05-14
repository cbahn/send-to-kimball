// Definitions of structs shared by multiple packages
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

// It is more convenient to encapsulate this list before handing it 
// to the templating engine
type TaskList struct {
	List []Task
}
