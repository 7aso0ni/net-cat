package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

var (
	Header        string
	Count         int
	GUI           *gocui.Gui
	ClientList    []string
	latestMessage string
)

func OpenUI() {
	Count = 0
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		g.Close()
		os.Exit(0)
	}()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		fmt.Println(err.Error())
	}

	GUI = g
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Println(err.Error())
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// Header Viewer
	v, err := g.SetView("header", maxX/2-len(Header)/2-1, -1, maxX/2+len(Header)/2+1, 2)
	if err == gocui.ErrUnknownView {
		v.Frame = false
		fmt.Fprintln(v, Header)
	} else if err != nil {
		return err
	}
	// Client Viewer
	v, err = g.SetView("clients", 0, 1, maxX/2-1, maxY-1)
	if err == gocui.ErrUnknownView {
		v.Title = "Clients"
	} else if err == nil {
		v.Clear()
		fmt.Fprintf(v, "Count: %v\n", len(ClientList))
		for _, c := range ClientList {
			fmt.Fprintln(v, c)
		}
	} else {
		return err
	}
	// Logs Viewer
	v, err = g.SetView("logs", maxX/2, 1, maxX-1, maxY-1)
	if err == gocui.ErrUnknownView {
		v.Title = "Logs"
	} else if err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func AddClient(client string) {
	ClientList = append(ClientList, client)
	GUI.Update(UpdateClients)
}

func DeleteClient(UID int) {
	ClientList = append(ClientList[:UID], ClientList[UID+1:]...)
	GUI.Update(UpdateClients)
}

func UpdateClients(g *gocui.Gui) error {
	v, err := g.View("clients")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(v, "Count: %v\n", len(ClientList))
	for _, c := range ClientList {
		fmt.Fprintln(v, c)
	}
	return nil
}

func AddMessage(message string) {
	latestMessage = message
	GUI.Update(UpdateLogs)
}

func UpdateLogs(g *gocui.Gui) error {
	v, err := GUI.View("logs")
	if err != nil {
		return err
	}
	fmt.Fprint(v, latestMessage)
	return nil
}
