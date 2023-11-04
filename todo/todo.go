package todo

import (
	"fmt"
	"errors"
	
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/sqlite"

	"log"
)

const (
	ItemsPerPage = 20
)

var db *gorm.DB = nil
type TodoItem struct {
	gorm.Model
	Description string
	Done bool `gorm:"default:false"`
}

// Starts todo given an identifier (for the db)
func Start(identifier string) {
	// Initialize logger
	log.SetPrefix("todo: ")
	log.SetFlags(0)

	// Initialize database
	db_init(identifier)
}

// Initializes db given a DB name
// (Connects to the database and runs any migrations)
func db_init(name string) {
	// Connect to the Database
	var err error
	db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%v.db", name)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	db.AutoMigrate(&TodoItem{});
}

// Adds a todo list item
func AddTodo(item *TodoItem) {
	db.Create(item)
}

// Deletes a todo list item
func DeleteTodo(item *TodoItem) {
	db.Delete(item)
}

// Updates a todo list item
func UpdateTodo(item *TodoItem) {
	db.Save(&item)
}

// Lists all todo list items
func ListTodos(pageNum int) ([]TodoItem, error) {
	var items []TodoItem
	
	result := db.Find(&items).Limit(ItemsPerPage).Offset(ItemsPerPage * (pageNum - 1))
	if result.Error != nil {
		return nil, errors.New("Database error (possibly no rows found)")
	}
	
	return items, nil
}

// Finds a single todo list item by description
// (returns first match to search on description)
func FindTodo(searchterm string) (*TodoItem, error) {
	var item TodoItem
	
	result := db.Where("Description = ?", searchterm).First(&item)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("No rows found")
		}
		return nil, errors.New("Database error")
	}

	return &item, nil
}
