package main

import (
	"embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os/exec"
	"strings"
	"time"
)

//go:embed warp-logo.jpeg
var f embed.FS
var file, _ = f.ReadFile("warp-logo.jpeg")

type Ui struct {
	canvas *fyne.Container
}

func NewUi() *Ui {
	box := container.NewVBox()
	return &Ui{canvas: box}
}

func (u Ui) render() {
	statusLabel := widget.NewLabel("Status: Unknown")
	u.canvas.Add(statusLabel)

	button := widget.NewButton("Connect", func() {
	})

	u.canvas.Add(button)

	go func() {
		for {
			status, err := u.getWarpStatus()
			if err != nil {
				log.Fatal(err)
			}
			statusLabel.SetText("Status: " + status)
			statusLabel.Refresh()

			if status == "Connected" {
				button.SetText("Disconnect")
				button.OnTapped = func() {
					u.DisconnectWarp()
				}
				button.Refresh()
			} else {
				button.SetText("Connect")
				button.OnTapped = func() {
					u.ConnectWarp()
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

func (u Ui) getWarpStatus() (string, error) {
	cmd := exec.Command("warp-cli", "status")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	status := strings.Split(string(stdout), ":")[1]
	status = strings.TrimSpace(status)

	return status, err
}

func (u Ui) ConnectWarp() (string, error) {
	cmd := exec.Command("warp-cli", "connect")
	stdout, err := cmd.Output()
	fmt.Println(string(stdout))
	return string(stdout), err
}

func (u Ui) DisconnectWarp() (string, error) {
	cmd := exec.Command("warp-cli", "disconnect")
	stdout, err := cmd.Output()
	fmt.Println(string(stdout))
	return string(stdout), err
}

func main() {
	a := app.New()
	w := a.NewWindow("Warp")

	w.SetIcon(fyne.NewStaticResource("icon", file))
	w.Resize(fyne.NewSize(300, 100))
	w.SetFixedSize(true)

	ui := NewUi()
	ui.render()
	w.SetContent(ui.canvas)
	w.ShowAndRun()
}
