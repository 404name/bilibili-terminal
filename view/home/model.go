package home

import (
	"fmt"
	"io"

	"github.com/404name/bilibili-terminal/components"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/image"
	"github.com/knipferrc/teacup/statusbar"
)

// 页面状态
type ViewState uint
type errMsg error
type LoadImgMsg string
type DataReadyMsg []*RcmdVideo

const (
	mainView ViewState = iota
	dialogView
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = lipgloss.NewStyle().PaddingLeft(4).PaddingBottom(1)
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
	// textStyle         = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func (i RcmdVideo) Title() string { return i.Headline }
func (i RcmdVideo) Description() string {
	return fmt.Sprintf("🧑%20s 👀%d 👍%d 💬%d", i.Owner.Name, i.Stat.View, i.Stat.Like, i.Stat.Danmaku)
}
func (i *RcmdVideo) FilterValue() string {
	return i.Headline
}

var HomeView *Tui

type Tui struct {
	state     ViewState
	dialog    *components.Dialog
	statusbar statusbar.Bubble
	image     image.Bubble
	spinner   *spinner.Model
	list      *list.Model
	selected  *RcmdVideo
	loading   bool
	err       error
}

// 用bubbletea写的一个终端界面，用来选择搜索结果，然后打开链接
func NewTuiModel() Tui {
	if HomeView != nil {
		return *HomeView
	}
	// 用自己实现的rowDelegate来渲染列表
	// c := list.New(nil, rowDelegate{}, 0, 0)
	// 用默认的rowDelegate来渲染列表
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "🔥百度热搜Trending（回车访问）"
	l.SetShowStatusBar(true)
	l.SetShowTitle(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	statusbarModel := statusbar.New(
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Dark: "#cc241d", Light: "#cc241d"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Dark: "#3c3836", Light: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Dark: "#A550DF", Light: "#A550DF"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Dark: "#6124DF", Light: "#6124DF"},
		},
	)
	statusbarModel.SetContent(
		"上下左右移动光标",
		"esc/q退出",
		"空格选取",
		"回车确认",
	)
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	// TODO Implement a progressive loading list
	HomeView = &Tui{
		state:     mainView,
		dialog:    components.NewDialog("你要访问吗", "是的", "不要"),
		image:     image.New(false, true, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}),
		statusbar: statusbarModel,
		spinner:   &s,
		list:      &l,
		loading:   true,
	}
	return *HomeView
}

func (m Tui) Init() tea.Cmd {

	loading := func() tea.Msg {
		return spinner.TickMsg{}
	}
	init := func() tea.Msg {
		if data, err := GetTrendingList(); err != nil {
			return err
		} else {
			return DataReadyMsg(data)
		}
	}
	loadImg := func() tea.Msg {
		return LoadImgMsg("logo1.png")
	}
	// 注册加载中的动画和初始化事件
	return tea.Batch(loading, init, loadImg)
}

// 用于渲染列表单个内容的样式
type rowDelegate struct{}

func (d rowDelegate) Height() int {
	return 1
}

func (d rowDelegate) Spacing() int {
	return 0
}

func (d rowDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {

	return nil
}

// 渲染列表
func (d rowDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*RcmdVideo)
	if !ok {
		return
	}
	q := (*RcmdVideo)(i)
	str := fmt.Sprintf("%s [👀%d 👍%d 💬%d] ", q.Headline, q.Stat.View, q.Stat.Like, q.Stat.Danmaku)
	if index == m.Index() {
		str = selectedItemStyle.Render("> " + str)
	} else {
		str = itemStyle.Render(str)
	}
	_, _ = fmt.Fprint(w, str)
}
