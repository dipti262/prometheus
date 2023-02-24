// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
        "fmt"
        "log"
        "net/http"

        "github.com/prometheus/common/model"

        "github.com/prometheus/prometheus/storage/remote"
)

var (
        username = "abc"
        password = "123"
)


func main() {
        http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
                u, p, ok := r.BasicAuth()
                if !ok {
                        fmt.Println("Error parsing basic auth")
                        w.WriteHeader(401)
                        return
                        }
                if u != username {
                        fmt.Printf("Username provided is correct: %s\n", u)
                        w.WriteHeader(401)
                        return
                        }
                if p != password {
                        fmt.Printf("Password provided is correct: %s\n", u)
                        w.WriteHeader(401)
                        return
                        }
                fmt.Printf("Username: %s\n", u)
                fmt.Printf("Password: %s\n", p)
                req, err := remote.DecodeWriteRequest(r.Body)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusBadRequest)
                        return
                }

                for _, ts := range req.Timeseries {
                        m := make(model.Metric, len(ts.Labels))
                        for _, l := range ts.Labels {
                                m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
                        }
                        fmt.Println(m)

                        for _, s := range ts.Samples {
                                fmt.Printf("\tSample:  %f %d\n", s.Value, s.Timestamp)
                        }

                        for _, e := range ts.Exemplars {
                                m := make(model.Metric, len(e.Labels))
                                for _, l := range e.Labels {
                                        m[model.LabelName(l.Name)] = model.LabelValue(l.Value)
                                }
                                fmt.Printf("\tExemplar:  %+v %f %d\n", m, e.Value, e.Timestamp)
                        }

                        for _, hp := range ts.Histograms {
                                h := remote.HistogramProtoToHistogram(hp)
                                fmt.Printf("\tHistogram:  %s\n", h.String())
                        }
                }


      })
      //http.Handle("/example", handler)
        log.Fatal(http.ListenAndServe("<FIXME_localhostIP:port>", nil))
}
