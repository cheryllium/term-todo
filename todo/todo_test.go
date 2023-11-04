package todo

import (
	"os"
	"testing"
)

// Setup/cleanup for each test case.
// Call by doing defer setupDB()() before each test case.
func setupDB() func() {
	Start("testing")
	return func() {
		err := os.Remove("testing.db")
		if err != nil {
			panic("failed to delete testing.db file")
		}
	}
}

// Test adding a todo list item
func TestAddTodo(t *testing.T) {
	defer setupDB()()
	
	item := TodoItem{Description: "Tickle the cat"}
	AddTodo(&item)
	
	result, err := FindTodo("Tickle the cat")
	if err != nil {
		t.Fatalf("Finding todo failed with error: %v", err)
	}

	if result.Description != "Tickle the cat" {
		t.Fatalf("Got the wrong result from DB, %v", result)
	}

	if result.Done {
		t.Fatalf("Newly added item is marked done by default")
	}
}

// Test deleting a todo list item
func TestDeleteTodo(t *testing.T) {
	defer setupDB()()

	item := TodoItem{Description: "Tickle the cat"}
	AddTodo(&item)

	DeleteTodo(&item)

	_, err := FindTodo("Tickle the cat")
	if err == nil {
		t.Fatalf("Got unexpected result after deleting todo")
	}
}

// Test completing a todo list item
func TestCompleteTodo(t *testing.T) {
	defer setupDB()()

	item := TodoItem{Description: "Tickle the cat"}
	AddTodo(&item)

	item.Done = true

	UpdateTodo(&item)

	result, err := FindTodo("Tickle the cat")
	if err != nil {
		t.Fatalf("Finding failed with error: %v", err)
	}

	if !result.Done {
		t.Fatalf("Todo not marked done after attempting to mark it done")
	}
}

// Test listing all todo list items
func TestListTodos(t *testing.T) {
	defer setupDB()()

	item := TodoItem{Description: "Tickle the cat"}
	AddTodo(&item)

	item2 := TodoItem{Description: "Buy fishy treats"}
	AddTodo(&item2)

	items, err := ListTodos(1) // List page 1 of items

	if err != nil {
		t.Fatalf("Finding failed with error: %v", err)
	}
	
	if len(items) != 2 {
		t.Fatalf("Did not list 2 todo items, listed %v instead", len(items))
	}

	if items[0].Description != "Tickle the cat" {
		t.Fatalf("First item not found in list")
	}

	if items[1].Description != "Buy fishy treats" {
		t.Fatalf("Second item not found in list")
	}
}
