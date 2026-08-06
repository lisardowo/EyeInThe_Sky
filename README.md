# EyeInThe_Sky

Eye in the sky is a TUI to simple monitor my home servers, it is inspired in SIEM tools but instead of prepare and add logs it just provide useful intel about the state of the machines. Some of the information that this tool is expected to provide is:

    - Logs -- tools is inspired on vim keybindings, so with commands you can change which specific logs you want to look at (network, OS, Application, Firewall etc)
    - Resources Usage -- Current resource usage reported by the machine. Thinking of adding some type of memory that tells you previous reported resource usage so it can help to deduct if any malware is in the machine 
    - Active processes -- Enumerate, classify and highlight the active processes in the server
    - Security of the network -- Still thinking bout this, it was supposed to detect suspicious activity in the servers and default to a more protected, less permissive connection between the TUI and server. 

This tool is thought to be compiled and used in a microcomputer such as a raspberry pico 2w so its default recommended used is usb connection but Im looking forward to implement ssh connection as well. Connection directed to a machine is possible.

This project is just for fun and the sake of learning

## Bugs

Current errors Ive found 
1. ~~ Timestamp boot time evaluating to an error number ~~
2. Extra width in dashboard parameters..?
3. Home view does not consider the height of the screen to construct the render while dash does
 4.~~ i THINK THIS IS MORE OF A DEBUG ERROR BUT TERMINAL DASHBOARD DEFAULTS TO UNSAFE ENVIROMENT EVEN THO THE HOME DASHBOARD WAS MARKED AS SAFE ~~
5. ~~ UNSAFE ENVIROMENT DOES NOT RENDER RESOURCE CONSUMPTION ~~
6. ~~ In the dashboard view, when moving in between focused menus, if you press either l or h more than 3 times the focused highlight starts to behalf weird, sometimes not focusing the correct screen or sometimes jumping one/not rendering properly ~~

** 4 and 5 errors both were debug errors happening because i was passing in the main file an empty --debug only-- home and dashboard model

7. The event stream keeps on growing when appending strings blocking the TUI completely, expected behavior is to leave behind old data keeping just fresh logs, the data is not deleted just displaced, when someone explictly enters via command in the log manager window then they can scroll trough the whole log

8. ~~ Home screen does not change color between secure and unsecure, dashboard works tho. TrustLevel parameter seems to not be connected what could be problematic ~~

9. ~~ Home screen still not changing colors ~~

10. The Time counter just updates when interacted with the TUI which is weird, implement the counter as a goroutine and reduce the amount of seconds displayed 

## Code Reference

This section is a compact map of the most important types, variables, constants, and functions.

## Types

### `Model` - Root UI router and shared terminal state - BubbleTea core dependency
Keeps the active screen, the home state, the dashboard state, and the terminal size.
This is the main object everything else depends on for screen switching and rendering.
Source: [TUI/model.go](TUI/model.go#L18)

```go
type Model struct {
	WhichScreen int
	Home        HomeState
	Dash        DashState
	Width,Height int
}
```

### `HomeState` - Home screen data container
Stores the welcome screen data such as operator, mode, boot time, and uptime.
The home renderer reads this structure directly.
Source: [TUI/home.go](TUI/home.go#L32)

```go
type HomeState struct {
	TrustLevel connection.TrustLevel
	Operator     string
	VLAN         int
	BootAt       time.Time
	Uptime       time.Duration
}
```

### `DashState` - Dashboard layout and telemetry state
Holds dashboard-specific data such as focus, resource usage, and logs.
The dash renderer depends on this state for panel borders and content.
Source: [TUI/dash.go](TUI/dash.go#L32)

```go
type DashState struct {
	TrustLevel     bool
	FocusedPanel FocusPanel
	CPUUsage     float64
	RAMUsage     float64
	LogsBuffer   []string
}
```

### `FocusPanel` - Dashboard focus enum
Defines the allowed focused panels for the dashboard.
This avoids string typos when switching the active panel.
Source: [TUI/dash.go](TUI/dash.go#L10)

```go
type FocusPanel int

const (
	PanelTelemetry FocusPanel = iota
	PanelCommands
	PanelLogs
)
```

## Constants

### `HomeScreen` / `DashScreen` - Screen mode selectors
Enumerator representing which view should be rendered.
The main model reads them in `Update` and `View`.
Source: [TUI/model.go](TUI/model.go#L10)

```go
const (
	HomeScreen int = iota
	DashScreen
)
```

## Functions

### `main` - App entry point
Builds the initial model, starts Bubble Tea, and launches the TUI.
Every screen depends on the state created here.
Source: [main.go](main.go#L23)

```go
initialModel := tui.Model{
	WhichScreen: 0,
	Home:        tui.HomeState{},
	Dash:        tui.DashState{},
	Width:       120,
	Height:      240, // DEBUG START TUI MODEL 
}
```

### `Model.Init` - Requests terminal size on startup
Asks Bubble Tea for the first `WindowSizeMsg`.
This is important because all layout math depends on width and height.
Source: [TUI/model.go](TUI/model.go#L31)

```go
func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}
```

### `Model.Update` - Handles keys and screen changes
Updates terminal size, quits on Ctrl+Q, and switches between home and dash.
This is the only place where user input should mutate the model.
Source: [TUI/model.go](TUI/model.go#L35)

```go
case tea.KeyEnter:
	m.WhichScreen = DashScreen
	m.LastAction = "start handshake"
	return m, nil
```

### `Model.View` - Selects the active renderer
Chooses between the home view and the dashboard view.
This is the routing layer for the whole TUI.
Source: [TUI/model.go](TUI/model.go#L76)

```go
if m.WhichScreen == DashScreen {
	return renderDash(m.Dash, m.Height, m.Width)
}

return renderHome(m.Home, m.Height, m.Width)
```

### `renderHome` - Renders the welcome screen
Builds the startup layout, operator panel, and help text.
The home screen depends on this for its final visual structure.
Source: [TUI/home.go](TUI/home.go#L60)

```go
func renderHome(m HomeState, TerminalHeight int, TerminalWidth int) string {
	operator := m.Operator
	if operator == "" {
		operator = "operator"
	}
}
```

### `renderDash` - Renders the dashboard layout
Builds telemetry, commands, and logs panels from the dashboard state.
The dashboard depends on terminal size and focused panel selection.
Source: [TUI/dash.go](TUI/dash.go#L43)

```go
func renderDash(state DashState, TerminalHeight int, TerminalWidth int) string {
	topHalfHeight := (TerminalHeight / 2) - 2
	panelAwidth := (TerminalWidth / 2) - 2
}
```

### `onFocus` - Small conditional color helper
Returns one color when the condition is true, another when false.
Used by the dashboard to keep border colors readable.
Source: [TUI/dash.go](TUI/dash.go#L179)

```go
func onFocus(cond bool, t, f lipgloss.Color) lipgloss.Color {
	if cond {
		return t
	}
	return f
}
```

### `renderProgressBar` - Builds the CPU and RAM bars
Converts a percentage into a text bar.
The dashboard depends on it to show resource usage.
Source: [TUI/dash.go](TUI/dash.go#L186)

```go
func renderProgressBar(percent float64, width int) string {
	bars := int((percent / 100.0) * float64(width))
	return strings.Repeat("█", bars)
}
```
