package main

import (
	"github.com/godbus/dbus/v5"
)

func main() {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rules := []string{
		"member='Notify',path='/org/freedesktop/UDisks2'",
	}

	var flag uint

	call := conn.BusObject().Call("org.freedesktop.Monitoring.BecomeMonitor", 0, rules, flag)
	if call.Err != nil {
		panic(call.Err)
	}

	ch := make(chan *dbus.Message, 10)
	conn.Eavesdrop(ch)
	for v := range ch {
		println(v)
	}
}
