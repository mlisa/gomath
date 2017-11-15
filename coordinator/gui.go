package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

type Gui struct {
	Controller *Controller
}

func setLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	if view, err := g.SetView("log", 0, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "Log"
		view.SelBgColor = gocui.ColorGreen
		view.SelFgColor = gocui.ColorBlack
		g.SetCurrentView("log")
	}
	if view, err := g.SetView("peers", 0, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "Peers"
		view.SelBgColor = gocui.ColorGreen
		view.SelFgColor = gocui.ColorBlack
	}
	return nil
}

func StartGui(c *Controller) {
	g, err := gocui.NewGui(gocui.Output256)
	c.Gui = g
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
