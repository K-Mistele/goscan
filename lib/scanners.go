package lib

import(
	"fmt"
	"github.com/sparrc/go-ping"
)

func PortScan() (error) {
	Debug("Running Port Scan")
	

	return nil
}

func PingScan(outFileName string) (error) {
	Debug("Running Ping Scan")

	hostList, err := makeHostList(RFC1918Subnets)
	if err != nil {
		return err
	}

	Debug(fmt.Sprintf("Ping scanning %d hosts :)", len(hostList)))

	// CREATE CHANNEL FOR RESPONSES
	pongs := make(chan *Pong)

	// START ALL GOROUTINES
	for _, host := range hostList {
		go Ping(host, pongs)
	}

	var liveHosts []string

	// CATCH ALL RESPONSES
	for i := 0; i < len(hostList); i++ {
		pong := <-pongs
		if pong.Alive == true {
			Debug("Host " + pong.IP + " is alive!")
			liveHosts = append(liveHosts, pong.IP)
		}
		if pong.Error != nil {
		} else if pong.Alive == false {
			//Debug("Host " + pong.IP + " is dead!")
		}
	}

	writeLiveHosts(liveHosts, outFileName)
	return nil
}

// GOROUTING TO DO A PING
func Ping(IP string, pongs chan *Pong) {

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
	}
	pinger.Run()


}