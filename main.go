package main

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/getlantern/systray"
	"github.com/lxn/win"
	"gopkg.in/yaml.v3"
)

type Window struct {
	Hwnd  win.HWND
	Title string
}

const (
	SPI_GETWORKAREA = 48
	SPI_SETWORKAREA = 47
)

const (
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	MOD_WIN     = 0x0008
	WM_HOTKEY   = 0x0312
)

type Config struct {
	Padding int32 `yaml:"padding"`
	Hotkeys struct {
		Workspace1  string `yaml:"workspace1"`
		Workspace2  string `yaml:"workspace2"`
		Workspace3  string `yaml:"workspace3"`
		Workspace4  string `yaml:"workspace4"`
		Workspace5  string `yaml:"workspace5"`
		Workspace6  string `yaml:"workspace6"`
		Workspace7  string `yaml:"workspace7"`
		Workspace8  string `yaml:"workspace8"`
		Workspace9  string `yaml:"workspace9"`
		MoveTo1     string `yaml:"moveTo1"`
		MoveTo2     string `yaml:"moveTo2"`
		MoveTo3     string `yaml:"moveTo3"`
		MoveTo4     string `yaml:"moveTo4"`
		MoveTo5     string `yaml:"moveTo5"`
		MoveTo6     string `yaml:"moveTo6"`
		MoveTo7     string `yaml:"moveTo7"`
		MoveTo8     string `yaml:"moveTo8"`
		MoveTo9     string `yaml:"moveTo9"`
		ResizeFull  string `yaml:"resizeFull"`
		ResizeHalf  string `yaml:"resizeHalf"`
		ResizeSmall string `yaml:"resizeSmall"`
	} `yaml:"hotkeys"`
}

var currentHash = make(map[int][]Window)
var currentConfig Config

func main() {
	systray.Run(onReady, onExit)
}

func loadConfig() {
	currentConfig = Config{}
	
	// Default config
	currentConfig.Padding = 15
	currentConfig.Hotkeys.Workspace1 = "Alt+1"
	currentConfig.Hotkeys.Workspace2 = "Alt+2"
	currentConfig.Hotkeys.Workspace3 = "Alt+3"
	currentConfig.Hotkeys.Workspace4 = "Alt+4"
	currentConfig.Hotkeys.Workspace5 = "Alt+5"
	currentConfig.Hotkeys.Workspace6 = "Alt+6"
	currentConfig.Hotkeys.Workspace7 = "Alt+7"
	currentConfig.Hotkeys.Workspace8 = "Alt+8"
	currentConfig.Hotkeys.Workspace9 = "Alt+9"
	currentConfig.Hotkeys.MoveTo1 = "Alt+Shift+1"
	currentConfig.Hotkeys.MoveTo2 = "Alt+Shift+2"
	currentConfig.Hotkeys.MoveTo3 = "Alt+Shift+3"
	currentConfig.Hotkeys.MoveTo4 = "Alt+Shift+4"
	currentConfig.Hotkeys.MoveTo5 = "Alt+Shift+5"
	currentConfig.Hotkeys.MoveTo6 = "Alt+Shift+6"
	currentConfig.Hotkeys.MoveTo7 = "Alt+Shift+7"
	currentConfig.Hotkeys.MoveTo8 = "Alt+Shift+8"
	currentConfig.Hotkeys.MoveTo9 = "Alt+Shift+9"
	currentConfig.Hotkeys.ResizeFull = "Alt+Q"
	currentConfig.Hotkeys.ResizeHalf = "Alt+W"
	currentConfig.Hotkeys.ResizeSmall = "Alt+E"

	data, err := os.ReadFile("config.yaml")
	if err == nil {
		yaml.Unmarshal(data, &currentConfig)
	} else {
		out, _ := yaml.Marshal(currentConfig)
		os.WriteFile("config.yaml", out, 0644)
	}
}

