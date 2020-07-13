package main

import (
	"net/http"

	"github.com/yuval-k/chainlink-twitter/adapter/pkg/adapter"
	"github.com/yuval-k/chainlink-twitter/adapter/pkg/jobs"
	"github.com/yuval-k/chainlink-twitter/adapter/pkg/twitter"

	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	defer sugar.Sync()
	sugar.Infow("Starting adapter")
	tc := twitter.NewTwitterClientFromEnv(sugar)
	jm := jobs.NewJobManager(sugar, tc)
	adapter := adapter.NewFromEnv(sugar, jm.AddJob())
	go jm.Run()
	http.ListenAndServe(":8080", adapter)
}
