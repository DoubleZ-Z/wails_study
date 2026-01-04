package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) getMenu() *menu.Menu {
	m := menu.NewMenu()
	fileMenu := m.AddSubmenu("file")
	fileMenu.AddText("open file", keys.CmdOrCtrl("o"), func(data *menu.CallbackData) {
		fmt.Println("open file")
	})
	fileMenu.AddText("save file", keys.Control("s"), func(data *menu.CallbackData) {
		fmt.Println("save file")
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("exit", &keys.Accelerator{}, func(data *menu.CallbackData) {
		fmt.Println("exit")
	})
	return m
}

func (a *App) getFileMenu() *menu.Menu {
	m := menu.NewMenu()
	fileMenu := m.AddSubmenu("file")
	fileMenu.AddText("open file", keys.CmdOrCtrl("o"), func(data *menu.CallbackData) {
		fileDialog, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
			Title: "Open File",
			Filters: []runtime.FileFilter{
				{
					DisplayName: "Image Files (*.jpg, *.png)",
					Pattern:     "*.jpg;*.png",
				},
			},
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(fileDialog)
	})
	fileMenu.AddText("save file", keys.Control("s"), func(data *menu.CallbackData) {
		fmt.Println("save file")
	})
	return m
}

func (a *App) getScreen() *menu.Menu {
	m := menu.NewMenu()
	submenu := m.AddSubmenu("screen")
	submenu.AddText("全屏", keys.CmdOrCtrl("F11"), func(data *menu.CallbackData) {
		if runtime.WindowIsFullscreen(a.ctx) {
			runtime.WindowUnfullscreen(a.ctx)
		} else {
			runtime.WindowFullscreen(a.ctx)
		}
	})
	return m
}
