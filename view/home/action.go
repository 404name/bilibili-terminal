package home

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/icons"
)

func (m Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.statusbar.SetSize(msg.Width)
		m.list.SetSize(msg.Width/2, msg.Height-m.statusbar.Height)
		resizeImgCmd := m.image.SetSize(msg.Width/2+10, msg.Width/20*3)
		m.dialog.SetSize(msg.Width/2, msg.Height-m.statusbar.Height-m.image.Viewport.Height)

		return m, tea.Batch(resizeImgCmd, tea.ClearScreen)
	case DataReadyMsg:
		data := DataReadyMsg(msg)
		list := m.list.Items()

		if len(list) > 0 {
			list = list[:len(list)-1]
		}

		for _, q := range data {
			list = append(list, q)
		}
		loadMore := RcmdVideo{
			Headline: "【加载更多】",
			URI:      "https://www.bilibili.com",
			Stat: struct {
				View    int `json:"view"`    // 观看数量
				Like    int `json:"like"`    // 喜欢数
				Danmaku int `json:"danmaku"` // 弹幕数量
			}{
				View:    0,
				Like:    0,
				Danmaku: 0,
			},
		}
		list = append(list, &loadMore)

		m.list.SetItems(list)
		m.loading = false
		return m, nil
	case spinner.TickMsg:
		// 用来显示加载中的动画
		if m.loading == false {
			return m, nil
		}
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = &spinner
		return m, cmd
	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	switch m.state {
	// 交给不同的view去处理事件
	case mainView:
		m, cmd = m.updateMainView(msg)
	case dialogView:
		m.dialog, cmd = m.dialog.Update(msg)
	}

	if m.dialog.Active {
		m.state = dialogView
	} else {
		m.state = mainView
	}

	return m, cmd
}

func (m Tui) updateMainView(msg tea.Msg) (Tui, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "enter":
			// 最后一个是加载更多
			if m.list.Index() == len(m.list.Items())-1 {
				m.dialog.SetActive(true)
				m.dialog.SetQuestionAndCmd("确定要加载更多推荐吗？", LoadMore)
				return m, nil
			}

			// 如果是选中的就问是否要打开
			if m.list.Index() == m.lastSelectIndex {
				m.dialog.SetActive(true)
				m.dialog.SetQuestionAndCmd("你要访问【"+m.selected.Headline+"】吗？", open(m.selected.URI))
				return m, nil
			}
			// 第一次选就展示图片
			if m.list.SelectedItem() != nil {
				m.lastSelectIndex = m.list.Index()
				m.selected = (*RcmdVideo)(m.list.SelectedItem().(*RcmdVideo))
			}
			return m, m.getRemoteImg(m.list.SelectedItem().(*RcmdVideo).Pic)
		case "tab":
			m.list.SetDelegate(rowDelegate{})
			m.list.SetSize(m.list.Width(), m.list.Height())
		}
	case LoadImgMsg:
		cmd := m.image.SetFileName(string(msg))
		return m, cmd
	}

	list, _ := m.list.Update(msg)
	m.list = &list
	m.image, _ = m.image.Update(msg)

	if m.loading {
		return m, nil
	}

	logoText := fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "哔哩哔哩")
	m.statusbar.SetContent(
		"UP:"+m.list.SelectedItem().(*RcmdVideo).Owner.Name,
		m.list.SelectedItem().(*RcmdVideo).Headline,
		fmt.Sprintf("%d/%d", m.list.Index()+1, len(m.list.Items())),
		logoText,
	)
	return m, nil
}

func (m Tui) View() string {

	if m.err != nil {
		return m.err.Error()
	}
	if m.loading {
		return fmt.Sprintf("\n\n   %s 正在加载数据\n %s \n", m.spinner.View(), m.statusbar.View())
	} else {
		return "\n" +
			lipgloss.JoinVertical(lipgloss.Top,
				lipgloss.JoinHorizontal(lipgloss.Center,
					m.list.View(),
					lipgloss.JoinVertical(lipgloss.Top,
						m.image.View(),
						m.dialog.View(),
					)),
				m.statusbar.View(),
			)
	}
}

func (m *Tui) Selected() *RcmdVideo {
	return m.selected
}
