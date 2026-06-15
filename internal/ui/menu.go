package ui

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/caseymrm/menuet"
	"github.com/rob121/tunneler/internal/assets"
	"github.com/rob121/tunneler/internal/config"
	"github.com/rob121/tunneler/internal/tunnel"
)

func Run(mgr *tunnel.Manager) {
	app := menuet.App()
	app.Name = "Tunneler"
	app.Label = "app.tunneler"

	iconPath, err := assets.IconPath()
	if err != nil {
		log.Printf("icon unavailable: %v", err)
		app.SetMenuState(&menuet.MenuState{Title: "Tunnels"})
	} else {
		app.SetMenuState(&menuet.MenuState{Image: iconPath})
	}

	app.Children = func() []menuet.MenuItem {
		var items []menuet.MenuItem

		for _, t := range mgr.Tunnels() {
			tun := t

			title := fmt.Sprintf("⚪ %s (%d)", tun.Name, tun.Dynamic)
			if mgr.Running(tun.Name) {
				title = fmt.Sprintf("🟢 %s (%d)", tun.Name, tun.Dynamic)
			}

			items = append(items, menuet.MenuItem{
				Text: title,
				Clicked: func() {
					_ = mgr.Toggle(tun)
					menuet.App().MenuChanged()
				},
			})
		}

		items = append(items,
			menuet.MenuItem{Type: menuet.Separator},
			menuet.MenuItem{
				Text: "Manage",
				Clicked: func() {
					path, err := config.FilePath()
					if err != nil {
						return
					}
					_ = exec.Command("open", path).Run()
				},
			},
		)

		return items
	}

	app.RunApplication()
}
