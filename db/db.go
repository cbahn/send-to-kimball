// database.go

/* Wow, this site is some useful db stuff
 * https://www.alexedwards.net/blog/organising-database-access
 */

package db

import (
	"../structs"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/* Initiates the connection to a localhost database.
 * Yep, I should really move some of this connection information somewhere else :-\
 */
func MysqlLocalConnect(dbName string, username string, password string) (db *sql.DB) {

	db, err := sql.Open("mysql", username+":"+password+"@/"+dbName+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}

	return db
}

/* returns one arbitrary task. Useful for testing
 */
func GetExampleTask(db *sql.DB) (*structs.Task, error) {

	mytask := new(structs.Task)

	rows, err := db.Query("SELECT task_id, timestamp, deleted, description, ip_address, stamp FROM tasks LIMIT 1")

	if err != nil {
		return nil, err
	}

	rows.Next()

	err = rows.Scan(&mytask.Task_id, &mytask.Timestamp, &mytask.Deleted, &mytask.Description, &mytask.Ip_address, &mytask.Stamp)
	if err != nil {
		return nil, err
	}

	return mytask, nil
}

func InsertNewTask(db *sql.DB, ipAddress string, description string, stamp string) error {
	_, err := db.Query("INSERT INTO tasks (ip_address,description,stamp) VALUES (?,?,?)", ipAddress, description, stamp)
	return err
}

/* Should this have such a long, descriptive name? I'm not sure
 */
func SelectAllVisibleTaskDescriptions(db *sql.DB) (*structs.TaskList, error) {
	myTaskList := new(structs.TaskList)
	rows, err := db.Query("SELECT description FROM tasks WHERE deleted=0 ORDER BY task_id DESC")
	if err != nil {
		return nil, err
	}

	var task structs.Task
	for rows.Next() {

		err = rows.Scan(&task.Description)
		if err != nil {
			panic(err.Error())
		}

		myTaskList.List = append(myTaskList.List, task)
	}

	return myTaskList, nil
}

/* Returns the number of successful posts within the last 'since' minutes.
 * I kind of want this to take a time.Time parameter, but it doesn't sound
 * like fun to translate that into SQL.
 */
func NumberOfPostsSince(db *sql.DB, since int) int {
	rows, err := db.Query("SELECT COUNT(*) FROM tasks WHERE timestamp >= DATE_SUB(NOW(), INTERVAL ? MINUTE)", since)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Must be called once before scanning
	rows.Next()

	var numberOfResults int
	rows.Scan(&numberOfResults)

	// Wow, this was easy
	return numberOfResults
}