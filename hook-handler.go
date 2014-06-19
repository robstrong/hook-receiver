package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

type HookHandler struct {
	Config Config
}

func (h HookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ipParts := strings.Split(r.RemoteAddr, ":")
	ip := net.ParseIP(ipParts[0])

	if h.isValidIp(ip) {
		payload, err := formatPayload(r)
		if err != nil {
			fmt.Printf("Invalid payload\n")
		} else {
			h.runEvents(payload)
			fmt.Fprint(w, "WebHook Received")
		}
	} else {
		fmt.Fprint(w, "Reeeeeejected!!!")
	}
}
func (h HookHandler) runEvents(payload Payload) {
	//check if payload matches any of the rules
Rule:
	for _, rule := range h.Config.Rules {
		for _, criteria := range rule.Criteria {
			//check that types match
			if criteria.Type != "" && payload.Type != criteria.Type {
				continue Rule
			}
			//check that owners match
			if criteria.Owner != "" && payload.Owner != criteria.Owner {
				continue Rule
			}
			//check that repo names match
			if criteria.Repository != "" && payload.Repository != criteria.Repository {
				continue Rule
			}

			//we have a matching rule, run the command
			output, err := runCommand(rule.Command)
			if err != nil {
				fmt.Printf("Command Error:\n    %s\n", err)
			}
			//format the output
			outputStr := string(output)
			if strings.HasSuffix(outputStr, "\n") {
				outputStr = outputStr[:len(outputStr)-1]
			}
			outputStr = "    " + strings.Replace(string(outputStr), "\n", "\n    ", -1) + "\n"
			fmt.Printf("Command output:\n%s", outputStr)
		}
	}
}

func runCommand(cmd string) (output []byte, err error) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	output, err = exec.Command(head, parts...).CombinedOutput()
	return
}

func parseCIDRs(cidrs []string) []*net.IPNet {
	if len(cidrs) == 0 {
		log.Fatal("No CIDRs specified")
	}
	cidrNet := make([]*net.IPNet, 0)
	for _, cidr := range cidrs {
		_, netCidr, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatal(err)
		}
		cidrNet = append(cidrNet, netCidr)
	}
	return cidrNet
}

func (h HookHandler) isValidIp(ip net.IP) bool {
	for _, cidr := range h.Config.CIDRs {
		if cidr.Contains(ip) {
			return true
		}
		fmt.Printf("IP %s is not in Github CIDR: %s\n", ip, cidr.String())
	}
	return false
}
