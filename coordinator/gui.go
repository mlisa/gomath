package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mlisa/gomath/common"
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
	}
	if view, err := g.SetView("peers", 0, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "Peers"
		view.Highlight = true
		view.SelBgColor = gocui.ColorCyan
		view.SelFgColor = gocui.ColorYellow | gocui.AttrBold
		view.Autoscroll = true
		g.SetCurrentView("peers")
	}
	return nil
}

func (gui *GuiCoordinator) StartGui(c *Controller) {
	g, err := gocui.NewGui(gocui.Output256)
	gui.mainGui = g
	gui.Controller = c
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(setLayout)

	if err := gui.initKeybindings(g); err != nil {
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

func (gui *GuiCoordinator) UpdatePings(pings map[string]common.Pong) {
	type pong struct {
		value    int64
		complete bool
	}
	gui.mainGui.Update(func(g *gocui.Gui) error {
		if v, e := gui.mainGui.View("peers"); e == nil {
			v.Clear()
			for s, _ := range gui.Controller.GetPeers() {
				ping := pings[s]
				if ping.Complete {
					fmt.Fprintln(v, s+" "+strconv.FormatInt(ping.Value, 10)+" ms")
				} else {
					fmt.Fprintln(v, s+" N/A ms")
				}
			}
		}
		if v, e := gui.mainGui.View("peer"); e == nil {
			buf := v.Buffer()
			infos := strings.Split(buf, "\n")
			name := strings.Split(infos[1], " ")[1]
			address := strings.Split(infos[2], " ")[1]
			latency := strings.Split(infos[3], " ")[1]
			v.Clear()
			pong := pings[address+"/"+name]
			fmt.Fprintln(v, "Status: OK")
			fmt.Fprintln(v, "Name: "+name)
			fmt.Fprintln(v, "Address: "+address)
			if pong.Complete {
				fmt.Fprintln(v, "Latency: "+strconv.FormatInt(pong.Value, 10)+" ms")
			} else {
				fmt.Fprintln(v, "Latency: "+latency+" ms")
			}
		}
		return nil
	})
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (gui *GuiCoordinator) scrollUp(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		return gui.scrollView(g, v, -1)
	})
	return nil
}

func (gui *GuiCoordinator) scrollDown(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		return gui.scrollView(g, v, 1)
	})
	return nil
}

func (gui *GuiCoordinator) initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	// move up
	if err := g.SetKeybinding("peers", gocui.KeyArrowUp, gocui.ModNone, gui.scrollUp); err != nil {
		return err
	}
	// move up
	if err := g.SetKeybinding("peers", 'k', gocui.ModNone, gui.scrollUp); err != nil {
		return err
	}
	// move down
	if err := g.SetKeybinding("peers", gocui.KeyArrowDown, gocui.ModNone, gui.scrollDown); err != nil {
		return err
	}
	// move down
	if err := g.SetKeybinding("peers", 'j', gocui.ModNone, gui.scrollDown); err != nil {
		return err
	}
	// show peer details
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		gui.Controller.RunPing()
		g.Update(func(g *gocui.Gui) error {
			vp, _ := g.View("peers")
			cx, _ := vp.Cursor()
			l, _ := vp.Line(cx)
			if len(l) > 0 {
				gui.newView(l, g)
			}
			return nil
		})
		return nil
	}); err != nil {
		return err
	}
	// hide peer details
	if err := g.SetKeybinding("", gocui.KeyBackspace2, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		gui.Controller.RunPing()
		if _, err := g.View("peer"); err == nil {
			if err := g.DeleteView("peer"); err != nil {
				return err
			} else {
				g.SetCurrentView("peers")
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Move the cursor of 'peers' view and show the content of the highlighted bin
func (gui *GuiCoordinator) scrollView(g *gocui.Gui, v *gocui.View, dy int) error {
	if v != nil {
		gui.Controller.RunPing()
		g.Update(func(g *gocui.Gui) error {
			_, cy := v.Cursor()
			l, _ := v.Line(cy + dy)
			if len(l) > 0 {
				moveTo(v, dy)
			}
			return nil
		})
	}
	return nil
}

func moveTo(v *gocui.View, step int) {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	_, sy := v.Size()
	offset := (sy - 1) / 2
	// Start list
	if cy <= offset || (oy == 0 && step < 0) {
		v.SetCursor(cx, cy+step)
	} else {
		var l string
		var e error
		if step > 0 {
			// End list
			l, e = v.Line(sy)
		} else {
			// Middle list
			l, e = v.Line(cy + step + offset)
		}
		if e == nil && len(l) > 0 {
			v.SetOrigin(ox, oy+step)
		} else {
			v.SetCursor(cx, cy+step)
		}
	}
}

func (gui *GuiCoordinator) newView(l string, g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := "peer"
	if view, err := g.SetView(name, 5, 5, maxX-5, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		l = strings.Split(l, " ")[0]
		pong := gui.Controller.GetLatency(l)
		view.Title = "Peer Details"
		fmt.Fprintln(view, "Status: OK")
		fmt.Fprintln(view, "Name: "+strings.Split(l, "/")[1])
		fmt.Fprintln(view, "Address: "+strings.Split(l, "/")[0])
		if pong > 0 {
			fmt.Fprintln(view, "Latency: "+strconv.FormatInt(pong, 10)+" ms")
		} else {
			fmt.Fprintln(view, "Latency: N/A ms")
		}
	}
	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}
	return nil
}
