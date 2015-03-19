package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"text/template"
)

type HookHandler struct {
	Config Config
}

func (h HookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.isValidIp(r.RemoteAddr) {
		fmt.Fprint(w, "Rejected!!!")
		return
	}

	payload, err := parsePayload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go h.handlePayload(payload)

	fmt.Fprint(w, "WebHook Received")
}

func (h HookHandler) handlePayload(payload Payload) {
	//check if payload matches any of the rules
Rule:
	for _, rule := range h.Config.Rules {
		for _, criteria := range rule.Criteria {

			if !payload.IsMatch(criteria) {
				continue Rule
			}

			//we have a matching rule, run the command
			output, err := runCommand(rule.Command, payload)
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

func runCommand(cmd string, payload Payload) (output []byte, err error) {

	parsed, err := parseCommand(cmd, payload)
	if err != nil {
		return
	}

	parts := strings.Fields(parsed)
	head := parts[0]
	parts = parts[1:len(parts)]

	output, err = exec.Command(head, parts...).CombinedOutput()
	return
}

func parseCommand(cmd string, payload Payload) (string, error) {
	tmpl, err := template.New(cmd).Parse(cmd)
	if err != nil {
		return "", err
	}
	out := bytes.NewBuffer(make([]byte, 0))
	err = tmpl.Execute(out, payload)
	if err != nil {
		return "", err
	}
	return out.String(), nil
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

func (h HookHandler) isValidIp(addr string) bool {

	ipParts := strings.Split(addr, ":")
	ip := net.ParseIP(ipParts[0])

	for _, cidr := range h.Config.CIDRs {
		if cidr.Contains(ip) {
			return true
		}
		fmt.Printf("IP %s is not in Github CIDR: %s\n", ip, cidr.String())
	}
	return false
}
