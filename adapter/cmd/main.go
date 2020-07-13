package main

import (
	"net/http"

	"github.com/yuval-k/chainlink-twitter/adapter/pkg/adapter"
	"github.com/yuval-k/chainlink-twitter/adapter/pkg/jobs"
	"github.com/yuval-k/chainlink-twitter/adapter/pkg/twitter"
)

func main() {
	tc := twitter.NewTwitterClientFromEnv()
	jm := jobs.NewJobManager(tc)
	adapter := adapter.NewFromEnv(jm.AddJob())
	go jm.Run()
	http.ListenAndServe(":8080", adapter)
}
