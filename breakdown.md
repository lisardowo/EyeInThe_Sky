# Here's a small breakdown so I can organize my ideas while developing this
# I asked gemini to generate this with my ideas
---

## ⬢ 1. Conceptual Breakdown: Data Flow and Analysis

To operate safely and efficiently in resource-limited environments (such as a Raspberry Pi Zero), the tool uses a **"Process-as-you-go"** paradigm (process and discard/persist in bursts) rather than accumulating infinite data in RAM. The data flow is divided into three critical layers:

### A. Transport and Basic Telemetry Layer
*   **Link Origin (`Transport`):** Dynamic identification of the physical input medium (`PHYSICAL_USB`) or emergency remote fallback (`NETWORK_SSH`).
*   **Latency Metrics (`time.Duration`):** Precise measurement of packet round-trip time to adjust the UI visual refresh rate for high-latency links (such as satellite bursts).
*   **Heartbeat (Life Pulse):** Active background timer that triggers a "Link Lost" state if the node stops transmitting for more than $X$ seconds.

### B. Security Layer (Sanitization and Isolation)
*   **ANSI Escape Cleaning:** When connecting to a vulnerable lab, the system filters and sanitizes all incoming text strings to prevent terminal injection attacks (vulnerabilities where malicious data attempts to corrupt rendering or execute local commands).
*   **Automatic Trust Classification:** Deduction of `TrustLevel` (`Secure` or `Unsecure`) by analyzing connection metadata (such as IP addressing or serial ports) to dynamically restrict UI capabilities (for example, disabling write commands in untrusted zones).

### C. Specialized Parsing Layer (Payloads)
*   **Infrastructure Mode (m4c-ne):** Reading and interpreting plain text flows (`journalctl`) or JSON structures to extract telemetry from Proxmox servers (CPU usage, RAM, firewall alerts).
*   **Satellite Mode (Artemis):** Reception of raw byte arrays (`[]byte`), checksum validation (low-level CRC), and extraction of readable physical variables (coordinates, thermal telemetry).

---

## ⬢ 2. File Architecture (Library Structure)

In Go, modularity is organized through directories that act as independent, isolated packages. This is the structure that will keep your `main.go` clean and free from name collisions:

```text
EyeInThe_Sky/
├── go.mod                 # Module definition and external dependency management
├── main.go                # Single entry point (initializes and launches the TUI)
│
├── createConnection/      # Network Backend and Hardware Contracts
│   └── connection.go      # Interfaces (SourceNode), Handshake, Mocks, and BootNode
│
├── parser/                # Data Processing and Security Layer
│   └── sanitizer.go       # ANSI string sanitization and payload-specific parsing
│
└── tui/                   # Visual Nervous System (Bubble Tea + Lip Gloss)
    ├── model.go           # Global 'model' struct and UI Finite State Machine
    ├── update.go          # Keyboard event handler (Vim-style keybindings)
    ├── view_home.go       # Home Screen Rendering (Welcome, Flags, Keybinds)
    └── view_dash.go       # Operational Panel Rendering (Log Viewport, Commands)
```

---

## ⬢ 3. Implementation Order and Strategy

### **Phase 1: Foundation & Transport Layer** (Start here)
This layer is the lowest and most independent. Build it first without any UI dependencies.

**Files to create:**
- `createConnection/connection.go` — Define core interfaces and data structures:
  ```go
  type Transport string  // PHYSICAL_USB, NETWORK_SSH
  
  type TrustLevel string // Secure, Unsecure
  
  type SourceNode interface {
    Connect() error
    Disconnect() error
    Read() ([]byte, error)
    Write(data []byte) error
    Latency() time.Duration
  }
  
  type LinkState struct {
    Transport      Transport
    TrustLevel     TrustLevel
    LastHeartbeat  time.Time
    Latency        time.Duration
    Connected      bool
  }
  
  type BootNode struct {
    ID            string
    LinkState     LinkState
    HeartbeatTick <-chan time.Time
  }
  ```

