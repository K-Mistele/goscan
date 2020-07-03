package main

import(

	// REMOTE
	"fmt"
	//"net"
	"flag"
	"os/exec"

	// LOCAL
	. "./lib"
)

func main(){
	Debug("Configuring OS ping settings")
	exec.Command("sudo sysctl -w net.ipv4.ping_group_range=\"0 2147483647\"")
	exec.Command("ulimit -n 4096")
	Debug("Starting goscan!")
	fmt.Println(Banner)
	fmt.Println("Presented by Blacklight")

	// DO COMMAND LINE ARGS
	var doPortScan, doPingScan bool
	var outFileName string
	var workerCount int

	// FLAG DECLARATIONS
	flag.BoolVar(&doPortScan, "portscan", false, "Do a TCP port scan for common ports")
	flag.BoolVar(&doPingScan, "pingscan", false, "Do an ICMP echo ping scan for hosts")
	flag.StringVar(&outFileName, "outfile", "live_hosts.txt", "Specify the text file to write live hosts to")
	flag.IntVar(&workerCount, "workers", 512, "The number of workers to use. Default is 1024; limited by the number of file descriptors available to your user. Check with (ulimit -n)")

	flag.Parse()

	if doPortScan == false && doPingScan == false {
		Warn("Specify scan type with --portscan, --pingscan, or both!")
		return
	} 

	if doPortScan == true {
		PortScan()
	}

	if doPingScan == true {
		PingScan(outFileName, workerCount)
	}
}



