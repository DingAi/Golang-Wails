package service

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"go.bug.st/serial"
)

// SerialPortInfo represents a single serial port's basic info
type SerialPortInfo struct {
	Name string `json:"name"`
}

// SerialStatus represents the current status displayed in the Dashboard status bar
type SerialStatus struct {
	IsOpen   bool   `json:"is_open"`
	PortName string `json:"port_name"`
	BaudRate int    `json:"baud_rate"`
	RxBytes  int64  `json:"rx_bytes"`
	TxBytes  int64  `json:"tx_bytes"`
}

// SerialConfig represents the serial port configuration parameters
type SerialConfig struct {
	PortName string `json:"port_name"`
	BaudRate int    `json:"baud_rate"`
	DataBits int    `json:"data_bits"`
	Parity   string `json:"parity"`
	StopBits int    `json:"stop_bits"`
}

// SerialService handles all serial port operations for the application
type SerialService struct {
	mu      sync.Mutex
	port    serial.Port
	mode    *serial.Mode
	config  *SerialConfig
	isOpen  bool
	rxBytes int64
	txBytes int64

	// Receive buffer
	recvBuffer strings.Builder

	// Auto-send control
	autoSendCancel chan struct{}
	autoSendMu     sync.Mutex

	// Receive event buffer (for auto frame splitting)
	recvLastTime time.Time
	splitMs      int
	pendingData  []byte
	splitTimer   *time.Timer
	splitMu      sync.Mutex
}

// NewSerialService creates a new SerialService instance
func NewSerialService() *SerialService {
	return &SerialService{
		config: &SerialConfig{
			PortName: "COM5",
			BaudRate: 115200,
			DataBits: 8,
			Parity:   "无",
			StopBits: 1,
		},
		mode: &serial.Mode{
			BaudRate: 115200,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		},
	}
}

// GetSerialPorts lists all available serial ports on the system
func (s *SerialService) GetSerialPorts() ([]SerialPortInfo, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("failed to list serial ports: %w", err)
	}

	result := make([]SerialPortInfo, len(ports))
	for i, p := range ports {
		result[i] = SerialPortInfo{Name: p}
	}
	return result, nil
}

// GetStatus returns the current serial port status for the Dashboard status bar
func (s *SerialService) GetStatus() SerialStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	portName := ""
	baudRate := 115200
	if s.config != nil {
		portName = s.config.PortName
		baudRate = s.config.BaudRate
	}

	return SerialStatus{
		IsOpen:   s.isOpen,
		PortName: portName,
		BaudRate: baudRate,
		RxBytes:  atomic.LoadInt64(&s.rxBytes),
		TxBytes:  atomic.LoadInt64(&s.txBytes),
	}
}

// OpenSerialPort opens the serial port with the specified configuration
func (s *SerialService) OpenSerialPort(config SerialConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isOpen {
		return fmt.Errorf("serial port is already open")
	}

	// Parse parity
	parity, err := parseParity(config.Parity)
	if err != nil {
		return err
	}

	// Parse stop bits
	stopBits, err := parseStopBits(config.StopBits)
	if err != nil {
		return err
	}

	mode := &serial.Mode{
		BaudRate: config.BaudRate,
		DataBits: config.DataBits,
		Parity:   parity,
		StopBits: stopBits,
	}

	port, err := serial.Open(config.PortName, mode)
	if err != nil {
		return fmt.Errorf("failed to open port %s: %w", config.PortName, err)
	}

	// Set read timeout to prevent blocking forever
	if err := port.SetReadTimeout(100 * time.Millisecond); err != nil {
		port.Close()
		return fmt.Errorf("failed to set read timeout: %w", err)
	}

	s.port = port
	s.mode = mode
	s.config = &config
	s.isOpen = true
	s.recvBuffer.Reset()
	atomic.StoreInt64(&s.rxBytes, 0)
	atomic.StoreInt64(&s.txBytes, 0)

	// Start the read goroutine
	go s.readLoop()

	// Emit status update event
	s.emitStatusEvent()

	log.Printf("Serial port %s opened (baud=%d, data=%d, parity=%s, stop=%d)",
		config.PortName, config.BaudRate, config.DataBits, config.Parity, config.StopBits)

	return nil
}

// CloseSerialPort closes the currently open serial port
func (s *SerialService) CloseSerialPort() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isOpen {
		return nil // already closed, no error
	}

	// Stop auto-send if running
	s.stopAutoSendLocked()

	// Clear split timer
	s.splitMu.Lock()
	if s.splitTimer != nil {
		s.splitTimer.Stop()
		s.splitTimer = nil
	}
	// Flush any pending data
	if len(s.pendingData) > 0 {
		s.flushPendingData()
	}
	s.splitMu.Unlock()

	err := s.port.Close()
	s.port = nil
	s.isOpen = false

	// Emit status update event
	s.emitStatusEvent()

	if err != nil {
		log.Printf("Error closing serial port: %v", err)
		return fmt.Errorf("failed to close port: %w", err)
	}

	log.Println("Serial port closed")
	return nil
}

