/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	buildTime   string = "now"
	buildGitTag string = "dev"
)
var isReady = false

var listen *string = flag.String("listen", ":12380", "listen for health")
var url *string = flag.String("url", "http://127.0.0.1:3000", "target url")
var http_timeout *time.Duration = flag.Duration("http.timeout", 1*time.Second, "http.timeout")
var try_timeout *time.Duration = flag.Duration("try.timeout", 1*time.Second, "time before next reguest")
var wait_httpStatus *int = flag.Int("wait.http.status", 200, "wait for http status")
var wait_success_probes *int = flag.Int("wait.success_probes", 3, "max success probes")
var resultFile *string = flag.String("result.file", "", "print ok to file")
var wait_httpStatusCount int = 0
var host *string = flag.String("host", "", "change host")

func main() {
	log.Printf("starting %s-%s...\n", buildGitTag, buildTime)
	flag.Parse()

	go check()

	http.HandleFunc("/ready", ready)
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(*listen, nil)

	if err != nil {
		log.Panic(err)
	}
}

func check() {
	client := http.Client{
		Timeout: *http_timeout,
	}

	for {
		req, err := http.NewRequest("GET", *url, nil)
		if err != nil {
			log.Println(err)
		}

		if len(*host) > 0 {
			req.Host = *host
		}

		resp, err := client.Do(req)

		if err != nil {
			log.Println(err)
		}
		if resp != nil {
			log.Printf("resp.StatusCode=%d\n", resp.StatusCode)
			if resp.StatusCode == *wait_httpStatus {
				wait_httpStatusCount = wait_httpStatusCount + 1
				log.Printf("wait_httpStatusCount=%d\n", wait_httpStatusCount)
			} else {
				wait_httpStatusCount = 0
			}
		}
		if wait_httpStatusCount >= *wait_success_probes {
			isReady = true
			log.Println("condition completed")
			break
		}
		time.Sleep(*try_timeout)
	}
	if len(*resultFile) > 0 {
		err := ioutil.WriteFile(*resultFile, []byte("ok"), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func ready(w http.ResponseWriter, r *http.Request) {
	if !isReady {
		http.Error(w, "url not ready", http.StatusInternalServerError)
		return
	}
	_, err := w.Write([]byte("ok"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func healthz(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
