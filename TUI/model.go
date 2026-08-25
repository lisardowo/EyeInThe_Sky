package TUI

import (
	//connection "EyeInThe_Sky/createConnection"

	connection "EyeInThe_Sky/createConnection"
	createconnection "EyeInThe_Sky/createConnection"

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

type USBFrameMsg struct {
	Frame createconnection.RemoteFrame
	Err error
}

type Model struct {
	//General usage
	usbConn		*createconnection.USBConnection
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



// =================================
// Model render Logic
// =================================

func (m Model) Init() tea.Cmd {
	 return tea.Batch(tea.WindowSize(),
	  tickCmd(),
	 /* fetchMetricsCmd(m.prevCPUSample),
	  fetchAllCmd(),
	  collectTickCmd(), STOP READING FROM HOST*/
	)
	  
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	
	case USBFrameMsg:
		if msg.Err == nil {
			m.Dash.CPUUsage = msg.Frame.CPU
			m.Dash.RAMUsage = msg.Frame.RAM
			if msg.Frame.Log != "" {
				entry := helpers.LogEntry{
					Timestamp: helpers.SKYnow(),
					Category: helpers.LogCategory(msg.Frame.Category),
					Message: msg.Frame.Log,
				}
				m.Dash.LogsBuffer.Add(entry)
			}

			return m, waitForUSBData(m.usbConn)

		}
		return m, nil
	
	case FetchMsg:
		for _, entry := range msg.Entries{
			if evicted, ok := m.Dash.LogsBuffer.Add(entry); ok{
				helpers.WriteReport(evicted)
			}
		}
		return m, collectTickCmd() 

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
						m.Dash.ActiveFilter = &cat //filters current buffer content therefor not reaching the commands
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

			m.usbConn = createconnection.NewUSBConnection()
			if err := m.usbConn.Connect(); err != nil {
				return m, nil //TODO raise an alert for failed connection
			}
			if err := m.usbConn.PerformHandshake(); err != nil{
				m.usbConn.Close()
				return m,nil
			}

			return m, waitForUSBData(m.usbConn)
			
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


//  =================================
// Fetch Entries logic
// ================================= 

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


//  =================================
//   Utility functions
// ================================= 


func GetSize(TerminalHeight int, TerminalWidth int,) (leftWidth int, rightWidth int, stacked bool){

	availableWidth := TerminalWidth
	if availableWidth <= 0 {
		availableWidth = 120
	}
// TODO calculations below can be moved to a utility getSize function
	minPaneWidth := 28
	gap := 2
	stacked = availableWidth < 90

	leftWidth = availableWidth - 4
	rightWidth = leftWidth - 4
	if !stacked {
		leftWidth = int(float64(availableWidth) * 0.58)
		rightWidth = availableWidth - leftWidth - gap - 6
		if leftWidth < 34 {
			leftWidth = 34
		}
		if rightWidth < 24 {
			rightWidth = 24 - 6
		}
		if leftWidth+rightWidth+gap > availableWidth {
			rightWidth = availableWidth - leftWidth - gap - 6
		}
		if rightWidth < minPaneWidth {
			rightWidth = minPaneWidth
		}
	} else {
		if leftWidth < minPaneWidth {
			leftWidth = minPaneWidth
		}
		rightWidth = leftWidth - 6
	}

	return leftWidth, rightWidth, stacked
}

func waitForUSBData(connection *connection.USBConnection) tea.Cmd {
	return func() tea.Msg {
		frame, err := connection.ReadNextFrame()
		return USBFrameMsg{Frame: frame, Err: err}
	}
}

