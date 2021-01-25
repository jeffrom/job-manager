package commands

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob"
	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type consumerOpts struct {
	concurrency     int
	claims          []string
	sleep           time.Duration
	failTimes       int
	shutdownTimeout time.Duration
	jitter          time.Duration
	doPanic         bool
	consumers       int
}

type consumerCmd struct {
	*cobra.Command
	opts *consumerOpts
}

func newConsumerCmd(cfg *client.Config) *consumerCmd {
	opts := &consumerOpts{}
	c := &consumerCmd{
		Command: &cobra.Command{
			Use:  "consumer QUEUE...",
			Args: cobra.MinimumNArgs(1),
			// Aliases: []string{"wrk"},
		},
		opts: opts,
	}

	flags := c.Command.Flags()
	flags.IntVarP(&opts.concurrency, "concurrency", "C", 1, "max concurrent jobs")
	flags.IntVar(&opts.consumers, "consumers", 1, "number of consumers (total is consumers * concurrency)")
	flags.IntVar(&opts.failTimes, "fail-times", 0, "number of failures before success")
	flags.StringArrayVarP(&opts.claims, "claim", "c", nil, "claims for this consumer")
	flags.DurationVar(&opts.sleep, "sleep", 0, "sleep before completion")
	flags.DurationVar(&opts.shutdownTimeout, "shutdown-timeout", 15*time.Second, "graceful shutdown period")
	flags.DurationVar(&opts.jitter, "jitter", 0, "jitter effect for sleep")
	flags.BoolVar(&opts.doPanic, "panic", false, "panic when failing")

	return c
}

func (c *consumerCmd) Cmd() *cobra.Command { return c.Command }
func (c *consumerCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	var claims label.Claims
	if len(c.opts.claims) > 0 {
		var err error
		claims, err = label.ParseClaims(c.opts.claims)
		if err != nil {
			return err
		}
	}
	cl := clientFromContext(ctx)
	queues := args

	var done context.CancelFunc
	ctx, done = context.WithCancel(ctx)
	defer done()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		done()
	}()
	log.Print("Starting consumer on queues: ", strings.Join(queues, ", "))

	wg := sync.WaitGroup{}
	for i := 0; i < c.opts.consumers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumer := mjob.NewConsumer(cl, &runner{opts: c.opts}, mjob.ConsumerWithConfig(mjob.ConsumerConfig{
				ShutdownTimeout: c.opts.shutdownTimeout,
				Concurrency:     c.opts.concurrency,
				DequeueOpts: client.DequeueOpts{
					Claims: claims,
					Queues: queues,
				},
			}))
			err := consumer.Run(ctx)
			if err != nil {
				log.Print("consumer", i, err)
			}
		}()
	}

	wg.Wait()
	return nil
}

type runner struct {
	opts *consumerOpts
}

func (r *runner) Run(ctx context.Context, job *resource.Job) (*resource.JobResult, error) {
	log.Printf("consumer executing on job %s", job.ID)
	if r.opts.sleep > 0 {
		d := r.opts.sleep
		if r.opts.jitter != 0 {
			neg := rand.Intn(1) == 1
			jit := time.Duration(rand.Intn(int(r.opts.jitter)))
			if neg {
				d += jit
			} else {
				d -= jit
			}
		}
		time.Sleep(d)
	}

	if ft := r.opts.failTimes; ft > 0 && job.Attempt <= ft {
		log.Printf("failing job %s on attempt %d", job.ID, job.Attempt)
		err := fmt.Errorf("failing job attempt %d", job.Attempt)
		if r.opts.doPanic {
			panic(err)
		}
		return nil, err
	}
	log.Printf("job %s complete", job.ID)
	return nil, nil
}
