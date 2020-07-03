package lib

import(
	"fmt"
	"github.com/sparrc/go-ping"
)


func PortScan() (error) {
	Debug("Running Port Scan")
	

	return nil
}

// PERFORM A PING SCAN
func derpScan(outFileName string, workers int) (error) {
	Debug("Running Ping Scan")


	// CREATE CHANNEL FOR JOBS
	jobs := make(chan string, 10000)

	// CREATE CHANNEL FOR RESPONSES
	pongs := make(chan *pong, 10000)

	// CREATE THE TARGET LIST BY PUSHING ALL IPS (EXCEPT NETOWRK AND BROADCAST ADDRESSES) INTO A CHANNEL FOR WORKERS
	targets, err := makeHostList(RFC1918Subnets)

	go pushTargetsToChannel(targets, jobs)
	if err != nil {
		return err
	}

	// TRACK COMPLETION
	numPings = len(jobs)
	numPingsFinished = 0
	
	Debug(fmt.Sprintf("Ping scanning %d hosts :)", len(jobs)))
	

	// START ALL WORKERS, limited to prevent too many open file descriptors
	for i := 0; i < workers; i++ {
		go PingWorker(jobs, pongs)
	}
	

	var liveHosts []string

	// CATCH ALL RESPONSES
	for i := 0; i < numPings; i++ {
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
func PingWorker(IPs chan string, pongs chan *pong) {

	for (len(IPs) > 0) {
		IP := <- IPs
		// CREATE A PINGER
		pinger, err := ping.NewPinger(IP)
		if err != nil {
			pongs <- &pong {
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

			pongs <- &pong {
				IP: IP,
				Alive: alive,
				Error: nil,
			}
			numPingsFinished++
			if numPingsFinished % 10000 == 0 {
				Debug(fmt.Sprintf("Finished %d out of %d pings", numPingsFinished, numPings))
			}
		}
		pinger.Run()
	}
	


}