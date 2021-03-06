package lib

import(
	"fmt"
	"sync"
	"github.com/sparrc/go-ping"
	"time"
)
var numPingsFinished int
var numPings int

// CREATE CHANNELS FOR JOBS AND RESULTS
var jobs = make(chan string, 10000)
var pongs = make(chan *pong, 10000)

var timeout time.Duration

// RUN A PING SCAN
func PingScan(outFileName string, workerCount int, targetFileName string, timeoutSeconds int) (error) {

	timeout = time.Duration(timeoutSeconds) * time.Second
	// GET SUBNET LIST
	var targets[] string
	var err error

	if targetFileName == "" {
		// GET TARGET LIST
		Debug("Using RFC1918 subnets as targets")
		targets, err = makeHostList(RFC1918Subnets)
	}else {
		// PARSE FROM FILE
		Debug("Reading targets from file " + targetFileName)
		targets, err = getTargetsFromFile(targetFileName)
	}
	
	if err != nil{
		return err
	}
	Debug(fmt.Sprintf("Identified %d targets :)", len(targets)))
	numPings = len(targets)
	numPingsFinished = 0
	// ALLOCATE TASKS TO WORKERS
	Debug("Allocating tasks")
	go allocate(targets)

	// HANDLE RESULTS OF WORKER THREADS
	done := make(chan bool)
	go handlePongs(done, outFileName)

	// START WORKERS
	createWorkerPool(workerCount)
	<- done

	
	return nil
}

// FUNCTION TO DO THE PING 
func doPingScan(IP string)  {
	pinger, err := ping.NewPinger(IP)
	if err != nil {
		pongs <- &pong {
			IP: IP,
			Alive: false,
			Error: err,
		}
	}
	pinger.Count = 1 // ONLY PING ONCE
	pinger.Timeout = timeout
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
		if numPingsFinished % 170000 == 0 {
			Debug(fmt.Sprintf("Finished %f percent of %d pings", (float64(numPingsFinished)/float64(numPings))*100.0, numPings))
		}
	}
	pinger.Run()
}

// A WORKER
func worker(wg *sync.WaitGroup){
	
	for job := range jobs {
		doPingScan(job)
	}

	wg.Done()
}

// CREATE A WORKER POOL SIZED BY THE USER
func createWorkerPool(workerCount int) {

	var wg sync.WaitGroup
	for i:= 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	wg.Wait()
	close(pongs)
}

// PUSH JOBS INTO THE JOB CHANNEL
func allocate(jobArr []string) {
	for i := 0; i < len(jobArr); i++ {
		job := jobArr[i]
		jobs <- job
	}
	close(jobs)
}

// HANDLE ALL PONGS PUSHED TO THE CHANNEL BY WORKERS
func handlePongs(done chan bool, outFileName string) {

	var liveHosts[] string

	// CATCH ALL RESPONSES
	for pong := range pongs{

		if pong.Alive == true {
			Debug("Host " + pong.IP + " is alive!")
			liveHosts = append(liveHosts, pong.IP)
		}
		if pong.Error != nil {
		} else if pong.Alive == false {
			//Debug("Host " + pong.IP + " is dead!")
		}
	}

	// WRITE TO OUTPUT FILE
	writeLiveHosts(liveHosts, outFileName)
	done <- true
	
}