**Mock implementations:**
- Implement `SourceNode` with mock USB and SSH connectors for testing without real hardware.

---

### **Phase 2: Security & Parser Layer** (Build after Phase 1)
This layer processes the raw data from Phase 1 and makes it safe for display.

**Files to create:**
- `parser/sanitizer.go` — Define sanitization and parsing functions:
  ```go
  type PayloadMode string  // Infrastructure, Satellite
  
  type SanitizedPayload struct {
    Raw        []byte
    Cleaned    string
    Mode       PayloadMode
    Timestamp  time.Time
  }
  
  func SanitizeANSI(raw string) string {
    // Strip ANSI escape sequences
  }
  
  func ParsePayload(data []byte, mode PayloadMode) (SanitizedPayload, error) {
    // Route to appropriate parser (JSON, plaintext, binary)
  }
  
  func ValidateChecksum(data []byte) bool {
    // Validate CRC for satellite payloads
  }
  ```

---

### **Phase 3: TUI Layer** (Build after Phases 1 & 2)
This layer ties everything together and provides the user interface.

**Files to create:**
- `tui/model.go` — Define the central model:
  ```go
  type Model struct {
    BootNode      *createConnection.BootNode
    PayloadBuffer []parser.SanitizedPayload
    UIState       UIState
    CurrentView   ViewMode
  }
  
  type UIState string  // Welcome, Connected, Dashboard, Error
  type ViewMode string // Home, Dashboard, Logs
  ```

- `tui/update.go` — Handle messages and state transitions:
  ```go
  func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
      return handleKeyPress(m, msg)
    case LinkHeartbeatMsg:
      return handleHeartbeat(m)
    case PayloadReceivedMsg:
      return handlePayload(m, msg)
    }
    return m, nil
  }
  ```

- `tui/view_home.go` — Render the welcome screen:
  - Reuse or refactor existing `home.go` logic
  - Display operator info, VLAN, mode, uptime

- `tui/view_dash.go` — Render the operational dashboard:
  - Display real-time logs with scrolling
  - Show link metrics (latency, heartbeat status)
  - Command prompt

---

### **Phase 4: Integration in main.go** (Final step)
Wire everything together:

```go
package main

import (
    "EyeInThe_Sky/createConnection"
    "EyeInThe_Sky/parser"
    "EyeInThe_Sky/tui"
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    // 1. Initialize a BootNode with mocked or real transport
    bootNode := createConnection.NewBootNode("pi-zero-01", createConnection.PHYSICAL_USB)
    
    // 2. Create the TUI model
    model := tui.Model{
        BootNode: bootNode,
        UIState:  tui.Welcome,
    }
    
    // 3. Launch Bubble Tea
    p := tea.NewProgram(model)
    if _, err := p.Run(); err != nil {
        panic(err)
    }
}
```

---

## ⬢ 4. Key Structs and Variables Summary

| Layer | Struct/Type | Purpose |
|-------|-------------|---------|
| **Transport** | `SourceNode` | Interface for USB/SSH connections |
| **Transport** | `LinkState` | Current link metrics and status |
| **Security** | `SanitizedPayload` | Cleaned and validated data |
| **TUI** | `Model` | Central state machine for the UI |
| **TUI** | `UIState` | Enum: Welcome, Connected, Dashboard, Error |
| **TUI** | `ViewMode` | Enum: Home, Dashboard, Logs |

---

## ⬢ 5. Development Tips

1. **Test each layer independently** before moving to the next (use `go test`).
2. **Mock external dependencies** (real hardware) in Phase 1 and 2.
3. **Use channels** for heartbeat and payload reception (goroutines in background).
4. **Keep `main.go` small** — all logic lives in the packages.
5. **Refactor existing `TUI` and `home` packages** to align with the new architecture as you build Phases 1–3.
