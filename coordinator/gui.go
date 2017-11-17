package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jroimartin/gocui"
)

type GuiCoordinator struct {
	Controller *Controller
	mainGui    *gocui.Gui
}

func setLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	if view, err := g.SetView("log", 0, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "Log"
		view.FgColor = gocui.ColorCyan
		g.SetCurrentView("log")
	}
	if view, err := g.SetView("peers", 0, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "Peers"
		view.BgColor = gocui.ColorGreen
	}
	return nil
}

func (gui *GuiCoordinator) StartGui(c *Controller) {
	g, err := gocui.NewGui(gocui.Output256)
	gui.mainGui = g
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(setLayout)

	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (gui *GuiCoordinator) PrintToView(v string, s string) {
	gui.mainGui.Update(func(g *gocui.Gui) error {
		if v, e := gui.mainGui.View(v); e == nil {
			fmt.Fprintln(v, s)
			return nil
		} else {
			return e
		}
	})

}

func (gui *GuiCoordinator) UpdatePings(pings map[string]int64) {
	gui.mainGui.Update(func(g *gocui.Gui) error {
		if v, e := gui.mainGui.View("peers"); e == nil {
			v.Clear()
			for k, p := range pings {
				fmt.Fprintln(v, k+": "+strconv.FormatInt(p, 10)+" ms")
			}
			return nil
		} else {
			return e
		}
	})
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}
