// Copyright 2016 Iron.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package drivers

import (
	"testing"
	"time"
)

func TestAverage(t *testing.T) {
	start := time.Date(2016, 8, 11, 0, 0, 0, 0, time.UTC)
	stats := make([]Stat, 10)
	for i := 0; i < len(stats); i++ {
		stats[i] = Stat{
			Timestamp: start.Add(time.Duration(i) * time.Minute),
			Metrics:   map[string]uint64{"x": uint64(i)},
		}
	}

	res, ok := average(stats)
	if !ok {
		t.Error("Expected good record")
	}

	expectedV := uint64(4)
	if v, ok := res.Metrics["x"]; !ok || v != expectedV {
		t.Error("Actual average didn't match expected", "actual", v, "expected", expectedV)
	}

	expectedT := time.Unix(1470873870, 0)
	if res.Timestamp != expectedT {
		t.Error("Actual average didn't match expected", "actual", res.Timestamp, "expected", expectedT)
	}
}

func TestDecimate(t *testing.T) {
	start := time.Now()
	stats := make([]Stat, 480)
	for i := range stats {
		stats[i] = Stat{
			Timestamp: start.Add(time.Duration(i) * time.Second),
			Metrics:   map[string]uint64{"x": uint64(i)},
		}
		//		t.Log(stats[i])
	}

	stats = Decimate(240, stats)
	if len(stats) != 240 {
		t.Error("decimate function bad", len(stats))
	}

	//for i := range stats {
	//t.Log(stats[i])
	//}

	stats = make([]Stat, 700)
	for i := range stats {
		stats[i] = Stat{
			Timestamp: start.Add(time.Duration(i) * time.Second),
			Metrics:   map[string]uint64{"x": uint64(i)},
		}
	}
	stats = Decimate(240, stats)
	if len(stats) != 240 {
		t.Error("decimate function bad", len(stats))
	}

	stats = make([]Stat, 300)
	for i := range stats {
		stats[i] = Stat{
			Timestamp: start.Add(time.Duration(i) * time.Second),
			Metrics:   map[string]uint64{"x": uint64(i)},
		}
	}
	stats = Decimate(240, stats)
	if len(stats) != 240 {
		t.Error("decimate function bad", len(stats))
	}

	stats = make([]Stat, 300)
	for i := range stats {
		if i == 150 {
			// leave 1 large gap
			start = start.Add(20 * time.Minute)
		}
		stats[i] = Stat{
			Timestamp: start.Add(time.Duration(i) * time.Second),
			Metrics:   map[string]uint64{"x": uint64(i)},
		}
	}
	stats = Decimate(240, stats)
	if len(stats) != 49 {
		t.Error("decimate function bad", len(stats))
	}
}

func TestParseImage(t *testing.T) {
	cases := map[string][]string{
		"iron/hello":                                        {"", "iron/hello", "latest"},
		"iron/hello:v1":                                     {"", "iron/hello", "v1"},
		"my.registry/hello":                                 {"my.registry", "hello", "latest"},
		"my.registry/hello:v1":                              {"my.registry", "hello", "v1"},
		"mongo":                                             {"", "library/mongo", "latest"},
		"mongo:v1":                                          {"", "library/mongo", "v1"},
		"quay.com/iron/hello":                               {"quay.com", "iron/hello", "latest"},
		"quay.com:8080/iron/hello:v2":                       {"quay.com:8080", "iron/hello", "v2"},
		"localhost.localdomain:5000/samalba/hipache:latest": {"localhost.localdomain:5000", "samalba/hipache", "latest"},
	}

	for in, out := range cases {
		reg, repo, tag := ParseImage(in)
		if reg != out[0] || repo != out[1] || tag != out[2] {
			t.Errorf("Test input %q wasn't parsed as expected. Expected %q, got %q", in, out, []string{reg, repo, tag})
		}
	}
}
