package peer

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/mlisa/gomath/message"
	"github.com/mlisa/gomath/parser"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/jroimartin/gocui"
)

type Controller struct {
	Gui   *gocui.Gui
	Peer  *actor.PID
	Cache *CacheManager
}

func (controller *Controller) AskForResult(operation string) {
	var something = true
	if something {
		controller.Peer.Tell(&message.AskForResult{operation})
		controller.Gui.Update(func(g *gocui.Gui) error {
			output, _ := g.View("Log")
			output.Clear()
			fmt.Fprint(output, "Sent AskForResult message to peer")
			return nil
		})
	} else {
		result, err := parser.ParseReader("", bytes.NewBufferString(operation))
		if err == nil {
			controller.Gui.Update(func(g *gocui.Gui) error {
				output, _ := g.View("Output")
				output.Clear()
				fmt.Fprint(output, result)
				return nil
			})
			log.Println("[CONTROLLER] Saved result :" + strconv.Itoa(result.(int)))
			controller.Cache.addNewOperation(operation, strconv.Itoa(result.(int)))
		}
	}
}

func (controller *Controller) SearchInCache(operation string) string {
	controller.Gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Log")
		fmt.Fprint(output, "Received SearchInCache message from peer")
		return nil
	})

	if result, err := controller.Cache.retrieveResult(operation); err == nil {
		log.Println("Retrieved result " + result + " from cache")
		return result
	}
	return ""
}

func (controller *Controller) SetResult(result string) {
	log.Println("Received result " + result)
	controller.Gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Output")
		fmt.Fprint(output, result)
		return nil
	})
	controller.Gui.Update(func(g *gocui.Gui) error {
		output, _ := g.View("Log")
		fmt.Fprint(output, "Receive Response message from peer")
		return nil
	})
}
