package lib

const Banner = `
  ________           _________                        
 /  _____/   ____   /   _____/  ____  _____     ____  
/   \  ___  /  _ \  \_____  \ _/ ___\ \__  \   /    \ 
\    \_\  \(  <_> ) /        \\  \___  / __ \_|   |  \
 \______  / \____/ /_______  / \___  >(____  /|___|  /
        \/                 \/      \/      \/      \/                                       
`
var RFC1918Subnets = []string{
  "10.0.0.0/8",
  "172.16.0.0/12",
  "192.168.0.0/16",
}

var testSubnets = []string{
  "10.0.1.0/24",
}