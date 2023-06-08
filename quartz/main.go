package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/reugn/go-quartz/quartz"
	"log"
	"math/rand"
	"time"
)

func main() {
	log.Println("Start demo")
	rand.Seed(time.Now().UnixMicro())

	ctx := context.Background()
	scheduler := quartz.NewStdSchedulerWithOptions(quartz.StdSchedulerOptions{
		WorkerLimit: 2,
	})
	scheduler.Start(ctx)

	jobs := [5]*quartz.FunctionJob[int]{}

	for i := 0; i < 5; i++ {
		index := i
		jobs[i] = quartz.NewFunctionJobWithDesc(
			fmt.Sprintf("random number job #%d", i),
			func(ctx context.Context) (int, error) {
				sleep := rand.Intn(10) + 1
				log.Println("try run a job", index, "sleep for", sleep, "seconds")
				time.Sleep(time.Second * time.Duration(sleep))

				if number := rand.Intn(100); number > 50 {
					return number, nil
				} else {
					return 0, errors.New("random is below 50")
				}
			},
		)

		scheduler.ScheduleJob(
			ctx,
			jobs[i],
			quartz.NewRunOnceTrigger(time.Second),
		)
	}

	scheduler.ScheduleJob(
		ctx,
		quartz.NewFunctionJobWithDesc("display jobs", func(ctx context.Context) (bool, error) {
			log.Println("display jobs ...")
			displayJobs(scheduler, jobs)
			return true, nil
		}),
		quartz.NewSimpleTrigger(time.Second*5),
	)

	scheduler.Wait(ctx)
}

func displayJobs(scheduler quartz.Scheduler, jobs [5]*quartz.FunctionJob[int]) {
	log.Println("Available jobs", scheduler.GetJobKeys())

	if 0 == len(scheduler.GetJobKeys()) {
		scheduler.Stop()
	}

	for _, job := range jobs {
		scheduledJob, err := scheduler.GetScheduledJob(job.Key())
		if err != nil {
			log.Println("job does not exist anymore in scheduler:", err)
			continue
		}

		funcJob := scheduledJob.Job.(*quartz.FunctionJob[int])
		log.Printf("Job %s status: %d\n", funcJob.Description(), funcJob.JobStatus)
	}
}
