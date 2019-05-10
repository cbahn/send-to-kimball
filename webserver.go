package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"log"
	"net/http"
	"regexp"
	"path/filepath"
	"encoding/json"
	"io"
	"time"
	"github.com/gorilla/mux"
)

func SetMyCookie(response http.ResponseWriter){
	// Add a simplistic cookie to the response.
	cookie := http.Cookie{Name: "testcookiename", Value:"testcookievalue"}
	http.SetCookie(response, &cookie)
}

// Respond to URLs of the form /generic/...
func GenericHandler(response http.ResponseWriter, request *http.Request){

	// Set cookie and MIME type in the HTTP headers.
	SetMyCookie(response)
	response.Header().Set("Content-type", "text/plain")

	// Parse URL and POST data into the request.Form
	err := request.ParseForm()
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	// Send the text diagnostics to the client.
	fmt.Fprint(response,  "FooWebHandler says ... \n")
	fmt.Fprintf(response, " request.Method     '%v'\n", request.Method)
	fmt.Fprintf(response, " request.RequestURI '%v'\n", request.RequestURI)
	fmt.Fprintf(response, " request.URL.Path   '%v'\n", request.URL.Path)
	fmt.Fprintf(response, " request.Form       '%v'\n", request.Form)
	fmt.Fprintf(response, " request.Cookies()  '%v'\n", request.Cookies())
	fmt.Fprintf(response, " request.RemoteAddr '%v'\n", request.RemoteAddr)
}

func AboutHandler(response http.ResponseWriter, request *http.Request){
	http.ServeFile(response, request, "about.html")
}

// Respond to the URL /home with an html home page
func HomeHandler(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-type", "text/html")
	webpage, err := ioutil.ReadFile("index.html")
	if err != nil { 
		http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)
	}
	fmt.Fprint(response, string(webpage));
	fmt.Println("Sent response to /home")
}

// Serves the vote.html file
func VoteGETHandler(response http.ResponseWriter, request *http.Request){
	http.ServeFile(response, request, "vote.html")
}



// A utility function for converting the request.Body into a map[string]string.
/* This is quite fragile. If the json has a non-string type in it then the marshall fails.
   In general, converting JSON into a go object is awkward because of go's static typing.
   We could write our own Unmarshaler but I'm going to stick with a weaker, simpler solution
   and pass the numeric values as strings then convert them.
*/// https://engineering.bitnami.com/articles/dealing-with-json-with-non-homogeneous-types-in-go.html
func jsonReaderToMap(jsonReader io.ReadCloser) (map[string]string, error) {
	jsonBytes, err :=  ioutil.ReadAll(jsonReader)
	if err != nil {
		panic(err) // What could possibilty go wrong?
	}

	jsonMap := make(map[string]string)
	err = json.Unmarshal(jsonBytes, &jsonMap)

	// Whatever the error was, return it
	// They gotta deal with that shit upstream
	return jsonMap, err
}


// This recieves votes as POST requests to /vote and records them to the database
func VotePOSTHandler(response http.ResponseWriter, request *http.Request){

	// Render the raw post into postData, of type map[string]string
	postData, err := jsonReaderToMap(request.Body)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	// Count up the total number of votes, and also validate the vote json
	totalVotes := 0
	for _, v := range postData {
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("Error while converting POST json data to int: ",err)
			return;
		}
		totalVotes += i
	}

	// Read usernamecookie into string username
	// If error, default to ""
	var username string
	if usernameCookie, err := request.Cookie("username"); err == nil {
		// Validate the username's value against our requirements
		var re = regexp.MustCompile(`^[A-Za-z0-9 ]{1,20}$`)
		if re.MatchString(usernameCookie.Value) {
			username = usernameCookie.Value
		}
	}

	fmt.Printf("%s: %s,%s - %d votes\n",time.Now().Format("20060102150405"),request.RemoteAddr,username,totalVotes)

	// Send a string as a response
	response.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(response, "{\"reponse\":\"POST recieved - success! You voted %d times\"}", totalVotes)
}

// Loads up files from the /res folder
// WARNING - ALL FILES IN THAT FOLDER WILL BE PUBLIC
func ResHandler(response http.ResponseWriter, request *http.Request){
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
	http.ServeFile(response, request, filepath.Join(resourceFolder,resource[1]) )
}


// Loads up files from the /res/pic folder
// WARNING - ALL FILES IN THAT FOLDER WILL BE PUBLIC
func PicHandler(response http.ResponseWriter, request *http.Request){
	resourceFolder := "res/pic"
	// Only resources with characters from a-z, A-Z, 0-9, and the _ (underscore) character will be valid.
	var resURL = regexp.MustCompile(`^/res/pic/(\w+\.\w+)$`) 
	var resource = resURL.FindStringSubmatch(request.URL.Path)
	// resource is captured regex matches i.e. ["/res/file.txt", "file.txt"]

	if len(resource) == 0 { // If url could not be parsed, send 404
		fmt.Println("Could not parse /res request:", request.URL.Path)
		http.Error(response, "404 page not found", 404)
		return
	}

	// Everything's good, serve up the file
	http.ServeFile(response, request, filepath.Join(resourceFolder,resource[1]) )
}

// This recieves votes as POST requests to /vote and records them to the database
func SendHandler(response http.ResponseWriter, request *http.Request){

	request.ParseForm()

	for k,v := range request.Form {
		fmt.Printf("%s = %s\n",k,v)
	}

	fmt.Fprintln(response,"Registration successful. Your inner man is now aligned with nature")
}

func ListHandler(response http.ResponseWriter, request *http.Request){
	// Respond to the index url

	response.Header().Set("Content-type", "text/html")
	webpage, err := ioutil.ReadFile("list.html")
	if err != nil { 
		http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)
	}
	fmt.Fprint(response, string(webpage));
	fmt.Println("Sent response to /home")

}


func main(){
	port := 8097
	portstring := strconv.Itoa(port)

	// We're using gorilla/mux as the router because
	// it's not garbage like the default one.
	mux := mux.NewRouter()

	/*
	mux.Handle("/generic/", 		http.HandlerFunc( GenericHandler  ))
	mux.Handle("/about",			http.HandlerFunc( AboutHandler    ))
	mux.Handle("/vote",				http.HandlerFunc( VoteGETHandler  )).Methods("GET")
	mux.Handle("/vote",				http.HandlerFunc( VotePOSTHandler )).Methods("POST")
	*/

	mux.Handle("/res/{resource}",	http.HandlerFunc( ResHandler      )).Methods("GET")
	mux.Handle("/res/pic/{picture}",http.HandlerFunc( PicHandler	  )).Methods("GET")
	

	mux.Handle("/list",				http.HandlerFunc( ListHandler     )).Methods("GET")
	mux.Handle("/send",				http.HandlerFunc( SendHandler	  )).Methods("POST")
	mux.Handle("/", 				http.HandlerFunc( HomeHandler     )).Methods("GET")

	// Start listing on a given port with these routes on this server.
	log.Print("Listening on port " + portstring + " ... ")
	err := http.ListenAndServe(":" + portstring, mux)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

