// database.go

package db

import (
	"database/sql"
//	"time"
//	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"../structs"
)

/* Wow, this site is useful for go db stuff
 * https://www.alexedwards.net/blog/organising-database-access
 */

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

/*
	Task_id     int
	Timestamp   time.Time
	Deleted     bool
	Description string
	Ip_address  string
	Stamp       string


	`task_id` int(11) NOT NULL AUTO_INCREMENT,
	`timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`deleted` tinyint(1) NOT NULL DEFAULT '0',
	`description` text,
	`ip_address` varchar(40) DEFAULT NULL,
	`stamp` text DEFAULT NULL,
	PRIMARY KEY (`task_id`)
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

func SelectAllVisibleTaskDescriptions(db *sql.DB) (*structs.TaskList, error) {
	myTaskList := new(structs.TaskList)
	rows, err := db.Query("SELECT description FROM tasks WHERE deleted=0")
	if err != nil {
		return nil, err
	}

	var task structs.Task
	for rows.Next() {

		err = rows.Scan( &task.Description )
		if err != nil {
			panic(err.Error())
		}

		myTaskList.List = append(myTaskList.List, task)
	}

	return myTaskList, nil
}

/*

func main() { /* ------- MAIN ------- 
	
	db = MysqlConnect()

	defer db.Close()

	// Challenge form, suggested
	  <difficulty>:<timestamp>:<authentication>:<userdefined>
	*   4:20190509T071655Z:c09799c24f7d50f4:1KomQXRI9  <- valid stamp, where secret=="sendtokimball"
	*   difficulty = number of zeros needed
	*   timestamp = date+time as Format("20060102T150405Z")
	*	authentication = sha256("<difficulty>:<timestamp>:<secret>") truncated to the first 16 hex characters
	*	userdefined = any value that matches /[a-zA-Z0-9]{0,64}/
	


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

	fmt.Printf("%d, %s, %b, %s, %s, %s", task.Task_id, task.Timestamp.Format(time.RFC3339), task.Deleted, task.Description, task.Ip_address, task.Stamp)

}


*/