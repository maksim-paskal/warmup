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
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("starting...")
	url := flag.String("url", "http://127.0.0.1:3000", "target url")
	http_timeout := flag.String("http.timeout", "1s", "http.timeout")
	try_timeout := flag.String("try.timeout", "1s", "time before next reguest")

	flag.Parse()

	http_timeout_d, err := time.ParseDuration(*http_timeout)
	if err != nil {
		log.Fatal(err)
	}

	try_timeout_d, err := time.ParseDuration(*try_timeout)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: http_timeout_d,
	}

	for {
		req, err := http.NewRequest("GET", *url, nil)
		if err != nil {
			log.Println(err)
		}

		resp, err := client.Do(req)

		if err != nil {
			log.Println(err)
		}
		if resp != nil {
			log.Printf("resp.StatusCode=%d\n", resp.StatusCode)
			if resp.StatusCode == 404 {
				log.Println("condition completed")
				break
			}
		}
		time.Sleep(try_timeout_d)
	}

}
