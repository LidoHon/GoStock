package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LidoHon/GoStock/router"
	"github.com/joho/godotenv"
)




func main(){
	
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port ==""{
		port =`8080`
	}
 r := router.Router()
 fmt.Println("server running on port %s...\n", port)
 log.Fatal(http.ListenAndServe(":" +port, r))
}