// database.go

package main

import (
	"database/sql"
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

/* Wow, this site is useful for go db stuff
 * https://www.alexedwards.net/blog/organising-database-access
 */

type Task struct {
	TaskId int
	Title string
	Date time.Time
}

func MysqlConnect() (db *sql.DB) {
	username := "junco"
	password := "OverboardSkimmer3397"
	dbName := "testdb"

	db, err := sql.Open("mysql", username+":"+password+"@/"+dbName+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}

	return db
}

func GetExampleTask(db *sql.DB) (*Task, error) {

	mytask := new(Task)

	rows, err := db.Query("SELECT * FROM mytable LIMIT 1")

	fmt.Println("query worked fine")

    if err != nil {
        return nil, err
    }

    rows.Next()

	err = rows.Scan(&mytask.TaskId, &mytask.Title, &mytask.Date)
	if err != nil {
		return nil, err
	}

	return mytask, nil
}

func InsertNewTask(db *sql.DB, ipAddress string, description string, stamp string) error {
	_, err := db.Query("INSERT INTO tasks (ip_address,description,stamp) VALUES (?,?,?)", ipAddress, description, stamp)
	return err
}


func main() { /* ------- MAIN ------- */
	
	db := MysqlConnect()

	defer db.Close()

	// Challenge form, suggested
	/*  <difficulty>:<timestamp>:<authentication>:<userdefined>
	*   4:20190509T071655Z:c09799c24f7d50f4:1KomQXRI9  <- valid stamp, where secret=="sendtokimball"
	*   difficulty = number of zeros needed
	*   timestamp = date+time as Format("20060102T150405Z")
	*	authentication = sha256("<difficulty>:<timestamp>:<secret>") truncated to the first 16 hex characters
	*	userdefined = any value that matches /[a-zA-Z0-9]{0,64}/
	*/


	// perform a db.Query insert
	// time is formatted similar to: t.Format("20060102T150405Z")
	// Currently validates using sha256
	err := InsertNewTask(db,"192.168.0.1", "Send to my man kimballin", "nonce:8:20060102150405:LqXBpljZxFafr3")

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Insert successful")
	}
	

	task, err := GetExampleTask(db)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%d, %s, %s", task.TaskId, task.Title, task.Date.Format(time.RFC3339))

}


