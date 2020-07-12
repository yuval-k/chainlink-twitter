package jobs

import (
	"strings"
	"time"

	"github.com/yuval-k/chainlink-twitter/adapter/pkg/twitter"
)

const Interval = time.Minute / 2

type Job struct {
	Handle   string
	Text     string
	NotAfter time.Time
	Callback func(success bool)
}

type JobManager struct {
	TwitterClient twitter.TwitterClient

	jobs      map[*Job]struct{}
	jobsToAdd chan *Job
}

func NewJobManager(twitterClient twitter.TwitterClient) *JobManager {
	return &JobManager{
		TwitterClient: twitterClient,
		jobsToAdd:     make(chan *Job, 10),
	}
}

func (j *JobManager) AddJob() chan<- *Job {
	return j.jobsToAdd
}

func (j *JobManager) Run() {
	for {
		j.RunOnce()
		select {
		case <-time.After(Interval):
		case job, ok := <-j.jobsToAdd:
			if !ok {
				return
			}
			j.jobs[job] = struct{}{}
		}
	}

}
func (j *JobManager) RunOnce() {
	for job := range j.jobs {
		tweets, err := j.TwitterClient.GetTweetsFor(job.Handle)
		if err != nil {
			// TODO: log error
			continue
		}
		for _, t := range tweets {
			if strings.Contains(t, job.Text) {
				// success!
				job.Callback(true)
				delete(j.jobs, job)
			}
		}
		if job.NotAfter.After(time.Now()) {
			job.Callback(false)
			delete(j.jobs, job)
		}
	}
}
