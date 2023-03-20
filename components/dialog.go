package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	minHeight              = 1
	minWidth               = 2
	defaultHeight          = 7
	defaultWidth           = 100
	questionWidth          = 50
	defaultWhitespaceChars = "诶嘿"
)

var (
	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			MarginLeft(1).
			MarginRight(1).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				Underline(true)
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)

type Dialog struct {
	Active              bool
	ChooseOK            bool
	question            string
	okButtonDesc        string
	cancelButtonDesc    string
	withWhitespaceChars string
	questionWidth       int
	width               int
	height              int
	okCmd               tea.Cmd
}

func NewDialog(question, okButtonDesc, cancelButtonDesc string) *Dialog {
	return &Dialog{
		Active:              false,
		ChooseOK:            true,
		question:            question,
		okButtonDesc:        okButtonDesc,
		cancelButtonDesc:    cancelButtonDesc,
		withWhitespaceChars: defaultWhitespaceChars,
		width:               defaultWidth,
		height:              defaultHeight,
		questionWidth:       questionWidth,
	}
}
func (d *Dialog) SetActive(active bool) {
	d.Active = active
	d.ChooseOK = true
}
func (d *Dialog) SetSize(width, height int) {
	if width < minWidth {
		width = minWidth
	}
	if height < minHeight {
		height = minHeight
	}
	d.width = width
	d.height = height
	d.questionWidth = width / 2
}

func (d *Dialog) SetQuestionAndCmd(question string, okCmd tea.Cmd) {
	d.question = question
	d.okCmd = okCmd
}

func (d Dialog) Init() tea.Cmd {
	return nil
}

func (d *Dialog) Update(msg tea.Msg) (*Dialog, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "left":
			d.ChooseOK = true
		case "down", "right":
			d.ChooseOK = false
		case "enter", " ":
			d.Active = false
			if d.ChooseOK {
				cmd = d.okCmd
			}
		}
	}
	return d, cmd
}

func (d Dialog) View() string {

	var okButton, cancelButton string
	var ui string
	if d.ChooseOK {
		okButton = activeButtonStyle.Render(d.okButtonDesc)
		cancelButton = buttonStyle.Render(d.cancelButtonDesc)
	} else {
		okButton = buttonStyle.Render(d.okButtonDesc)
		cancelButton = activeButtonStyle.Render(d.cancelButtonDesc)
	}

	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	question := lipgloss.NewStyle().Width(d.questionWidth).Align(lipgloss.Center).Render(d.question)

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)

	if d.Active {
		ui = lipgloss.JoinVertical(lipgloss.Center, question, buttons)
	}
	dialog := lipgloss.Place(d.width, d.height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars(d.withWhitespaceChars),
		lipgloss.WithWhitespaceForeground(subtle),
	)
	return dialog + "\n"
}
