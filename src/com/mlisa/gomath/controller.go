package main

import (
	"bytes"
	"com/mlisa/gomath/message"
	"com/mlisa/gomath/parser"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/jroimartin/gocui"
	"log"
)

type Controller struct {
	gui  *gocui.Gui
	peer *actor.PID
}

func (controller *Controller) Receive(context actor.Context) {

	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Println("[CONTROLLER] Started, initialize actor here, I'm " + context.Self().Id + " " + context.Self().Address)

	case *message.AskForResult:
		//Policy per decidere se calcolarlo qui o meno
		var something = false
		if something {
			controller.peer.Request(msg, context.Self())
		} else {
			res, err := parser.ParseReader("", bytes.NewBufferString(msg.Operation))
			if err == nil {
				controller.gui.Update(func(g *gocui.Gui) error {
					output, _ := g.View("Output")
					fmt.Fprint(output, res)
					return nil
				})
			}
		}

	case *message.Response:
		//setResult(controller.gui, msg.Result)

	}
}
