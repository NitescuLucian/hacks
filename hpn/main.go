package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

//Struct to hold the Nmap output
type NmapRun struct {
	XMLName  xml.Name `xml:"nmaprun"`
	Text     string   `xml:",chardata"`
	Scanner  string   `xml:"scanner,attr"`
	Start    string   `xml:"start,attr"`
	StartStr string   `xml:"startstr,attr"`
	Version  string   `xml:"version,attr"`
	Scaninfo struct {
		Text        string `xml:",chardata"`
		Type        string `xml:"type,attr"`
		Protocol    string `xml:"protocol,attr"`
		NumServices string `xml:"numservices,attr"`
		Services    string `xml:"services,attr"`
	} `xml:"scaninfo"`
	Verbose struct {
		Text  string `xml:",chardata"`
		Level string `xml:"level,attr"`
	} `xml:"verbose"`
	Debugging struct {
		Text  string `xml:",chardata"`
		Level string `xml:"level,attr"`
	} `xml:"debugging"`
	Host []struct {
		Text      string `xml:",chardata"`
		Starttime string `xml:"starttime,attr"`
		Endtime   string `xml:"endtime,attr"`
		Status    struct {
			Text      string `xml:",chardata"`
			State     string `xml:"state,attr"`
			Reason    string `xml:"reason,attr"`
			ReasonTtl string `xml:"reason_ttl,attr"`
		} `xml:"status"`
		Address struct {
			Text     string `xml:",chardata"`
			Addr     string `xml:"addr,attr"`
			AddrType string `xml:"addrtype,attr"`
			Vendor   string `xml:"vendor,attr"`
		} `xml:"address"`
		Hostnames struct {
			Text     string `xml:",chardata"`
			Hostname struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				Type string `xml:"type,attr"`
			} `xml:"hostname"`
		} `xml:"hostnames"`
		Ports struct {
			Text       string `xml:",chardata"`
			ExtraPorts []struct {
				Text      string `xml:",chardata"`
				Count     string `xml:"count,attr"`
				State     string `xml:"state,attr"`
				Reason    string `xml:"reason,attr"`
				ReasonTtl string `xml:"reason_ttl,attr"`
			} `xml:"extraports"`
			Port []struct {
				Text      string `xml:",chardata"`
				Protocol  string `xml:"protocol,attr"`
				PortId    string `xml:"portid,attr"`
				State     string `xml:"state,attr"`
				Reason    string `xml:"reason,attr"`
				ReasonTtl string `xml:"reason_ttl,attr"`
			} `xml:"port"`
		} `xml:"ports"`
		Os struct {
			Text    string `xml:",chardata"`
			OsMatch []struct {
				Text     string `xml:",chardata"`
				Name     string `xml:"name,attr"`
				Accuracy string `xml:"accuracy,attr"`
				Line     string `xml:"line,attr"`
				OsClass  []struct {
					Text     string `xml:",chardata"`
					Type     string `xml:"type,attr"`
					Vendor   string `xml:"vendor,attr"`
					OsGen    string `xml:"osgen,attr"`
					Accuracy string `xml:"accuracy,attr"`
				} `xml:"osclass"`
			} `xml:"osmatch"`
		} `xml:"os"`
		Times struct {
			Text   string `xml:",chardata"`
			Srtt   string `xml:"srtt,attr"`
			RttVar string `xml:"rttvar,attr"`
			To     string `xml:"to,attr"`
		} `xml:"times"`
	} `xml:"host"`
	Runstats struct {
		Text     string `xml:",chardata"`
		Finished struct {
			Text     string `xml:",chardata"`
			Time     string `xml:"time,attr"`
			Timestr  string `xml:"timestr,attr"`
			Elapsed  string `xml:"elapsed,attr"`
			Summary  string `xml:"summary,attr"`
			Exit     string `xml:"exit,attr"`
			Errormsg string `xml:"errormsg,attr"`
			Hosts    string `xml:"hosts,attr"`
		} `xml:"finished"`
		Hosts struct {
			Text  string `xml:",chardata"`
			Down  string `xml:"down,attr"`
			Total string `xml:"total,attr"`
			Up    string `xml:"up,attr"`
		} `xml:"hosts"`
	} `xml:"runstats"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file name as an argument to the script.")
		return
	}
	fileName := os.Args[1]
	//Open the xml file
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	//Unmarshal the xml data into the NmapRun struct
	var n NmapRun
	err = xml.Unmarshal(file, &n)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	//Loop through the hosts to print out the ip address and port
	for _, host := range n.Host {
		for _, port := range host.Ports.Port {
			fmt.Printf("%s:%s\n", host.Address.Addr, port.PortId)
		}
	}
}
