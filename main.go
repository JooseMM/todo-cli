package main

import (
	"flag"
	"fmt"
	"todo-cli/utils"
)


func main() {
	db, err := utils.InitializeDB()
	defer utils.CloseDatabaseConnection(db)

	if err != nil {
		fmt.Println("Error: unexpected error when initializing the database")
		return
	}

	isBack := flag.Bool("back", false, "Back-end category")
	isFront := flag.Bool("front", false, "Front-end category")

	/* Modes */
	list := flag.Bool("list", false, "List all tasks")
	create := flag.String("create", "", "Create task")
	delete := flag.String("delete", "", "Delete task")
	complete := flag.String("complete", "", "Task status")

	flag.Parse()

	if *create != "" {
		utils.CreateOne(create, isFront, isBack, db)
	}

	if *list {
		utils.ListAllTasks(db, isFront, isBack)
	}

	if *delete != "" {
		utils.DeleteOne(db, delete)
	}

	if *complete != "" {
		utils.MarkAsComplete(db, complete)
	}
}
