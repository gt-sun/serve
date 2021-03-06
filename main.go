// Copyright 2014 Karan Misra.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var version = "0.2.3"

var (
	port        = flag.Int("p", 5000, "port to serve on")
	prefix      = flag.String("x", "/", "prefix to serve under")
	showVersion = flag.Bool("v", false, "show version info")
	openBrowser = flag.Bool("o", false, "open the url")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println("serve version", version)
		os.Exit(0)
	}

	var dir string
	// Get the dir to serve
	if flag.NArg() < 1 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Please provide the dir to serve as the last argument. A simple . will also do")
			os.Exit(1)
		}
		dir = cwd
	}
	dir = flag.Arg(0)
	portStr := fmt.Sprintf(":%v", *port)
	if !strings.HasPrefix(*prefix, "/") {
		*prefix = "/" + *prefix
	}
	if !strings.HasSuffix(*prefix, "/") {
		*prefix = *prefix + "/"
	}

	uri := fmt.Sprintf("http://localhost:%v%v", *port, *prefix)

	fmt.Printf("Service traffic from %v under port %v with prefix %v\n", dir, *port, *prefix)
	fmt.Printf("Or simply put, just open %v to get rocking!\n", uri)

	go func() {
		if *openBrowser {
			success := waitForWebserver()
			if !success {
				// We have waited too long for the webserver to start; bail.
				fmt.Fprintf(os.Stderr, "The webserver did not start within the required time. Cannot open the browser for you\n")
				return
			}
			fmt.Printf("Opening your browser to %v\n", uri)
			cmd := exec.Command("open", uri)
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Could not open the url in your browser\n%v\n", err)
			}
		}
	}()

	http.Handle(*prefix, http.StripPrefix(*prefix, http.FileServer(http.Dir(dir))))
	if err := http.ListenAndServe(portStr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error while starting the web server\n%v\n", err)
		os.Exit(1)
	}
}

func waitForWebserver() bool {
	timeout := time.After(1 * time.Second)
	connStr := fmt.Sprintf("127.0.0.1:%v", *port)
	for {
		select {
		case <-timeout:
			return false
		default:
			conn, err := net.DialTimeout("tcp", connStr, 50*time.Millisecond)
			if err != nil {
				continue
			}
			conn.Close()
			return true
		}
	}
}
