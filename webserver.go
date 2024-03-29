package main

import (
	"./db"
	"./stampmaster"
	"./structs"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"
)

const (
	SECRET = "sendtokimball"
)

// This function defines how the difficulty is computed and used
func computeDifficulty() int {
	/* Right now the difficulty is defined as
	16 + number of posts in the last  2 minutes + number of posts in the last 10 minutes
	Most likely this needs tweaking */

	difficulty := 16
	difficulty += db.NumberOfPostsSince(db_conn, 2)
	difficulty += db.NumberOfPostsSince(db_conn, 10)
	return difficulty
}

// Loads up files from the /res folder
// WARNING - ALL FILES IN THAT FOLDER WILL BE PUBLIC
func ResHandler(response http.ResponseWriter, request *http.Request) {
	resourceFolder := "res"
	// Only resources with characters from a-z, A-Z, 0-9, and the _ (underscore) character will be valid.
	var resURL = regexp.MustCompile(`^/res/(\w+\.\w+)$`)
	var resource = resURL.FindStringSubmatch(request.URL.Path)
	// resource is captured regex matches i.e. ["/res/file.txt", "file.txt"]

	if len(resource) == 0 { // If url could not be parsed, send 404
		fmt.Println("Could not parse /res request:", request.URL.Path)
		http.Error(response, "404 page not found", 404)
		return
	}

	// Everything's good, serve up the file
	http.ServeFile(response, request, filepath.Join(resourceFolder, resource[1]))
}

func ListHandler(response http.ResponseWriter, request *http.Request) {
	// Read values in from the database
	var myTaskList *structs.TaskList
	myTaskList, err := db.SelectAllVisibleTaskDescriptions(db_conn)
	if err != nil {
		http.Error(response, "500 Error reading database", 500)
		return
	}

	// Load in the list template
	// (someday this should be loaded only once at startup)
	t, err := template.ParseFiles("list.tmpl")
	if err != nil {
		http.Error(response, "500 Error could not parse list.html template", 500)
		return
	}

	// Execute template
	response.Header().Set("Content-type", "text/html")
	err = t.Execute(response, *myTaskList)
	if err != nil {
		http.Error(response, "500 Error could not execute list.html template", 500)
		return
	}
}

// This recieves votes as POST requests to /vote and records them to the database
func SendHandler(response http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	for k, v := range request.Form {
		fmt.Printf("%s = %s\n", k, v)
	}

	fmt.Fprintf(response, "In the last 100 minutes, there have been %d posts\n", db.NumberOfPostsSince(db_conn, 100))
	fmt.Fprintf(response, "Registration successful. Your inner man is now aligned with nature\n%s", stampmaster.CreateNewStamp(4, "k").ToString())

	err := db.InsertNewTask(db_conn, request.RemoteAddr, request.FormValue("description"), "this:is:a:stamp")
	if err != nil {
		http.Error(response, "422 Could not process submission", 422)
	}
}

//http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)

// Respond to the URL /home with an html home page
func HomeHandler(response http.ResponseWriter, request *http.Request) {
	templateFile := "index.tmpl"
	// Load in the list template
	// (someday this should be loaded only once at startup)
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(response, fmt.Sprintf("500 Error parsing %s: %v", templateFile, err), 500)
		return
	}

	stampToInsert := stampmaster.CreateNewStamp(computeDifficulty(), SECRET).ToString()

	// Execute template
	response.Header().Set("Content-type", "text/html")
	err = t.Execute(response, stampToInsert)
	if err != nil {
		http.Error(response, fmt.Sprintf("500 Error could not execute template: %v", err), 500)
		return
	}
}

/* ================================================== */

// Making the db variable global feels appropriate because there is only one database
var db_conn *sql.DB

func main() {
	port := 80
	portstring := strconv.Itoa(port)

	// We're using gorilla/mux as the router because
	//  it's not garbage like the default one.
	mux := mux.NewRouter()

	// Establish database connection
	dbName := "testdb"	
	username := "junco"
	password := "OverboardSkimmer3397" // These values should be moved to a separate config file
	db_conn = db.MysqlLocalConnect(dbName,username,password)
	defer db_conn.Close()

	mux.Handle("/res/{resource}", http.HandlerFunc(ResHandler)).Methods("GET")
	mux.Handle("/list", http.HandlerFunc(ListHandler)).Methods("GET")
	mux.Handle("/send", http.HandlerFunc(SendHandler)).Methods("POST")
	mux.Handle("/", http.HandlerFunc(HomeHandler)).Methods("GET")

	// Start listing on a given port with these routes on this server.
	log.Print("Listening on port " + portstring + " ... ")
	err := http.ListenAndServe(":"+portstring, mux)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
