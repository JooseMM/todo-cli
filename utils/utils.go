package utils

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/* Category Enum */
type Category string

const FrontEnd Category = "Front-end"
const BackEnd Category = "Back-end"

type Task struct {
	ID          uint8 `gorm:"primaryKey"`
	Task        string
	Category    Category
	IsCompleted bool
}

func CreateOne(task *string, isFront *bool, isBack *bool, db *gorm.DB) {
	var category Category

	if !*isFront && !*isBack {
		fmt.Println("Error: category is needed use --front, --backend to specify one")
		flag.Usage()
		return
	}

	if *isFront && *isBack {
		fmt.Println("Error: Tasks can have only one category")
		flag.Usage()
		return
	}

	if *isFront {
		category = FrontEnd
	} else {
		category = BackEnd
	}

	if task == nil || *task == "" {
		fmt.Println("Error: task description cannot be empty")
		return
	}

	taskModel := Task{
		Task:        *task,
		IsCompleted: false,
		Category:    category,
	}

	result := db.Create(&taskModel)

	if result.Error != nil {
		println(result.Error)
	} else {
		println(taskModel.ID)
	}
}

func printTask(title string, bank []Task, maxStringLength uint8) {
	reset := "\033[0m"
	yellow := "\033[33m"
	red := "\033[31m"
	green := "\033[32m"
	blue := "\033[34m"
	cya := "\033[36m"
	purple := "\033[35m"
	leftPadding := 5

	fmt.Printf("\n\033[1m%s\n\n", title)

	// Header with fixed width columns
	fmt.Printf(
		"%-*s %-5s%s  %s%-*s%s  %s%-10s%s\n",
		leftPadding, blue, "ID", reset,
		purple, maxStringLength, "Task", reset,
		cya, "Complete", reset,
	)

	line := strings.Repeat("-", int(maxStringLength)+20)
	fmt.Println(line)

	for _, t := range bank {
		status := red + " " + reset
		if t.IsCompleted {
			status = green + " " + reset
		}

		// Print each row with fixed width columns
		fmt.Printf(
			"%-*s %-5d%s  %s%-*s%s  %s%-10s%s\n",
			leftPadding, yellow, t.ID, reset,
			reset, maxStringLength+2, t.Task, reset,
			reset, status, reset,
		)
		fmt.Println(line)
	}
}

func MarkAsComplete(db *gorm.DB, idStr *string) {
	id, err := strconv.Atoi(*idStr)

	if err != nil {
		fmt.Printf("Invalid ID: %v", err)
	}
	error := db.Model(&Task{}).
		Where("id = ?", id).
		Update("is_completed", true).
		Error
	if error != nil {
		fmt.Printf("%v", error)
	}

	fmt.Printf("Task with ID: %d was mark as completed", id)
}

func DeleteOne(db *gorm.DB, idStr *string) {
	id, err := strconv.Atoi(*idStr)
	if err != nil {
		fmt.Printf("invalid ID: %v", err)
	}

	// Delete task by primary key
	result := db.Delete(&Task{}, id)
	if result.Error != nil {
		fmt.Printf("invalid ID: %v", result.Error)
	}
}

func ListAllTasks(db *gorm.DB, isFront *bool, isBack *bool) error {
	var tasks []Task
	var frontTask []Task
	var backTask []Task
	var maxLen uint8 = 0

	result := db.Find(&tasks)

	if result.Error != nil {
		return result.Error
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	for _, t := range tasks {
		var crt uint8 = uint8(len(t.Task))
		if crt > maxLen {
			maxLen = crt
		}

		if t.Category == FrontEnd {
			frontTask = append(frontTask, t)
		} else if t.Category == BackEnd {
			backTask = append(backTask, t)
		}
	}

	if *isBack == false && *isFront == false {
		printTask("Front-end", frontTask, maxLen)
		printTask("Back-end", backTask, maxLen)
	}

	if *isFront {
		printTask("Front-end", frontTask, maxLen)
	}
	if *isBack {
		printTask("Back-end", backTask, maxLen)
	}

	return nil
}

func InitializeDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("todo-list.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Task{})
	return db, nil
}

func CloseDatabaseConnection(db *gorm.DB) error {
	sql, err := db.DB()

	if err != nil {
		return err
	}

	sql.Close()
	return nil
}
