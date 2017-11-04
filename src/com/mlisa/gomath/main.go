package main

import (
	"log"

	"runtime"

	"com/mlisa/gomath/message"
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/jroimartin/gocui"
)

func main() {
	//myself = getConfig().Myself
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start(getConfig().Myself.Address)

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(Receive)

	peer, err := actor.SpawnNamed(props, getConfig().Myself.Name)

	if err != nil {
		println("[PEER] Name already in use")
	}
	console.ReadLine()

	g, _ := gocui.NewGui(gocui.Output256)
	defer g.Close()
	controller := actor.Spawn(actor.FromInstance(&Controller{gui: g, peer: peer}))

	g.SetManagerFunc(layout)
	g.Cursor = true
	g.Mouse = false

	if err := initKeybindings(g, controller); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("input", 0, 0, maxX/2, maxY/3*2); err != nil {
		v.Title = "Input"
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
		v.FgColor = gocui.ColorYellow | gocui.AttrBold
		v.BgColor = gocui.ColorCyan
		g.SetCurrentView("input")
	}
	if v, err := g.SetView("Output", 0, maxY/3*2, maxX/2, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Title = "Output"
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func initKeybindings(g *gocui.Gui, controller *actor.PID) error {
	// quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		//vdst, _ := g.View("log")
		//fmt.Fprint(vdst, controller)
		controller.Tell(&message.AskForResult{v.Buffer()})
		return nil
	}); err != nil {
		return err
	}
	return nil
}
