package main

import (
	_ "embed"
	"os"
	"os/exec"

	"github.com/getlantern/systray"
	"github.com/lxn/win"
	"golang.org/x/sys/windows/registry"
)

//go:embed documentation/spaceicon.ico
var iconData []byte

const registryAppName = "SpaceWorkspaceManager"

func isStartupEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	_, _, err = k.GetStringValue(registryAppName)
	return err == nil
}

func enableStartup() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	return k.SetStringValue(registryAppName, `"`+exePath+`"`)
}

func disableStartup() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.DeleteValue(registryAppName)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Space")
	systray.SetTooltip("Space Window Manager")

	mStartup := systray.AddMenuItemCheckbox("Run on Startup", "Start Space automatically when PC starts up", isStartupEnabled())
	mConfig := systray.AddMenuItem("Open Config", "Edit config.yaml")
	mQuit := systray.AddMenuItem("Quit", "Quit Space")

	loadConfig()

	go startPipeServer()
	go startHotkeyLoop()

	go func() {
		for {
			select {
			case <-mStartup.ClickedCh:
				if mStartup.Checked() {
					disableStartup()
					mStartup.Uncheck()
				} else {
					enableStartup()
					mStartup.Check()
				}
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
