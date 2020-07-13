package adapter

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/yuval-k/chainlink-twitter/adapter/pkg/jobs"
)

type BridgeData struct {
	Handle string `json:"handle`
	Text   string `json:"text`
}

// copied from and modified bridgeOutgoing from: github.com/smartcontractkit/chainlink/core/adapters/bridge.go as it is private
type BridgeInput struct {
	JobRunID    string       `json:"id"`
	Data        *BridgeData  `json:"data"`
	Meta        *models.JSON `json:"meta,omitempty"`
	ResponseURL string       `json:"responseURL,omitempty"`
}

type BridgeOutputData struct {
	Done     bool `json:"done"`
	Approved bool `json:"approved"`
}

type Bridge struct {
	fromChainlinkToken string

	nodeAddress      *url.URL
	toChainlinkToken string

	//	send jobs here
	jobManager chan<- *jobs.Job
}

func NewFromEnv(jobManager chan<- *jobs.Job) *Bridge {
	// http://host:port
	node := os.Getenv("CHAINLINK_NODE")
	nodeurl, err := url.Parse(node)
	if err != nil {
		panic(err)
	}
	return &Bridge{
		fromChainlinkToken: os.Getenv("OUTGOING_TOKEN"),
		nodeAddress:        nodeurl,
		toChainlinkToken:   os.Getenv("INCOMING_TOKEN"),
		jobManager:         jobManager,
	}
}

var _ http.Handler = new(Bridge)

func (b *Bridge) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	authToken := utils.StripBearer(r.Header.Get("Authorization"))
	if subtle.ConstantTimeCompare([]byte(b.fromChainlinkToken), []byte(authToken)) != 1 {
		// TODO: don't panic
		panic("unauthorized")
	}

	var input BridgeInput
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&input); err != nil {
		panic(err)
	}
	out, err := b.Run(&input)
	if err != nil {
		panic(err)
	}

	e := json.NewEncoder(rw)
	e.Encode(out)

}

func (b *Bridge) Run(input *BridgeInput) (*models.BridgeRunResult, error) {
	// might need to do retries, so return the job as pending,
	// and move to processing that will handler that

	respUrl := input.ResponseURL
	if respUrl == "" {
		runpatch, _ := url.Parse("/v2/runs/" + input.JobRunID)
		finalurl := b.nodeAddress.ResolveReference(runpatch)
		respUrl = finalurl.String()
	}

	job := &jobs.Job{
		Handle: input.Data.Handle,
		Text:   input.Data.Text,
		Callback: func(done, approved bool, err error) {
			var brr models.BridgeRunResult
			if err != nil {
				brr.Status = models.RunStatusErrored
			} else {
				var bod BridgeOutputData
				bod.Done = done
				bod.Approved = approved

				brr.Status = models.RunStatusCompleted

				jsn, err := json.Marshal(&bod)
				if err != nil {
					// TODO: log instead of panic
					panic(err)
				}
				brr.Data.UnmarshalJSON(jsn)
			}
			err = b.Patch(respUrl, &brr)
			if err != nil {
				// TODO: do some retries / or log instead of panic
				panic(err)
			}
		},
	}

	// TODO: should i select here?
	b.jobManager <- job

	return &models.BridgeRunResult{
		ExternalPending: true,
		Status:          models.RunStatusInProgress,
	}, nil
}

func (b *Bridge) Patch(respUrl string, result *models.BridgeRunResult) error {

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(result)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", respUrl, &buffer)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+b.toChainlinkToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed " + resp.Status)
	}
	return nil
}
