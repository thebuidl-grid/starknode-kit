package monitoring

import (
	"time"

	"github.com/rivo/tview"
)

type MonitorApp struct {
	App        *tview.Application
	Grid       *tview.Grid
	SystemBox  *tview.TextView
	ClientsBox *tview.TextView
	LogsBox    *tview.TextView
	PeersBox   *tview.TextView
	ChainBox   *tview.TextView
	StatusBox  *tview.TextView

	// Update channels
	SystemChan  chan string
	ClientsChan chan string
	LogsChan    chan string
	PeersChan   chan string
	ChainChan   chan string

	// Control
	StopChan   chan bool
	UpdateRate time.Duration
}
