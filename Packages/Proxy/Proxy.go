package Proxy

import (
	"fmt"
	"net"

	"chicha/Packages/Config"
)

func ProxyDataToAnotherHost(tagID string, unixTime int64, antennaNumber uint8) {
	conn, err := net.Dial("tcp", Config.PROXY_ADDRESS) 
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "%s, %d, %d\n", tagID, unixTime, antennaNumber)
}
