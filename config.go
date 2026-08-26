package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Padding int32 `yaml:"padding"`
	IgnoreList []string `yaml:"ignoreList"`
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

var currentConfig Config

func getConfigPath() string {
	exePath, err := os.Executable()
	if err != nil {
		return "config.yaml"
	}
	return filepath.Join(filepath.Dir(exePath), "config.yaml")
}

func loadConfig() {
	currentConfig = Config{}
	
	// Default config
	currentConfig.Padding = 15
	currentConfig.IgnoreList = []string{}
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

	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err == nil {
		err = yaml.Unmarshal(data, &currentConfig)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		out, _ := yaml.Marshal(currentConfig)
		err = os.WriteFile(configPath, out, 0644)
		if err != nil {
			fmt.Println(err)
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
