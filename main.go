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

//go:embed warp-gui.jpeg
var f embed.FS
var activeIcon, _ = f.ReadFile("warp-gui.jpeg")

//go:embed warp-gui-inactive.jpg
var f2 embed.FS
var inactiveIcon, _ = f2.ReadFile("warp-gui-inactive.jpg")

type Ui struct {
	canvas *fyne.Container
	reconnectEnabled bool
	window fyne.Window
}

func NewUi(window fyne.Window) *Ui {
	box := container.NewVBox()
	return &Ui{canvas: box, reconnectEnabled: false, window: window}
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
				u.window.SetIcon(fyne.NewStaticResource("icon", activeIcon))
				button.SetText("Disconnect")
				button.OnTapped = func() {
					u.DisconnectWarp()
				}
				button.Refresh()
			} else {
				u.window.SetIcon(fyne.NewStaticResource("icon", inactiveIcon))
				button.SetText("Connect")
				button.OnTapped = func() {
					u.ConnectWarp()
				}

				// Try to reconnect
				if u.reconnectEnabled {
					if _, err := u.ConnectWarp(); err != nil {
						log.Fatal(err)
					}
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
	w.SetIcon(fyne.NewStaticResource("icon", activeIcon))
	w.Resize(fyne.NewSize(300, 100))
	w.SetFixedSize(true)

	ui := NewUi(w)
	ui.render()
	w.SetContent(ui.canvas)
	w.ShowAndRun()
}
