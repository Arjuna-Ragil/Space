package main

import (
	_ "embed"
	"os/exec"

	"github.com/getlantern/systray"
	"github.com/lxn/win"
)

//go:embed spaceicon.ico
var iconData []byte

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
