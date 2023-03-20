package main

//
import (
	"fmt"
	"os"

	"github.com/404name/bilibili-terminal/view"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// m := home.NewTuiModel()
	m := view.NewApp()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(0)
	}
}
