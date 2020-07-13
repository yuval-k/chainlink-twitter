package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/yuval-k/chainlink-twitter/adapter/pkg/jobs"
	"go.uber.org/zap"
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
	logger     *zap.SugaredLogger
}

func NewFromEnv(logger *zap.SugaredLogger, jobManager chan<- *jobs.Job) *Bridge {
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

		logger: logger,
	}
}

var _ http.Handler = new(Bridge)

func (b *Bridge) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	/*
		authToken := utils.StripBearer(r.Header.Get("Authorization"))
		if subtle.ConstantTimeCompare([]byte(b.fromChainlinkToken), []byte(authToken)) != 1 {
			b.logger.Error("not authorized!")
			http.Error(rw, "not authorized", http.StatusUnauthorized)
			return
		}
	*/
	var input BridgeInput
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&input); err != nil {
		b.logger.Error("can't decode input!")
		http.Error(rw, "can't decode input!", http.StatusUnprocessableEntity)
		return
	}
	out, err := b.Run(&input)
	if err != nil {
		b.logger.Errorw("can't process job run!", "error", err)
		http.Error(rw, "can't process job run!", http.StatusInternalServerError)
		return
	}
	b.logger.Infow("job run success", "out", out)

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
	b.logger.Infow("job run start", "respUrl", respUrl)

	job := &jobs.Job{
		Handle: input.Data.Handle,
		Text:   input.Data.Text,
		Callback: func(done, approved bool, err error) {
			b.logger.Infow("job callback", "done", done, "approved", approved, "err", err)
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
				// TODO: do some retries? or is that the node's responsibility?/s
				b.logger.Errorw("can't patch job!", "err", err)
			}
		},
	}

	// TODO: should i select here?
	b.jobManager <- job
	b.logger.Infow("job sent to manager", "job", job)

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
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errordata, _ := ioutil.ReadAll(resp.Body)
		b.logger.Errorw("can't patch job!", "errordata", string(errordata), "statusCode", resp.StatusCode)
		return fmt.Errorf("request failed " + resp.Status)
	}
	return nil
}
