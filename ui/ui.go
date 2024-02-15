package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

var (
	Header        string
	MessageCount  int
	GUI           *gocui.Gui
	ClientList    []string
	latestMessage string
)

func OpenUI() {
	MessageCount = 0
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
	if maxY < 3 {
		maxY = 3
	}
	if maxX < 4 {
		maxX = 4
	}
	v, err = g.SetView("clients", 0, 1, maxX/2-1, maxY-1)
	if err == gocui.ErrUnknownView {
		v.Title = "Clients[0]"
		v.FgColor = gocui.ColorMagenta
		v.Autoscroll = true
		v.Wrap = true
	} else if err == nil {
		v.Clear()
		for _, c := range ClientList {
			fmt.Fprintln(v, c)
		}
	} else {
		return err
	}
	// Logs Viewer
	v, err = g.SetView("logs", maxX/2, 1, maxX-1, maxY-1)
	if err == gocui.ErrUnknownView {
		v.Title = "Logs[0]"
		v.Autoscroll = true
		v.FgColor = gocui.ColorGreen
		v.Wrap = true
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

func ReplaceClient(UID int, client string) {
	ClientList[UID] = client
	GUI.Update(UpdateClients)
}

func UpdateClients(g *gocui.Gui) error {
	v, err := g.View("clients")
	if err != nil {
		return err
	}
	v.Title = fmt.Sprintf("Clients[%v]", len(ClientList))
	v.Clear()
	for _, c := range ClientList {
		fmt.Fprintln(v, c)
	}
	return nil
}

func AddMessage(message string) {
	latestMessage = message
	MessageCount++
	GUI.Update(UpdateLogs)
}

func UpdateLogs(g *gocui.Gui) error {
	v, err := GUI.View("logs")
	v.Title = fmt.Sprintf("Logs[%v]", MessageCount)
	if err != nil {
		return err
	}
	fmt.Fprint(v, latestMessage)
	return nil
}
