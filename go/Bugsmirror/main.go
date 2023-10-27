package main

import (
	"github.com/gin-gonic/gin"
	"errors"
	"net/http"
)

type todo struct{
	//how data structure will look using struct
	//define properties and IDs for unique identification each todo item
	ID 			string `JSON:"id"`
	First_Name 		string `JSON:"fname"`
	Last_Name 	string `JSON:"lname"`
	DOB string `JSON: "dob"`
	Email string `JSON:"email"`
	Mobile string `JSON:"mobile"`
	Completed bool `JSON: "completed"`


}

var todos = []todo{
	{ID: "1", First_Name :"Harsh",Last_Name: "Mavi", DOB: "20/05/2002", Email: "mavi@gmail.com", Mobile: "326598",Completed:true},
	{ID: "2", First_Name :"Nittyansh",Last_Name: "Srivastava", DOB: "23/10/2000",Email: "nsrivas@gmail.com", Mobile: "426598",Completed:true},
	{ID: "3", First_Name :"Rakshita",Last_Name: "Sharma", DOB: "4/8/2002",Email: "sharma@gmail.com", Mobile: "325598",Completed:true},
	{ID: "4", First_Name :"Rashmi",Last_Name: "Bhargava", DOB: "8/11/2002",Email: "rashmi@gmail.com", Mobile: "316598",Completed:true},
	{ID: "5", First_Name :"Pratham",Last_Name: "Jain", DOB: "10/6/2002",Email: "jain@gmail.com", Mobile: "326588",Completed:true},
	{ID: "6", First_Name :"Harsh",Last_Name: "Panwar", DOB: "19/2/2002",Email: "panwar@gmail.com", Mobile: "323598",Completed:true},

}

func getTodos(context *gin.Context){
	context.IndentedJSON(http.StatusOK, todos)
}

func addTodo(context *gin.Context){
	var newTodo todo

	if err := context.BindJSON(&newTodo); err != nil {
		return 
	}
	todos = append(todos,newTodo)
	context.IndentedJSON(http.StatusCreated, newTodo)
}


 

 func getTodo(context *gin.Context){
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil{
		context.IndentedJSON(http.StatusNotFound, gin.H{"message":"todo not found"})
		return 
	}
	context.IndentedJSON(http.StatusOK, todo)
 }

 func getTodoById(id string)(*todo,error){
	for i, t := range todos{
		if t.ID == id{
			return &todos[i], nil
		}
	}

   return nil, errors.New("todo not foound")

}

func toggletodoStatus(context *gin.Context){
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil{
		context.IndentedJSON(http.StatusNotFound, gin.H{"message":"todo not found"})
		return 
	}

	todo.Completed = !todo.Completed

	context.IndentedJSON(http.StatusOK, todo)
}


func main(){
	router := gin.Default() //router is server
	router.GET("/todos",getTodos)
	router.GET("/todos/:id",getTodo)
	router.PATCH("/todos/:id",toggletodoStatus)
	router.POST("/todos",addTodo)
	router.Run("localhost:9090")
}