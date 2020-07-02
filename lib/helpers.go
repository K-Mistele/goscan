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

func HostsInCidr(cidr string) ([]string, error) {
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

func increment(ip net.IP) {
	for j:= len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func makeHostList(cidrRanges []string) ([]string, error) {

	Debug("Creating host list")
	var hostList []string 

	for _, subnetCIDR := range cidrRanges {
		Debug("Adding subnet " + subnetCIDR + " targets")
		hostsInSubnet, err := HostsInCidr(subnetCIDR)
		if err != nil {
			return nil, err
		}
		hostList = append(hostList, hostsInSubnet...)
		
	}

	return hostList, nil

}

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

func check (e error) {
	if e != nil {
		panic(e)
	}
}