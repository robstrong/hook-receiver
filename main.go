package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

var config Config
var CIDRs []*net.IPNet

func main() {
    config := loadConfig()

    //get valid CIDRs from Github
    if len(config.CIDROverride) != 0 {
        CIDRs = parseCIDRs(config.CIDROverride)
    } else {
        CIDRs = getGithubCIDRs()
    }
    fmt.Println("CIDRs: ", CIDRs)

    http.HandleFunc("/", requestHandler)
    fmt.Printf("Starting server on port %d\n", config.Port)

    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
    ipParts := strings.Split(r.RemoteAddr, ":")
    ip := net.ParseIP(ipParts[0])

    if isValidIp(ip) {
        payload, err := formatPayload(r)
        if err != nil {
            fmt.Printf("Invalid payload\n")
        } else {
            runEvents(payload)
            fmt.Fprint(w, "WebHook Received")
        }
    } else {
        fmt.Fprint(w, "Reeeeeejected!!!")
    }
}

type Config struct {
    Port         int
    Rules        []Rule
    CIDROverride []string `json:"cidr_override"`
}

type Rule struct {
    Command  string
    Criteria []Criteria
}

type Criteria struct {
    Type       string
    Owner      string
    Repository string
}

type Payload struct {
    Type       string
    Owner      string
    Repository string
}

func formatPayload(req *http.Request) (payload Payload, err error) {
    var jsonBody struct {
        Repository struct {
            Owner struct {
                Login string //these two hold the same data, but the push event format is different
                Name  string //from all the other events
            }
            Name string
        }
    }
    defer req.Body.Close()
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        return
    }

    err = json.Unmarshal(body, &jsonBody)
    if err != nil {
        return
    }

    payload.Type = req.Header.Get("X-Github-Event")
    payload.Owner = jsonBody.Repository.Owner.Name
    if payload.Owner == "" {
        payload.Owner = jsonBody.Repository.Owner.Login
    }
    payload.Repository = jsonBody.Repository.Name
    json, err := json.Marshal(payload)
    fmt.Println("Payload: ", string(json))

    return
}

func runEvents(payload Payload) {
    //check if payload matches any of the rules
Rule:
    for _, rule := range config.Rules {
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
            } else {
                //format the output
                output := string(output)
                if strings.HasSuffix(output, "\n") {
                    output = output[:len(output)-1]
                }
                output = "    " + strings.Replace(string(output), "\n", "\n    ", -1) + "\n"
                fmt.Printf("Command output:\n%s", output)
            }
        }
    }
}

func runCommand(cmd string) (output []byte, err error) {
    parts := strings.Fields(cmd)
    head := parts[0]
    parts = parts[1:len(parts)]

    output, err = exec.Command(head, parts...).Output()
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

func loadConfig() Config {
    file, err := os.Open("config.json")
    if err != nil {
        log.Fatal("Error loading config file: ", err)
    }
    decoder := json.NewDecoder(file)
    config = Config{}
    err = decoder.Decode(&config)
    if err != nil {
        log.Fatal("Invalid config file: ", err)
    }

    return config
}

func getGithubCIDRs() []*net.IPNet {
    //request CIDRs from Github
    resp, err := http.Get("https://api.github.com/meta")
    if err != nil {
        log.Fatal("Could not load Github CIDRs")
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    var data struct {
        Hooks []string //we only really care about the Hooks
    }
    json.Unmarshal(body, &data)

    //convert the response into net.IPNet slice
    cidrs := parseCIDRs(data.Hooks)

    return cidrs
}

func isValidIp(ip net.IP) bool {
    for _, cidr := range CIDRs {
        if cidr.Contains(ip) {
            return true
        }
        fmt.Printf("IP %s is not in Github CIDR: %s\n", ip, cidr.String())
    }
    return false
}
