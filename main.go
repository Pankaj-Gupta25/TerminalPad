package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	newfile textinput.Model
	fileloc string
	isEditing bool

}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			fmt.Println("user clicled", msg.String())
			return m, tea.Quit
		case "ctrl+n":
			m.isEditing=true
			return m, nil
		case "enter":
			m.isEditing=true
			loc:=filepath.Join(m.fileloc,m.newfile.Value())
			os.Create(loc)
			fmt.Println("file creates at ",loc)
			return m,tea.Quit
		}
		

	}
	if m.isEditing{
		m.newfile,cmd = m.newfile.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		// PaddingTop(2).
		PaddingLeft(2).
		PaddingRight(2)

	welcome := style.Render("Welcome to TermianlPad")

	view := ""
	if m.isEditing{
		view=style.Render(m.newfile.View())
	}

	help := "ctrl+c : exit"
	// style.Render(welcome)
	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}

func initailizeModel() model {

	// we will use os package to create a new file at this locaion
	
	filPath:=`C:\Users\PANKAJ\OneDrive\Desktop\coding\GO\file\`

	// we are initailizing new file input

	ti := textinput.New()
	ti.Placeholder = "File name"
	ti.Focus()
	ti.CharLimit = 156
	// ti.Width = 20
	ti.PromptStyle.Blink(true)

	return model{
		newfile: ti,
		isEditing: false,
		fileloc: filPath,
	}
}

func main() {

	p := tea.NewProgram(initailizeModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("the error while initailizing is ", err)
		os.Exit(1)
	}

	fmt.Println("Termina UI")
}
