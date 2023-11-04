package main

import (
  "go.bug.bz/todo"

	"fmt"
	"os"
	"strings"
	
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// tab components
	"go.bug.bz/term-todo/itemstab"
)

const (
	windowHeight = 20
	windowWidth = 80
)

type model struct {
	// Tabs structure
	tabs []string
	activeTab int

	// Tab 0: Account
	account string // the current todo account
	accountInput textinput.Model

	// Tab 1: Items
	itemsTab itemstab.Model

	// Tab 2: Add new Item
	descriptionTextarea textarea.Model
}

// Initial state of application
func initialModel() model {
	// Tab 0
	accountInput := textinput.New()
	accountInput.Focus()
	accountInput.CharLimit = 12
	accountInput.Width = 12
	
	// Tab 2
	descriptionTextarea := textarea.New()
	descriptionTextarea.Placeholder = "What do you need to do?"
	descriptionTextarea.Focus()

	descriptionTextarea.Prompt = "┃ "
	descriptionTextarea.CharLimit = 100
	descriptionTextarea.SetWidth(30) // replace with const from todo library
	descriptionTextarea.SetHeight(1)
	descriptionTextarea.ShowLineNumbers = false

	// Create the initial state of the application
	return model{
		// Start with only the account tab
		tabs: []string{"Account"},
		activeTab: 0,
		
		// Account is blank initially
		account: "",
		accountInput: accountInput,
		
		// Items are not loaded until account is selected
		itemsTab: itemstab.New(),

		// Components and state of the new item form
		descriptionTextarea: descriptionTextarea, 
	}
}

func (m model) Init() tea.Cmd {
	if(m.activeTab == 0) {
		return textinput.Blink
	} else if(m.activeTab == 2) {
		return textarea.Blink
	}

	return nil
}

// Cmds and msgs for todo functionality

// Load the account
type loadAccountMsg string
func loadAccount(accountName string) tea.Cmd {
	todo.Start(accountName)
	
	return func() tea.Msg {
		return loadAccountMsg(accountName)
	}
}

// Toggle whether a todo list item is done
func toggleTodo(description string) tea.Cmd {
	selectedTodo, _ := todo.FindTodo(description)
	selectedTodo.Done = !selectedTodo.Done
	todo.UpdateTodo(selectedTodo)
	
	return func() tea.Msg {
		return itemstab.UpdateTodoMsg(-1)
	}
}

// Add a todo list item
type resetAddFormMsg int // Resets the add form (clears inputs, etc)
func addTodo(item todo.TodoItem) tea.Cmd {
	todo.AddTodo(&item)
	
	return func() tea.Msg {
		return resetAddFormMsg(-1)
	}
}

// Delete a todo list item
func deleteTodo(description string) tea.Cmd {
	selectedTodo, _ := todo.FindTodo(description)
	todo.DeleteTodo(selectedTodo)

	return func() tea.Msg {
		return itemstab.UpdateTodoMsg(-1)
	}
}

// Update loop
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
		
	// Handle todo messages
	case loadAccountMsg:
		m.account = string(msg)
		return m, itemstab.UpdateTodoList
	case itemstab.UpdateTodoMsg:
		m.itemsTab, cmd = m.itemsTab.Update(msg)
		
		if len(m.tabs) < 3 {
			// Add todo tab and switch to it to view items
			m.tabs = []string{"Account", "@TODO", "New @TODO"}
			m.activeTab = 1
		} else if m.activeTab != 1 {
			m.activeTab = 1
		}
		
		return m, cmd
	case resetAddFormMsg:
		m.descriptionTextarea.Reset()
		return m, itemstab.UpdateTodoList
		
	// Handle keypresses
	case tea.KeyMsg:
		switch msg.Type {

		// Quit the program
		case tea.KeyCtrlC:
			return m, tea.Quit

		// Handle enter key
		case tea.KeyEnter:
			if m.activeTab == 0{
				// Load account (and initialize items)
				return m, loadAccount(m.accountInput.Value())
			} else if m.activeTab == 1 {
				// Toggle if the todo list item is done
				return m, toggleTodo(m.itemsTab.SelectedItem())
			} else if m.activeTab == 2 {
				// Submit add item form
				newItem := todo.TodoItem{Description: m.descriptionTextarea.Value()}
				return m, addTodo(newItem)
			}
			break;

		// Handle delete/backspace keys
		case tea.KeyDelete:
		case tea.KeyBackspace:
			if m.activeTab == 1 {
				return m, deleteTodo(m.itemsTab.SelectedItem())
			}
			break;

		// Handle tab key
		case tea.KeyTab:
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			return m, cmd
		}
	}

	// If none of the casse above, then propagate msg down to children
	if(m.activeTab == 0) {
		m.accountInput, cmd = m.accountInput.Update(msg)
	}

	if(m.activeTab == 1) {
		m.itemsTab, cmd = m.itemsTab.Update(msg)
	}

	if(m.activeTab == 2) {
		m.descriptionTextarea, cmd = m.descriptionTextarea.Update(msg)
	}

	return m, cmd
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	extraPaddingStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(highlightColor).BorderTop(false).BorderLeft(false).BorderRight(false)
	windowStyle       = lipgloss.NewStyle().Height(windowHeight).Width(windowWidth).BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m model) View() string {
	// First, create the tabs structure
	tabs := m.tabs
	tabContent := []string{}
	
	// Account tab
	if(m.account == "") {
		tabContent = append(tabContent, fmt.Sprintf(
			"Account name \n%s", 
			m.accountInput.View(),
		))
	} else {
		tabContent = append(tabContent, fmt.Sprintf(
			"Viewing todo for account: %s",
			m.account,
		))
	}

	// Display todo items tab
	if(len(tabs) > 1) { // we have an items tab
		tabContent = append(tabContent, fmt.Sprintf(
			"%s's todo list:\n\n%s",
			m.account,
			m.itemsTab.View(),
		))
	} else {
		tabContent = append(tabContent, "Input an account in the Accounts tab to view items")
	}

	// Add new todo item tab
	if(len(tabs) > 2) {
		tabContent = append(tabContent, fmt.Sprintf(
			"Add new todo item:\n\n%s\n",
			m.descriptionTextarea.View(),
		))
	}

	// Render the tabs structure
	doc := strings.Builder{}
	var renderedTabs []string

	for i, t := range tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isLast{
			border.BottomLeft = "│"
			border.BottomRight = "└"
		} else	if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast {
			if isActive {
				border.BottomRight = "└"
			} else {
				border.BottomRight = "┴"
			}
		} 
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	renderedTabs = append(renderedTabs, extraPaddingStyle.Width(windowWidth + 1 - lipgloss.Width(row)).Height(lipgloss.Height(row) - 1).Render(""))
	row = lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Render(tabContent[m.activeTab]))
	
	return docStyle.Render(doc.String())
}

// Main function
func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