// Placeholder icon, dont forget to change it :)
var iconData = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, 
	0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0xf3, 0xff, 
	0x61, 0x00, 0x00, 0x00, 0x19, 0x49, 0x44, 0x41, 0x54, 0x38, 0x11, 0x63, 0xf8, 0xff, 0xff, 0x3f, 
	0x03, 0x0a, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x00, 0x48, 0x15, 
	0x04, 0x08, 0x93, 0x76, 0x36, 0x75, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 
	0x60, 0x82,
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Space")
	systray.SetTooltip("Space Workspace Manager")

	mConfig := systray.AddMenuItem("Open Config", "Edit config.yaml")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	loadConfig()

	go startHotkeyLoop()

	go func() {
		for {
			select {
			case <-mConfig.ClickedCh:
				cmd := exec.Command("cmd", "/c", "start", "config.yaml")
				cmd.Start()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	for _, windows := range currentHash {
		for _, w := range windows {
			win.ShowWindow(w.Hwnd, win.SW_SHOW)
			win.SetForegroundWindow(w.Hwnd)
		}
	}
}

func parseHotkey(hotkeyStr string) (uintptr, uintptr) {
	var modifiers uintptr
	var key uintptr
	parts := strings.SplitSeq(strings.ToUpper(hotkeyStr), "+")
	for p := range parts {
		switch p {
		case "ALT":
			modifiers |= MOD_ALT
		case "SHIFT":
			modifiers |= MOD_SHIFT
		case "CTRL":
			modifiers |= MOD_CONTROL
		case "WIN":
			modifiers |= MOD_WIN
		default:
			if len(p) == 1 {
				key = uintptr(p[0])
			}
		}
	}
	return modifiers, key
}

func startHotkeyLoop() {
	runtime.LockOSThread()

	user32 := syscall.NewLazyDLL("user32.dll")
	registerHotKey := user32.NewProc("RegisterHotKey")

	reg := func(id int, hotkey string) {
		if hotkey == "" {
			return
		}
		mod, key := parseHotkey(hotkey)
		registerHotKey.Call(0, uintptr(id), mod, key)
	}

	// Register workspaces 1-9
	reg(1, currentConfig.Hotkeys.Workspace1)
	reg(2, currentConfig.Hotkeys.Workspace2)
	reg(3, currentConfig.Hotkeys.Workspace3)
	reg(4, currentConfig.Hotkeys.Workspace4)
	reg(5, currentConfig.Hotkeys.Workspace5)
	reg(6, currentConfig.Hotkeys.Workspace6)
	reg(7, currentConfig.Hotkeys.Workspace7)
	reg(8, currentConfig.Hotkeys.Workspace8)
	reg(9, currentConfig.Hotkeys.Workspace9)

	// Register move workspaces 10-18
	reg(10, currentConfig.Hotkeys.MoveTo1)
	reg(11, currentConfig.Hotkeys.MoveTo2)
	reg(12, currentConfig.Hotkeys.MoveTo3)
	reg(13, currentConfig.Hotkeys.MoveTo4)
	reg(14, currentConfig.Hotkeys.MoveTo5)
	reg(15, currentConfig.Hotkeys.MoveTo6)
	reg(16, currentConfig.Hotkeys.MoveTo7)
	reg(17, currentConfig.Hotkeys.MoveTo8)
	reg(18, currentConfig.Hotkeys.MoveTo9)

	// Register resizing scripts 19-21
	reg(19, currentConfig.Hotkeys.ResizeFull)
	reg(20, currentConfig.Hotkeys.ResizeHalf)
	reg(21, currentConfig.Hotkeys.ResizeSmall)

	var windowList []Window

	enumFunc := syscall.NewCallback(func(hwnd win.HWND, lParam uintptr) uintptr {
		if !win.IsWindowVisible(hwnd) {
			return 1
		}
		length := win.SendMessage(hwnd, win.WM_GETTEXTLENGTH, 0, 0)
		if length == 0 {
			return 1
		}
		buffer := make([]uint16, length+1)
		win.SendMessage(
			hwnd,
			win.WM_GETTEXT,
			uintptr(length+1),
			uintptr(unsafe.Pointer(&buffer[0])),
		)
		title := syscall.UTF16ToString(buffer)

		lowerTitle := strings.ToLower(title)
		if strings.Contains(lowerTitle, "windows input experience") ||
			strings.Contains(lowerTitle, "nvidia geforce overlay") ||
			strings.Contains(lowerTitle, "program manager") ||
			strings.Contains(lowerTitle, "yasbbar") ||
			strings.Contains(lowerTitle, "nahimic") {
			return 1
		}

		windowList = append(windowList, Window{Hwnd: hwnd, Title: title})
		return 1
	})

	enumWindows := user32.NewProc("EnumWindows")
	currentWorkspace := 1
	var msg win.MSG

	for {
		ret := win.GetMessage(&msg, 0, 0, 0)
		if ret == 0 {
			break
		}

		if msg.Message == WM_HOTKEY {
			pressedID := int(msg.WParam)

			if pressedID >= 19 && pressedID <= 21 {
				hwnd := win.GetForegroundWindow()
				if hwnd != 0 {
					win.ShowWindow(hwnd, win.SW_RESTORE)

					var workArea win.RECT
					win.SystemParametersInfo(SPI_GETWORKAREA, 0, unsafe.Pointer(&workArea), 0)

					pad := currentConfig.Padding

					var x, y, w, h int32
					fullW := workArea.Right - workArea.Left
					fullH := workArea.Bottom - workArea.Top

					switch pressedID {
					case 19:
						// Full padded
						w = fullW - (pad * 2)
						h = fullH - (pad * 2)
						x = workArea.Left + pad
						y = workArea.Top + pad
					case 20:
						// half
						w = (fullW - (pad * 3)) / 2
						h = fullH - (pad * 2)
						x = workArea.Left + pad
						y = workArea.Top + pad
					case 21:
						// small
						w = fullW / 2
						h = fullH / 2
						x = workArea.Left + (fullW - w) / 2
						y = workArea.Top + (fullH - h) / 2
					}

					win.SetWindowPos(hwnd, 0, x, y, w, h, win.SWP_NOZORDER|win.SWP_NOACTIVATE)
				}
				continue
			}

			targetWorkspace := pressedID
			isMove := false
			if pressedID >= 10 && pressedID <= 18 {
				targetWorkspace = pressedID - 9
				isMove = true
			}

			if targetWorkspace >= 1 && targetWorkspace <= 9 && targetWorkspace != currentWorkspace {
				if isMove {
					hwnd := win.GetForegroundWindow()
					if hwnd != 0 {
						length := win.SendMessage(hwnd, win.WM_GETTEXTLENGTH, 0, 0)
						title := ""
						if length > 0 {
							buffer := make([]uint16, length+1)
							win.SendMessage(hwnd, win.WM_GETTEXT, uintptr(length+1), uintptr(unsafe.Pointer(&buffer[0])))
							title = syscall.UTF16ToString(buffer)
						}

						lowerTitle := strings.ToLower(title)
						isSystemWindow := false
						if strings.Contains(lowerTitle, "windows input experience") ||
							strings.Contains(lowerTitle, "nvidia geforce overlay") ||
							strings.Contains(lowerTitle, "program manager") ||
							strings.Contains(lowerTitle, "yasbbar") ||
							strings.Contains(lowerTitle, "nahimic") || length == 0 {
							isSystemWindow = true
						}

						if !isSystemWindow {
							win.ShowWindow(hwnd, win.SW_HIDE)
							currentHash[targetWorkspace] = append(currentHash[targetWorkspace], Window{Hwnd: hwnd, Title: title})
						}
					}
				}

				windowList = nil
				enumWindows.Call(enumFunc, 0)

				currentWindows := make([]Window, len(windowList))
				copy(currentWindows, windowList)
				currentHash[currentWorkspace] = currentWindows

				for _, w := range currentHash[currentWorkspace] {
					win.ShowWindow(w.Hwnd, win.SW_HIDE)
				}

				for _, w := range currentHash[targetWorkspace] {
					win.ShowWindow(w.Hwnd, win.SW_SHOW)
					win.SetForegroundWindow(w.Hwnd)
				}

				currentWorkspace = targetWorkspace
			}
		}

		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}
