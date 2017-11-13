package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func setLayout(g *gocui.Gui) error {
	g.Cursor = true
	maxX, maxY := g.Size()
	if view, err := g.SetView("log", 0, 0, maxX-1, maxY/2-1); err != nil {
		view.Title = "Log"
		if err != gocui.ErrUnknownView {
			return err
		}
		view.SelBgColor = gocui.ColorGreen
		view.SelFgColor = gocui.ColorBlack
		view.Editable = true
	} else {
		return err
	}
	if view, err := g.SetView("peers", 0, maxY/2, maxX-1, maxY-1); err != nil {
		view.Title = "Peers"
		if err != gocui.ErrUnknownView {
			return err
		}
		view.SelBgColor = gocui.ColorGreen
		view.SelFgColor = gocui.ColorBlack
		view.Editable = true
	} else {
		return err
	}
	g.SetCurrentView("log")
	return nil
}

func StartGui() *gocui.Gui {
	g, err := gocui.NewGui(gocui.Output256)
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
	return g
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
