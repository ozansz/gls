package gui

import (
	"fmt"
	"time"
)

func setLastLog(text string) {
	if currLastLogView == nil {
		return
	}
	currLastLogView.SetText(fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), text))
}

func setError(text string) {
	setLastLog(fmt.Sprintf("[ERROR] %s", text))
}

func setInfo(text string) {
	setLastLog(fmt.Sprintf("[INFO] %s", text))
}
