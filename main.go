package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	hookHandler := HookHandler{loadConfig("config.json")}

	http.Handle("/", hookHandler)
	fmt.Printf("Starting server on port %d\n", hookHandler.Config.Port)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(hookHandler.Config.Port), nil))
}
