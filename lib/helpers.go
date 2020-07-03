package lib

import(
	"log"
	"net"
	"bufio"
	"os"
)

func Debug(msg string) {
	log.Println("[+]", msg)
}

func Warn(msg string) {
	log.Println("[!]", msg)
}

// RETURN A SLICE OF STRINGS OF ALL THE IP ADDRESSES IN A CIDR BLOCK, SANS THE NETWORK AND BROADCAST ADDRESS
func hostsInCidr(cidr string) ([]string, error) {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	// ADD IP ADDRESS
	var ipAddresses []string
	for ip := ip.Mask(network.Mask); network.Contains(ip); increment(ip) {
		ipString := ip.String()

		// EXCLUDE BROADCAST ADDRESSES
		if ipString[len(ipString)-4 : len(ipString)] != ".255" {
			ipAddresses = append(ipAddresses, ipString)
		}
		
	}

	// REMOVE IP ADDRESS AND BROADCAST ADDRESS
	return ipAddresses, nil

}

// INCREMENT AN IP ADDRESS
func increment(ip net.IP) {
	for j:= len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// CREATE A HOST LIST FROM AN ARRAY OF CIDR NOTATIONS
func makeHostList(cidrRanges []string) ([]string, error) {

	Debug("Creating host list")
	var hostList []string 

	for _, subnetCIDR := range cidrRanges {
		Debug("Adding subnet " + subnetCIDR + " targets")
		hostsInSubnet, err := hostsInCidr(subnetCIDR)
		if err != nil {
			return nil, err
		}
		hostList = append(hostList, hostsInSubnet...)
		
	}

	return hostList, nil

}

// WRITE ALL THE HOSTS THAT WERE DISCOVERED TO BE ALIVE TO AN OUTPUT FILE
func writeLiveHosts(addresses []string, filename string) {
	
	Debug("Writing out live hosts to file " + filename)
	fileHandle, err := os.Create(filename)
	check(err)
	writer := bufio.NewWriter(fileHandle)

	defer fileHandle.Close()

	for _, address := range addresses {

		_, err := writer.WriteString(address + "\n")
		check(err)
		
	}
	writer.Flush()
}

// PANIC IF THERE'S AN ERROR
func check (e error) {
	if e != nil {
		panic(e)
	}
}

// PUSH TARGETS TO THE JOBS CHANNEL
func pushTargetsToChannel(targets []string, jobs chan string) { 
	for _, target := range targets {
		jobs <- target
	}

	close(jobs)
}

// GET TARGETS FROM FILE
func getTargetsFromFile(fileName string) ([]string, error) {
	var targetCIDRs []string
	var err error
	
	file, err := os.Open(fileName)
	if err != nil {
		Warn("Failed to open target file!")
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cidr := scanner.Text()
		targetCIDRs = append(targetCIDRs, cidr)
		Debug("Read subnet " + cidr + " from file")
	}

	if err := scanner.Err(); err != nil {
		Warn("Error reading file!")
		panic(err)
	}

	var hostList []string
	// NOW WE HAVE A SLICE OF CIDRS
	for _, subnetCIDR := range targetCIDRs {
		Debug("Adding subnet " + subnetCIDR + " targets")
		hostsInSubnet, err := hostsInCidr(subnetCIDR)
		if err != nil {
			return nil, err
		}
		hostList = append(hostList, hostsInSubnet...)
		
	}

	return hostList, nil
}