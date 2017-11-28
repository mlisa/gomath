package main

import (
	"log"
	"runtime"

	"github.com/mlisa/gomath/common"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/jroimartin/gocui"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	config = kingpin.Flag("config", "Configuration file for peer").Short('c').Default("config_peer.json").String()
)

func main() {
	kingpin.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	var config, err = common.GetFileConfig(*config)
	if err != nil {
		kingpin.FatalUsage("Wrong usage, please see the help")
	}
	remote.Start(config.Address)

	g, _ := gocui.NewGui(gocui.Output256)
	controller := Controller{Gui: g, Cache: &CacheManager{}, Config: config, Connected: true}
	defer g.Close()

	props := actor.FromInstance(&Peer{Controller: &controller, computeCapability: config.ComputeCapability})
	peer, err := actor.SpawnNamed(props, config.Id)
	if err != nil {
		println("[PEER] Name already in use")
	}

	controller.Peer = peer

	g.SetManagerFunc(layout)
	if err := initKeybindings(g, &controller); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true
	g.Mouse = false

	if v, err := g.SetView("input", 0, 0, maxX-1, maxY/6-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Input"
		v.Editable = true
		v.Wrap = true
		v.FgColor = gocui.AttrBold
		g.SetCurrentView("input")
	}
	if v, err := g.SetView("Output", 0, maxY/6, maxX-1, maxY/3-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Title = "Output"
		v.Autoscroll = true
	}
	if v, err := g.SetView("Log", 0, maxY/3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Title = "Log"
		v.Autoscroll = true
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func initKeybindings(g *gocui.Gui, controller *Controller) error {
	// quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.Update(func(g *gocui.Gui) error {
			controller.AskForResult(v.Buffer())
			return nil
		})
		return nil
	}); err != nil {
		return err
	}
	return nil
}
