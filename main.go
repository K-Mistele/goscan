package main

import(

	// REMOTE
	"fmt"
	//"net"
	"flag"

	// LOCAL
	. "./lib"
)

func main(){
	
	Debug("Starting goscan!")
	fmt.Println(Banner)
	fmt.Println("Presented by Blacklight")

	// DO COMMAND LINE ARGS
	var doPortScan, doPingScan bool

	// FLAG DECLARATIONS
	flag.BoolVar(&doPortScan, "portscan", false, "Do a TCP port scan for common ports")
	flag.BoolVar(&doPingScan, "pingscan", false, "Do an ICMP echo ping scan for hosts")

	flag.Parse()

	if doPortScan == false && doPingScan == false {
		Warn("Specify scan type with --portscan, --pingscan, or both!")
		return
	} 

	if doPortScan == true {
		PortScan()
	}

	if doPingScan == true {
		PingScan()
	}
}



