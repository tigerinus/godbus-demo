package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/godbus/dbus/v5"
)

func main() {
	println("godbus-demo")

	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	matchRules := []string{
		"path_namespace='/'",
	}

	var flag uint

	println("trying to become a monitor...")
	call := conn.BusObject().Call("org.freedesktop.DBus.Monitoring.BecomeMonitor", 0, matchRules, flag)
	if call.Err != nil {
		panic(call.Err)
	}

	println("eavesdropping to channel...")
	ch := make(chan *dbus.Message, 10)
	conn.Eavesdrop(ch)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for {
		select {
		case <-ctx.Done():
		case msg, ok := <-ch:
			if !ok {
				println("channel closed")
				return
			}

			fmt.Fprintf(w, "type:\t%v\n", msg.Type)
			fmt.Fprintf(w, "flags:\t%v\n", msg.Flags)
			fmt.Fprintf(w, "headers:\t%v\n", msg.Headers)
			fmt.Fprintf(w, "body:\t%v\n", msg.Body)
			fmt.Fprintf(w, "---\n")
			w.Flush()
		}
	}
}
