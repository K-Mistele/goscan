package main

import(
	"fmt"
	"github.com/sparrc/go-ping"

)
type Pong struct {
	IP		string
	Alive	bool
	Error	error
}
func main () {
	IP := "10.0.2.12"
	pongs := make(chan *Pong, 12)
		// CREATE A PINGER
		pinger, err := ping.NewPinger(IP)
		if err != nil {
			pongs <- &Pong {
				IP: IP,
				Alive: false,
				Error: err,
			}
		}
		pinger.Count = 1 // ONLY PING ONCE
		pinger.Timeout = 200000000
		// WHEN PING IS DONE
		pinger.OnFinish = func(stats *ping.Statistics) {
			var alive bool;
			alive = stats.PacketsRecv > 0
	
			pongs <- &Pong {
				IP: IP,
				Alive: alive,
				Error: nil,
			}

				fmt.Println(fmt.Sprintf("Finished %d out of %d pings", 1, 1))


		}
		pinger.Run()
		fmt.Println("Test")
}