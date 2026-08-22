package TUI

import (
	//connection "EyeInThe_Sky/createConnection"

	connection "EyeInThe_Sky/createConnection"

	"EyeInThe_Sky/helpers"
	"EyeInThe_Sky/sysinfo"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	HomeScreen int = iota
	DashScreen
)

type tickMsg time.Time
type FetchMsg struct{
	Entries []helpers.LogEntry
}

type MetricMsg struct{
	CPU float64
	RAM float64
	Sample sysinfo.CPUSample
}

type Model struct {
	//General usage
	prevCPUSample	sysinfo.CPUSample
	TrustLevel connection.TrustLevel
	Uptime			int
	//Logs
	LastKey      string	
	LastAction   string // Last change
	//Screen Information
	Width,Height	int
	WhichScreen int
	Home HomeState
	Dash DashState
}




func (m Model) Init() tea.Cmd {
	 return tea.Batch(tea.WindowSize(),
	  tickCmd(),
	  fetchMetricsCmd(m.prevCPUSample),
	  fetchAllCmd(),
	  collectTickCmd(),
	)
	  
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	
	
	case FetchMsg:
		for _, entry := range msg.Entries{
			if evicted, ok := m.Dash.LogsBuffer.Add(entry); ok{
				helpers.WriteReport(evicted)
			}
		}
		return m, collectTickCmd() //TODO fetch all the entries from processesTickCmd -> ReadProcessLogs -> ReadProcessLog

	case MetricMsg:
		m.Dash.CPUUsage = msg.CPU//TODO Move the cases to a separate function to avoid having long ass code in the update func
		m.Dash.RAMUsage = msg.RAM
		m.prevCPUSample = msg.Sample

		return m, tea.Tick(1 * time.Second, func(_ time.Time) tea.Msg{
			return fetchMetricsCmd(m.prevCPUSample)()
		})

	case tickMsg:
		
		m.Home.Uptime++
		return m, tickCmd()

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	
	case tea.KeyMsg:
		m.LastKey = msg.String()

		if msg.Type == tea.KeyCtrlQ {
			return m, tea.Quit
		}

		if m.Dash.CommandMode{
			switch msg.Type{
				case tea.KeyEnter:
					cmd := m.Dash.CommandInput
					m.Dash.CommandMode = false
					m.Dash.CommandInput = ""

					trimmed := strings.TrimPrefix(cmd, ":")
					if trimmed == "all" || trimmed == "" {
						m.Dash.ActiveFilter = nil
					} else if cat, ok := helpers.SelectCategory(cmd); ok{
						m.Dash.ActiveFilter = &cat
					}

					return m, nil

				case tea.KeyEsc:
					m.Dash.CommandMode = false
					m.Dash.CommandInput = ""
					return m,nil
				case tea.KeyBackspace:
					if len(m.Dash.CommandInput) > 1 {
						m.Dash.CommandInput = m.Dash.CommandInput[:len(m.Dash.CommandInput) - 1]
					}
					return m,nil
				case tea.KeyRunes:
					m.Dash.CommandInput += msg.String()
					return m, nil
				default:
					return m, nil //TODO raise a warning of not recognized commands
			}
		}

		switch msg.Type {
		case tea.KeyEnter:
			m.WhichScreen = DashScreen
			m.LastAction = "start handshake"
			return m, nil
		case tea.KeyEsc:
			m.WhichScreen = HomeScreen
			m.LastAction = "return home"
			return m, nil
		
		}
		//TODO move this to a wrapping fuction to keep update clean
		if(m.WhichScreen == DashScreen){
			switch msg.String() {
				
		
				case ":":
					
					m.Dash.CommandMode = true
					m.Dash.CommandInput = ":"
					
					return m, nil

				case "l":

					if(m.Dash.FocusedPanel < 2){
						m.Dash.FocusedPanel += 1
					} else
					{

						m.Dash.FocusedPanel = 0
					}
				
				case "h":

					if(m.Dash.FocusedPanel == 0){
						m.Dash.FocusedPanel = 2
					} else
					{

						m.Dash.FocusedPanel -= 1 
					}
				default: // TODO implement an "warning" notification letter not recognized, this goes over the layout in the corner of commands
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	//m.Home.Uptime = time.Since(m.Home.BootAt)
	
	if m.WhichScreen == DashScreen { 
		return renderDash(m.Dash, m.Height, m.Width, m.TrustLevel)
	}

	return renderHome(m.Home, m.Height, m.Width, m.TrustLevel, m.Home.Uptime)  
}

func fetchMetricsCmd(prevSample sysinfo.CPUSample) tea.Cmd{
	return func() tea.Msg{
		currSample, _ := sysinfo.GetCPUSample()
		cpu, _ := sysinfo.CalculateCPUusage(prevSample, currSample)
		ram, _ := sysinfo.GetRamUsage() //TODO add if cpu/ram err != nil check

		return MetricMsg{
			CPU: cpu,
			RAM: ram,
			Sample: currSample,
		}
	}
}

func fetchAllCmd() tea.Cmd{
	
	return func() tea.Msg {
	var all []helpers.LogEntry

		if e, err := helpers.ReadProcessLogs(); err == nil {
			all = append(all, e...)
		}
		if e, err := helpers.ReadTCPLogs(); err == nil {
			all = append(all, e...)
		}
		if e, err := helpers.ReadUDPLogs(); err == nil {
			all = append(all, e...)
		}
		if e, err := helpers.ReadModulesLogs(); err == nil {
			all = append(all, e...)
		}
		if e, err := helpers.ReadDiskstatsLogs(); err == nil {
			all = append(all, e...)
		}

		return FetchMsg{Entries: all}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(1 * time.Second, func(t time.Time) tea.Msg{
		return tickMsg(t)
	})
}

func collectTickCmd() tea.Cmd {
	return tea.Tick(8*time.Second, func(t time.Time) tea.Msg {
		return fetchAllCmd()()
	})
}