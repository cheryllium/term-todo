package itemstab

import (
	"go.bug.bz/todo"
	"fmt"
	
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// Type representing model state
type Model struct {
	items []todo.TodoItem // todo list items
	itemsTable table.Model
}

// Default / initial model state
func New() Model {
	return Model{
		items: []todo.TodoItem{},
		itemsTable: updateItemsTable([]todo.TodoItem{}),
	}
}

// Cmd: Update the list of todo items
type UpdateTodoMsg int
func UpdateTodoList() tea.Msg {
	return UpdateTodoMsg(-1)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case UpdateTodoMsg:
		cursor := m.itemsTable.Cursor()
		m.items, _ = todo.ListTodos(1)
		m.itemsTable = updateItemsTable(m.items)

		if cursor != m.itemsTable.Cursor() {
			if cursor >= len(m.items) {
				m.itemsTable.SetCursor(len(m.items) - 1)
			} else {
				m.itemsTable.SetCursor(cursor)
			}
		}

		return m, cmd
	}

	m.itemsTable, cmd = m.itemsTable.Update(msg)
	
	return m, cmd
}

// Returns identifier for selected item, currently description
func (m Model) SelectedItem() (string) {
	return m.itemsTable.SelectedRow()[1]
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n\n%s\n%s\n%s", 
		m.itemsTable.View(),
		"Use ↑ and ↓ to change selection     ",
		"Press ENTER to toggle selected item ",
		"Press DELETE to delete selected item",
	)
}

func updateItemsTable(items []todo.TodoItem) table.Model {
	columns := []table.Column{
		{Title: "Done", Width: 4},
		{Title: "@TODO", Width: 30},
	}

	rows := []table.Row{}
	for _, item := range items {
		checked := "no"
		if item.Done{
			checked = "yes"
		}
		rows = append(rows, table.Row{
			checked,
			item.Description,
		})
	}
	
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(8),
	)

	return t
}

