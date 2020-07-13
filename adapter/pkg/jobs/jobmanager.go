package jobs

import (
	"net/http"
	"strings"
	"time"

	"github.com/yuval-k/chainlink-twitter/adapter/pkg/twitter"
	"go.uber.org/zap"
)

const Interval = time.Minute / 2

type Job struct {
	Handle   string
	Text     string
	Callback func(done, approved bool, err error)
}

type JobManager struct {
	TwitterClient twitter.TwitterClient

	jobs      map[*Job]struct{}
	jobsToAdd chan *Job

	backoffDuration time.Duration

	logger *zap.SugaredLogger
}

func NewJobManager(logger *zap.SugaredLogger, twitterClient twitter.TwitterClient) *JobManager {
	return &JobManager{
		TwitterClient: twitterClient,
		jobs:          make(map[*Job]struct{}),
		jobsToAdd:     make(chan *Job, 100),
		logger:        logger,
	}
}

func (j *JobManager) AddJob() chan<- *Job {
	return j.jobsToAdd
}

func (j *JobManager) Run() {
	for {
		select {
		case <-j.backoff():
			// zbam
			for job := range j.jobs {
				j.logger.Infow("unqueueing job and retrying", "job", job)
				delete(j.jobs, job)
				j.runJobRatelimited(job)
				// do it one job at the time, so we also poll the job channel
				break
			}
		case job, ok := <-j.jobsToAdd:
			if !ok {
				return
			}
			j.logger.Infow("got new job", "job", job)
			j.runJobRatelimited(job)
		}
	}
}
func (j *JobManager) backoff() <-chan time.Time {
	var zero time.Duration
	if j.backoffDuration == zero {
		return nil
	}
	return time.After(j.backoffDuration)
}
func (j *JobManager) resetBackoff() {
	var zero time.Duration
	j.backoffDuration = zero
}
func (j *JobManager) increaseBackoff() {
	var zero time.Duration
	if j.backoffDuration == zero {
		j.backoffDuration = time.Second * 10
	}
	j.backoffDuration = j.backoffDuration * 2
}

func (j *JobManager) runJobRatelimited(job *Job) {
	j.logger.Infow("running job", "job", job)
	approved, resp, err := j.runJob(job)
	if err == nil {
		j.logger.Infow("running job success!", "approved", approved)
		j.resetBackoff()
		job.Callback(true, approved, nil)
		return
	}

	// if resp is nil, it might be a connection error, so let's retry.
	// error is anything but rate limit, respond with error.
	if resp != nil && resp.StatusCode != http.StatusTooManyRequests {
		j.resetBackoff()
		j.logger.Infow("running job error!", "err", err, "statuscode", resp.StatusCode)
		job.Callback(false, false, err)
		return
	}
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}
	j.logger.Infow("running job recoverable error", "err", err, "statuscode", statusCode)

	// if we are here we are either rate limited, or some connection error happened.
	// either way, do exp back-off
	// save the map

	j.logger.Infow("queueing job", "job", job)
	j.jobs[job] = struct{}{}
	j.increaseBackoff()
}

func (j *JobManager) runJob(job *Job) (bool, *http.Response, error) {
	tweets, resp, err := j.TwitterClient.GetTweetsFor(job.Handle)
	if err != nil {
		// TODO: return error; if it is rate limit then sleep / try again later
		return false, resp, err
	}

	for _, t := range tweets {
		if strings.Contains(t, job.Text) {
			// success!
			return true, nil, nil
		}
	}
	return false, nil, nil
}
