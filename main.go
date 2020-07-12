package main

// import "os"

func main() {
	mainApp := App{}
	initDB	:= InitDB{}

	initDB.dataInitDB("postgres","123456","packform")
	// mainApp.Initialize(
	// 	os.Getenv("APP_DB_USERNAME"),
	// 	os.Getenv("APP_DB_PASSWORD"),
	// 	os.Getenv("APP_DB_NAME"))

	mainApp.Initialize("postgres","123456","packform")

	mainApp.Run(":8010")
}