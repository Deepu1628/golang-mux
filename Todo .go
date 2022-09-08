package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//The only important thing is that you then convert from string to int and vice versa in the right places
//and also save after the delete
type Todo struct {
	//and if you set the strtuct at ID to string and adapt the code for it
	Id   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type ResponseStruct struct {
	Message string `json:"message"`
}

func main() {
	fmt.Println("Todo List Backend")
	fileCreation()
	createRoutes()

}
func createRoutes() {

	// rout handlers / Endpoint
	router := mux.NewRouter()
	router.HandleFunc("/todos", GetAllTodos).Methods("GET")
	router.HandleFunc("/todos", PostTodo).Methods("POST")
	router.HandleFunc("/todos", UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{todo_id}", DeleteTodo).Methods("DELETE")

	router.Use(mux.CORSMethodMiddleware(router))

	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), handlers.AllowedOrigins([]string{"*"}))(router)))

	//handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	//handlers.AllowedOrigins([]string{"*"})
	//log.Fatal(http.ListenAndServe(":5000", handlers.CORS()(router)))

}
func DeleteTodo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") // handler for cors
	todoId := mux.Vars(r)["todo_id"]                   //get todo id from url
	todoIdInt, _ := strconv.Atoi(todoId)               // convert the todo id to int
	todos := GetAllTodosHelper()                       // this function returns all todos in the file
	//todos = todos[1:]                                  // i tell just leave 0 and start from 1
	for _, todo := range todos {
		fmt.Printf("%+v\n", todo)

	}
	var indexToDelete int
	for ind, todo := range todos {
		todoIdx, _ := strconv.Atoi(todo.Id)

		if todoIdx == todoIdInt {
			indexToDelete = ind // this is the index which has to be deleted from the list
			break
		}
	}
	var newTodoList []Todo
	newTodoList = todos[:indexToDelete] //appending the list of todos to a new list excluding the one to be deleted
	newTodoList = append(newTodoList, todos[(indexToDelete+1):]...)
	_ = TodoCreateHelper(newTodoList)      //inserting new list of todos into csv
	json.NewEncoder(w).Encode(newTodoList) //returning response
}

func UpdateTodo(w http.ResponseWriter, h *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	handlers.AllowedOrigins([]string{"*"})
	var request Todo
	json.NewDecoder(h.Body).Decode(&request) // Decode the request body to find the todo id
	todos := GetAllTodosHelper()             // get all todos in the file

	for ind, todo := range todos {
		if todo.Id == request.Id {
			todos[ind].Done = request.Done
		}
	}

	json.NewEncoder(w).Encode(todos) // sending response of all todos in the file
	//log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
}

// hello

func PostTodo(w http.ResponseWriter, h *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	handlers.AllowedOrigins([]string{"*"})

	todos := GetAllTodosHelper() // get all existing todos
	lastInd := len(todos) - 1
	if lastInd < 0 {
		lastInd = 0
	}
	newId, _ := strconv.Atoi(todos[lastInd].Id) // find the id for the new todo
	var request Todo
	json.NewDecoder(h.Body).Decode(&request) // decoding the request to get the data
	request.Id = strconv.Itoa(newId + 1)
	todos = append(todos, request)
	_ = TodoCreateHelper(todos)      // creating the new todo and adding it to file
	json.NewEncoder(w).Encode(todos) // sending response of all todos in the file

}

//
func TodoCreateHelper(todos []Todo) ResponseStruct {

	csvFile, _ := os.Create("database.csv") // create csv file
	csvFile.WriteString("id,name,done")     // insert heading
	writer := csv.NewWriter(csvFile)
	var response ResponseStruct

	for ind, todo := range todos {
		var strArr []string

		if ind == 0 {
			_ = writer.Write(strArr) // add empty line below heading
		}
		strArr = append(strArr, fmt.Sprint(todo.Id))
		strArr = append(strArr, todo.Name)
		strArr = append(strArr, fmt.Sprint(todo.Done))
		err := writer.Write(strArr) // write the current todo to csv file
		if err != nil {
			response.Message = "Couldn't create a todo"
		} else {
			response.Message = "Todo created successfully"
		}

	}

	writer.Flush() // writes all the content to the file which are in the object "writer"

	return response
}

// 3
func GetAllTodos(w http.ResponseWriter, h *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //CONNECTION COrs policy
	handlers.AllowedOrigins([]string{"*"})             // cors policy

	todos := GetAllTodosHelper() // Get all todos inside the file
	newTodos := todos
	json.NewEncoder(w).Encode(newTodos)

}

func GetAllTodosHelper() []Todo {

	file, _ := os.Open("database.csv") // opens the file "database.csv" in read mode
	csvReader := csv.NewReader(file)   // csv reader to read the contents of csv file
	data, err := csvReader.ReadAll()   // read all returns the contents of the csv file in multi-dimensional string array format
	if err != nil {                    //not equal operator
		log.Fatal(err) // print
	}
	var todaArr []Todo
	//Loop through books and find ID
	for _, str := range data {

		var todo Todo
		todo.Id = str[0]                         // converts type  string to int
		todo.Name = str[1]                       //noooo
		todo.Done, _ = strconv.ParseBool(str[2]) // returns the bool value represented by string
		todaArr = append(todaArr, todo)
	}

	// converts all the multi-dimensional array contents to struct and returns it
	return todaArr[1:]

}

func fileCreation() *os.File {

	var file *os.File
	if _, err := os.Stat("database.csv"); err == nil { // checks wheather file is present locally
		fmt.Println("File already exists")
		file, _ := os.OpenFile("database.csv", os.O_CREATE, 0666) // opens file in create mode (default)
		return file
	} else {
		fmt.Println("File not found creating one")
		file, _ = os.Create("database.csv")
		file.WriteString("id,name,done")

	}

	return file
}
