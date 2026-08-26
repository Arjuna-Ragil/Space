package main

import (
	_ "embed"
	"fmt"
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
	defer func(){
		err := k.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	_, _, err = k.GetStringValue(registryAppName)
	return err == nil
}

func enableStartup() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func(){
		err := k.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
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
	defer func(){
		err := k.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
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
					err := disableStartup()
					if err != nil {
						fmt.Println(err)
					}
					mStartup.Uncheck()
				} else {
					err := enableStartup()
					if err != nil {
						fmt.Println(err)
					}
					mStartup.Check()
				}
			case <-mConfig.ClickedCh:
				cmd := exec.Command("cmd", "/c", "start", "", getConfigPath())
				err := cmd.Start()
				if err != nil {
					fmt.Println(err)
				}
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
