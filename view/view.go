package view

import (
	"github.com/404name/bilibili-terminal/view/home"
	tea "github.com/charmbracelet/bubbletea"
)

// 页面状态
type ViewState uint
type errMsg error

const (
	startView ViewState = iota
	homeView
)

// TODO 后面做一个router存储，可以返回上一级，也有保活机制
// TODO 做一个全局loading
type App struct {
	state ViewState

	Tui   home.Tui
	Start interface{}

	err errMsg
}

func NewApp() *App {
	return &App{
		state: homeView,
		Tui:   home.NewTuiModel(),
	}
}

func (app *App) Init() tea.Cmd {
	return app.Tui.Init()
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// 主app只处理优先级最高的一些指令,切换页面，以及报错页面处理
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return app, tea.Quit
		}
	case ViewState:
		app.state = ViewState(msg)
		return app, nil
	case errMsg:
		// TODO 做一个报错页面
		app.err = msg
		return app, nil
	}
	var cmd tea.Cmd
	var model tea.Model
	// 其它按键事件、resize事件等等，交给当前view去处理
	switch app.state {
	// 交给不同的view去处理事件
	case homeView:
		model, cmd = app.Tui.Update(msg)
		app.Tui = model.(home.Tui)
	case startView:
		model, cmd = app.Tui.Update(msg)
		app.Tui = model.(home.Tui)
	}

	return app, cmd
}

func (app *App) View() string {
	switch app.state {
	// 交给不同的view去处理事件
	case homeView:
		return app.Tui.View()
	case startView:
		return ""
	}
	return ""
}
