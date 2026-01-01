package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	fileHome    = `C:\Users\PANKAJ\OneDrive\Desktop\coding\GO\file`
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

// struct for the items of the list
type item struct {
	title, desc string
}
// methods of the items 
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

//model of the things we are doing in the terminal
type model struct {
	newfile     textinput.Model
	data        textarea.Model
	dataLen     int
	dataEditing bool
	isEditing   bool
	currentFile *os.File
	list        list.Model
	showingList bool
}

// methods of the model
func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			fmt.Println("user clicled", msg.String())
			return m, tea.Quit
		
		case "esc":
			
			if m.isEditing{
				m.newfile.SetValue("")
				m.isEditing=false
			}
			if m.currentFile!=nil{
				m.data.SetValue("")
				m.currentFile=nil
			}
			if m.showingList{
				// to exit the filter window in the search
				if m.list.FilterState() == list.Filtering{
					break
				}
				m.showingList=false
			}
			return m,nil

		case "ctrl+l":
			// to show the list
			// updating the list
			itemList:=listFile()
			m.list.SetItems(itemList)
			m.showingList=true
			return m,nil

		case "ctrl+n":
			m.isEditing = true
			// m.dataEditing= true

			return m, nil
		case "enter":

			if m.currentFile != nil {
				break
			}

			if m.showingList{
				item,ok:=m.list.SelectedItem().(item)
				if ok{
					// fileLoc := fmt.Sprintf("%s%s",fileHome,item.title)
					fileLoc := filepath.Join(fileHome,item.title)
					content,err:=os.ReadFile(fileLoc)
					if err!=nil{
						log.Println("Error reading the file ",err)
						return m,nil
					}
					m.data.SetValue(string(content))

					f,err:=os.OpenFile(fileLoc,os.O_RDWR,0644)
					if err!=nil{
						log.Fatal("Can't open the file ",err)
						return m,nil
					}
					m.currentFile=f
					m.showingList=false

				}
				return m,nil
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
			
			if m.currentFile == nil {
				break
			}
			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("Can't save the file")
				return m, nil
			}

			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Println("can not save the file")
				return m, nil
			}

			if m.data.Value() != "" {
				filedata := m.data.Value()
				len, err := m.currentFile.Write([]byte(filedata))
				if err != nil {
					log.Fatal(err)
				}
				m.dataLen = len
				m.dataEditing = false

				// m.data.SetValue("")
				// m.currentFile=nil
			}

			if err := m.currentFile.Close(); err != nil {
				fmt.Println("can not close the file")
			}

			m.data.SetValue("")
			m.currentFile=nil
			return m, nil

		case "ctrl+d":

			if m.showingList{
				item,ok:=m.list.SelectedItem().(item)
				if ok{
					
					fileLoc:=filepath.Join(fileHome,item.title)
					os.Remove(fileLoc)

					id:=m.list.Index()
					f:=m.list.Items()

					f = append(f[:id],f[id+1:]...)
					m.list.SetItems(f)

					// we cant re-inisilise the slice again and again it take more time but can edit the exixitng slice
					
					// itemList:=listFile()
					// m.list.SetItems(itemList)
				}
			}

			return m,nil
		}

	}
	if m.isEditing {
		m.newfile, cmd = m.newfile.Update(msg)
	}
	if m.currentFile != nil {
		m.data, cmd = m.data.Update(msg)
	}
	if m.showingList{
		m.list,cmd=m.list.Update(msg)
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
	if m.currentFile != nil {
		view = m.data.View()
	}
	if m.showingList{
		view=m.list.View()
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
	tii.ShowLineNumbers = false
	tii.Focus()
	tii.Cursor.Style = cursorStyle
	tii.FocusedStyle.Prompt = cursorStyle
	tii.FocusedStyle.CursorLineNumber.Blink(true)

	fileList:=listFile()
	fmt.Println(fileList)

	finalList:=list.New(fileList,list.NewDefaultDelegate(),0,0)
	finalList.Title = "Files"
	finalList.Styles.Title=lipgloss.NewStyle().
	Foreground(lipgloss.Color("16")).
	Background(lipgloss.Color("205")).
	Padding(0,1)

	return model{
		newfile:     ti,
		isEditing:   false,
		data:        tii,
		dataEditing: false,
		list: finalList,
		showingList: false,
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


func listFile()[]list.Item{
	items:=make([]list.Item,0)

	entries,err:=os.ReadDir(fileHome)
	if err!=nil{
		log.Fatal("Error reading the location ",err)
	}

	for _, entry:=range entries{
		
		if !entry.IsDir(){
			info,err:=entry.Info()
			if err!=nil{
				// log.Fatal("its not a directory ",err)
				continue
			}
			modTime:= info.ModTime().Format("2006-01-02 15:04")
			items = append(items,item{
				title: info.Name(),
				desc: fmt.Sprintf("Modified: %s",modTime),
			})

		}
	}
	return items
}