// SendData sends data through the serial port
// format can be: "hex", "ascii", "bin", "dec", "bcd"
func (s *SerialService) SendData(data string, format string) (int, error) {
	s.mu.Lock()
	if !s.isOpen || s.port == nil {
		s.mu.Unlock()
		return 0, fmt.Errorf("serial port is not open")
	}
	port := s.port
	s.mu.Unlock()

	// Convert data based on format
	rawBytes, err := convertSendData(data, format)
	if err != nil {
		return 0, fmt.Errorf("data conversion error: %w", err)
	}

	n, err := port.Write(rawBytes)
	if err != nil {
		return 0, fmt.Errorf("write failed: %w", err)
	}

	atomic.AddInt64(&s.txBytes, int64(n))
	s.emitStatusEvent()

	log.Printf("Sent %d bytes: %s", n, hex.EncodeToString(rawBytes))
	return n, nil
}

// StartAutoSend starts periodic auto-send with the given data, format, and interval (in seconds)
func (s *SerialService) StartAutoSend(data string, format string, intervalSec float64) error {
	s.autoSendMu.Lock()
	defer s.autoSendMu.Unlock()

	// Cancel existing auto-send if running
	if s.autoSendCancel != nil {
		close(s.autoSendCancel)
	}

	s.mu.Lock()
	if !s.isOpen {
		s.mu.Unlock()
		return fmt.Errorf("serial port is not open")
	}
	s.mu.Unlock()

	cancel := make(chan struct{})
	s.autoSendCancel = cancel

	go func() {
		ticker := time.NewTicker(time.Duration(intervalSec * float64(time.Second)))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_, err := s.SendData(data, format)
				if err != nil {
					log.Printf("Auto-send error: %v", err)
				}
			case <-cancel:
				return
			}
		}
	}()

	log.Printf("Auto-send started (interval=%.1fs)", intervalSec)
	return nil
}

// StopAutoSend stops the periodic auto-send
func (s *SerialService) StopAutoSend() {
	s.autoSendMu.Lock()
	defer s.autoSendMu.Unlock()
	s.stopAutoSendLocked()
}

// ClearReceiveBuffer clears the receive data buffer
func (s *SerialService) ClearReceiveBuffer() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recvBuffer.Reset()
}

// SetAutoSplit configures auto frame splitting parameters
// splitMs: the idle time in milliseconds to consider a frame complete
func (s *SerialService) SetAutoSplit(enabled bool, ms int) {
	s.splitMu.Lock()
	defer s.splitMu.Unlock()

	s.splitMs = ms
	if !enabled {
		if s.splitTimer != nil {
			s.splitTimer.Stop()
			s.splitTimer = nil
		}
		// Flush remaining data as a single frame
		if len(s.pendingData) > 0 {
			s.flushPendingData()
		}
	}
}

// ---------- internal methods ----------

