package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	fileHome = `C:\Users\PANKAJ\OneDrive\Desktop\coding\GO\file`
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type model struct {
	newfile     textinput.Model
	data        textarea.Model
	dataLen		int
	dataEditing bool
	isEditing   bool
	currentFile *os.File
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
			m.isEditing = true
			// m.dataEditing= true
		
			return m, nil
		case "enter":

			if m.currentFile!=nil{
				break
			}

			// file creation
			if m.newfile.Value() != "" {

				loc := filepath.Join(fileHome, m.newfile.Value())

				if _, err := os.Stat(loc); err == nil {
					return m, nil
				}

				f, err := os.Create(loc)
				if err != nil {
					log.Fatal(err)
				}
				m.currentFile = f
				m.isEditing = false
				// m.dataEditing=true
				m.newfile.SetValue("")
				// m.data.Focus()
				// fmt.Println("file creates at ", loc)
			}
			return m, nil
		case "ctrl+s":

			if m.currentFile==nil{
				break
			}
			if err:=m.currentFile.Truncate(0); err!=nil{
				fmt.Println("Can't save the file")
				return m,nil
			}

			if _,err:=m.currentFile.Seek(0,0); err!=nil{
				fmt.Println("can not save the file")
				return m,nil
			}


			if m.data.Value()!=""{
				filedata:=m.data.Value()
				len,err:=m.currentFile.Write([]byte(filedata))
				if err!=nil{
					log.Fatal(err)
				}
				m.dataLen=len
				m.dataEditing=false
				// m.data.SetValue("")
			}

			if err:=m.currentFile.Close();err!=nil{
				fmt.Println("can not close the file")
			}

			m.currentFile=nil
			m.data.SetValue("")
			return m,nil
		}

	}
	if m.isEditing {
		m.newfile, cmd = m.newfile.Update(msg)
	}
	if m.currentFile!=nil{
		m.data,cmd=m.data.Update(msg)
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
	if m.isEditing {
		view = m.newfile.View()
	}
	if m.currentFile!=nil{
		view=m.data.View()
		}

	help := "ctrl+c : exit ctrl+n : new file ctrl+s : save the file"
	// style.Render(welcome)
	return fmt.Sprintf("\n%s\n\n%s\n\n%s", welcome, view, help)
}

func initailizeModel() model {

	// we will use os package to create a new file at this locaion
	err := os.MkdirAll(fileHome, 0755)
	if err != nil {
		log.Fatal(err)
	}


	// we are initailizing new file with name

	ti := textinput.New()
	ti.Placeholder = "File name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.PromptStyle.Blink(true)
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = cursorStyle

	// we are making the data input in the created file
	// text area
	tii := textarea.New()
	tii.Placeholder = "Write your note hear"
	tii.ShowLineNumbers=false
	tii.Focus()
	tii.Cursor.Style = cursorStyle
	tii.FocusedStyle.Prompt = cursorStyle
	tii.FocusedStyle.CursorLineNumber.Blink(true)

	return model{
		newfile:   ti,
		isEditing: false,
		data: tii,
		dataEditing: false,
		dataLen: 0,
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
