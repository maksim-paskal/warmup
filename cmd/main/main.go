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
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	warmupHeader            = "X-Warmup-Request"
	defaultMaxSuccessProbes = 3
)

var (
	gitVersion          = "dev"
	isReady             = false
	listen              = flag.String("listen", ":12380", "listen for health")
	url                 = flag.String("url", "http://127.0.0.1:3000", "target url")
	httpTimeout         = flag.Duration("http.timeout", 1*time.Second, "http.timeout")
	tryTimeout          = flag.Duration("try.timeout", 1*time.Second, "time before next reguest")
	waitHTTPStatus      = flag.Int("wait.http.status", http.StatusOK, "wait for http status")
	waitSuccessProbes   = flag.Int("wait.success_probes", defaultMaxSuccessProbes, "max success probes")
	resultFile          = flag.String("result.file", "", "print ok to file")
	waitHTTPStatusCount = 0
	host                = flag.String("host", "", "change host")
	insertWarmupHeader  = flag.Bool("insert-warmup-header", true, "inserts to request header "+warmupHeader)
	headers             = flag.String("headers", "", "add headers to request - example X-Test1=1,X-Test2=2")
)

func main() {
	log.Printf("starting %s...\n", gitVersion)
	flag.Parse()

	go check()

	http.HandleFunc("/ready", ready)
	http.HandleFunc("/healthz", healthz)

	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Panic(err)
	}
}

func check() { //nolint:funlen,cyclop
	client := http.Client{
		Timeout: *httpTimeout,
	}

	ctx := context.Background()

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", *url, nil)
		if err != nil {
			log.Println(err)

			continue
		}

		if len(*host) > 0 {
			req.Host = *host
		}

		if *insertWarmupHeader {
			req.Header.Add(warmupHeader, "true")
		}

		if len(*headers) > 0 {
			for _, hValue := range strings.Split(*headers, ",") {
				k := strings.Split(hValue, "=")
				if len(k) == 2 { //nolint:gomnd
					req.Header.Add(k[0], k[1])
				} else {
					log.Printf("WARN header %s - invalid\n", hValue)
				}
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)

			continue
		}

		resp.Body.Close()

		log.Printf("resp.StatusCode=%d\n", resp.StatusCode)

		if resp.StatusCode == *waitHTTPStatus {
			waitHTTPStatusCount++
			log.Printf("waitHTTPStatusCount=%d\n", waitHTTPStatusCount)
		} else {
			waitHTTPStatusCount = 0
		}

		if waitHTTPStatusCount >= *waitSuccessProbes {
			isReady = true

			log.Println("condition completed")

			break
		}

		time.Sleep(*tryTimeout)
	}

	if len(*resultFile) > 0 {
		err := ioutil.WriteFile(*resultFile, []byte("ok"), 0o644) //nolint:gosec,gomnd
		if err != nil {
			log.Println(err)
		}
	}
}

func ready(w http.ResponseWriter, r *http.Request) {
	if !isReady {
		http.Error(w, "url not ready", http.StatusInternalServerError)

		return
	}

	if _, err := w.Write([]byte("ok")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("ok")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