// readLoop continuously reads data from the serial port
func (s *SerialService) readLoop() {
	buf := make([]byte, 4096)

	for {
		s.mu.Lock()
		port := s.port
		isOpen := s.isOpen
		s.mu.Unlock()

		if !isOpen || port == nil {
			return
		}

		n, err := port.Read(buf)
		if err != nil {
			// Timeout is expected, just continue
			if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "Timeout") {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			// Check if port is still open (might have been closed by user)
			s.mu.Lock()
			stillOpen := s.isOpen
			s.mu.Unlock()

			if !stillOpen {
				return
			}

			log.Printf("Serial read error: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if n == 0 {
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		atomic.AddInt64(&s.rxBytes, int64(n))

		// Append to receive buffer
		s.mu.Lock()
		s.recvBuffer.Write(data)
		hexStr := hex.EncodeToString(data)
		s.mu.Unlock()

		// Handle auto frame splitting
		s.splitMu.Lock()
		if s.splitMs > 0 {
			s.pendingData = append(s.pendingData, data...)
			s.recvLastTime = time.Now()

			if s.splitTimer != nil {
				s.splitTimer.Stop()
			}

			s.splitTimer = time.AfterFunc(time.Duration(s.splitMs)*time.Millisecond, func() {
				s.splitMu.Lock()
				s.flushPendingData()
				s.splitMu.Unlock()
			})
		} else {
			// Emit raw received data directly
			s.emitDataEvent(hexStr, data)
		}
		s.splitMu.Unlock()

		// Update status
		s.emitStatusEvent()
	}
}

// flushPendingData emits pending split-frame data as a single event
func (s *SerialService) flushPendingData() {
	if len(s.pendingData) == 0 {
		return
	}

	data := make([]byte, len(s.pendingData))
	copy(data, s.pendingData)
	s.pendingData = s.pendingData[:0]
	s.splitTimer = nil

	hexStr := hex.EncodeToString(data)
	s.emitDataEvent(hexStr, data)
}

// emitStatusEvent emits the current serial status to the frontend
func (s *SerialService) emitStatusEvent() {
	app := application.Get()
	if app == nil {
		return
	}
	status := s.GetStatus()
	app.Event.Emit("serial:status", status)
}

// emitDataEvent emits received data to the frontend
func (s *SerialService) emitDataEvent(hexStr string, raw []byte) {
	app := application.Get()
	if app == nil {
		return
	}

	// Also provide ASCII representation
	asciiStr := strings.Map(func(r rune) rune {
		if r >= 32 && r <= 126 {
			return r
		}
		return '.'
	}, string(raw))

	app.Event.Emit("serial:data", map[string]interface{}{
		"hex":   hexStr,
		"ascii": asciiStr,
		"raw":   raw,
		"len":   len(raw),
	})
}

// stopAutoSendLocked stops auto-send (must hold autoSendMu)
func (s *SerialService) stopAutoSendLocked() {
	if s.autoSendCancel != nil {
		close(s.autoSendCancel)
		s.autoSendCancel = nil
	}
}

// ---------- helper functions ----------

func parseParity(parity string) (serial.Parity, error) {
	switch strings.ToUpper(parity) {
	case "无", "NONE", "N":
		return serial.NoParity, nil
	case "奇", "ODD", "O":
		return serial.OddParity, nil
	case "偶", "EVEN", "E":
		return serial.EvenParity, nil
	case "MARK", "M":
		return serial.MarkParity, nil
	case "SPACE", "S":
		return serial.SpaceParity, nil
	default:
		return serial.NoParity, fmt.Errorf("unknown parity: %s", parity)
	}
}

func parseStopBits(stopBits int) (serial.StopBits, error) {
	switch stopBits {
	case 1:
		return serial.OneStopBit, nil
	case 2:
		return serial.TwoStopBits, nil
	default:
		return serial.OneStopBit, fmt.Errorf("unsupported stop bits: %d", stopBits)
	}
}

// convertSendData converts input string data to byte slice based on format
func convertSendData(data string, format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "hex":
		// Remove spaces and convert hex string to bytes
		clean := strings.ReplaceAll(data, " ", "")
		clean = strings.ReplaceAll(clean, "\n", "")
		clean = strings.ReplaceAll(clean, "\r", "")
		if len(clean)%2 != 0 {
			clean = "0" + clean
		}
		return hex.DecodeString(clean)

	case "ascii":
		return []byte(data), nil

	case "bin":
		// Binary format: space-separated 0/1 groups, e.g. "01000001 01000010"
		clean := strings.ReplaceAll(data, "\n", "")
		clean = strings.ReplaceAll(clean, "\r", "")
		parts := strings.Fields(clean)
		bytes := make([]byte, len(parts))
		for i, p := range parts {
			var b byte
			for j := 0; j < 8 && j < len(p); j++ {
				if p[j] == '1' {
					b |= 1 << (7 - j)
				}
			}
			bytes[i] = b
		}
		return bytes, nil

	case "dec":
		// Decimal format: space-separated numbers, e.g. "65 66 67"
		clean := strings.ReplaceAll(data, "\n", "")
		clean = strings.ReplaceAll(clean, "\r", "")
		parts := strings.Fields(clean)
		bytes := make([]byte, len(parts))
		for i, p := range parts {
			var val int
			if _, err := fmt.Sscanf(p, "%d", &val); err != nil {
				return nil, fmt.Errorf("invalid decimal value at position %d: %s", i, p)
			}
			if val < 0 || val > 255 {
				return nil, fmt.Errorf("decimal value out of range (0-255) at position %d: %d", i, val)
			}
			bytes[i] = byte(val)
		}
		return bytes, nil

	case "bcd":
		// BCD format: hex string where each byte represents two BCD digits
		clean := strings.ReplaceAll(data, " ", "")
		clean = strings.ReplaceAll(clean, "\n", "")
		clean = strings.ReplaceAll(clean, "\r", "")
		if len(clean)%2 != 0 {
			clean = "0" + clean
		}
		return hex.DecodeString(clean)

	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// init registers the serial:status and serial:data events with the binding generator
func init() {
	application.RegisterEvent[SerialStatus]("serial:status")
	application.RegisterEvent[map[string]interface{}]("serial:data")
}
