/*
Copyright 2019 vChain, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package suite

import (
	"runtime"
	"strconv"

	"github.com/codenotary/immudb/pkg/bm"
	"github.com/codenotary/immudb/pkg/db"
)

var maxProcs int
var Concurrency = runtime.NumCPU()
var V = []byte{0, 1, 3, 4, 5, 6, 7}

var FunctionBenchmarks = []bm.Bm{
	{
		CreateTopic: true,
		Name:        "sequential write (fine tuned / experimental)",
		Concurrency: 1_000_000, // Concurrency,
		Iterations:  1_000_000,
		Before: func(bm *bm.Bm) {
			maxProcs = runtime.GOMAXPROCS(Concurrency)
		},
		After: func(bm *bm.Bm) {
			runtime.GOMAXPROCS(maxProcs)
		},
		Work: func(bm *bm.Bm, start int, end int) {
			for i := start; i < end; i++ {
				key := []byte(strconv.FormatUint(uint64(i), 10))
				bm.Topic.Set(key, V)
			}
		},
	},
	{
		CreateTopic: true,
		Name:        "sequential write (baseline)",
		Concurrency: Concurrency,
		Iterations:  1_000_000,
		Work: func(bm *bm.Bm, start int, end int) {
			for i := start; i < end; i++ {
				key := []byte(strconv.FormatUint(uint64(i), 10))
				_ = bm.Topic.Set(key, V)
			}
		},
	},
	{
		CreateTopic: true,
		Name:        "sequential write (concurrency++)",
		Concurrency: Concurrency * 8,
		Iterations:  1_000_000,
		Work: func(bm *bm.Bm, start int, end int) {
			for i := start; i < end; i++ {
				key := []byte(strconv.FormatUint(uint64(i), 10))
				_ = bm.Topic.Set(key, V)
			}
		},
	},
	{
		CreateTopic: true,
		Name:        "sequential write (GOMAXPROCS=128)",
		Concurrency: Concurrency,
		Iterations:  1_000_000,
		Before: func(bm *bm.Bm) {
			maxProcs = runtime.GOMAXPROCS(128)
		},
		After: func(bm *bm.Bm) {
			runtime.GOMAXPROCS(maxProcs)
		},
		Work: func(bm *bm.Bm, start int, end int) {
			for i := start; i < end; i++ {
				key := []byte(strconv.FormatUint(uint64(i), 10))
				_ = bm.Topic.Set(key, V)
			}
		},
	},
	{
		CreateTopic: true,
		Name:        "batch write",
		Concurrency: Concurrency,
		Iterations:  1_000_000,
		Work: func(bm *bm.Bm, start int, end int) {
			var kvPairs []db.KVPair
			for i := start; i < end; i++ {
				kvPairs = append(kvPairs, db.KVPair{
					Key:   []byte(strconv.FormatUint(uint64(i), 10)),
					Value: V,
				})
			}
			_ = bm.Topic.SetBatch(kvPairs)
		},
	},
}