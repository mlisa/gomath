package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/jroimartin/gocui"
)

type Controller struct {
	gui  *gocui.Gui
	peer *actor.PID
}

func (controller *Controller) computeResult(operation string) {

	//Policy per decidere se calcolarla qui o meno
}
