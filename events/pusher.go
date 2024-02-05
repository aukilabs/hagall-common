package events

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/aukilabs/hagall-common/errors"
	"github.com/aukilabs/hagall-common/logs"
)

const (
	DefaultFlushInterval = time.Second * 30
	DefaultBatchSize     = 50
	DefaultQueueSize     = 4080
)

var (
	// The timeout duration to push events.
	PushTimeout = time.Second * 10
)

// A Pusher that pushes events to a remote endpoint.
type Pusher struct {
	// The endpoint where events are sent.
	Endpoint string

	// The duration between each event flush.
	FlushInterval time.Duration

	// The maximum number of event sent at once. Default is 50.
	BatchSize int

	// The size of the queue where events are stored. Default is 4080.
	QueueSize int

	// The HTTP transport to send events. Default is http.DefaultTransport.
	Transport http.RoundTripper

	// The function to encode events. Default is json.Marshal.
	Encode func(any) ([]byte, error)

	initOnce  sync.Once
	startOnce sync.Once
	wg        sync.WaitGroup
	cancel    func()
	events    chan Event
}

// Start starts logging events asynchronously.
func (l *Pusher) Start() {
	l.initOnce.Do(l.init)

	l.startOnce.Do(func() {
		go l.start()
	})
}

// NewEvent logs the given event.
func (l *Pusher) NewEvent(e Event) {
	l.initOnce.Do(l.init)
	l.events <- e
}

// Close ensures all events had been flushed and release allocated resources.
func (l *Pusher) Close() {
	if l.cancel != nil {
		l.cancel()
	}

	l.wg.Wait()
}

func (l *Pusher) init() {
	if l.FlushInterval == 0 {
		l.FlushInterval = DefaultFlushInterval
	}

	if l.BatchSize == 0 {
		l.BatchSize = DefaultBatchSize
	}

	if l.QueueSize == 0 {
		l.QueueSize = DefaultQueueSize
	}

	if l.Transport == nil {
		l.Transport = http.DefaultTransport
	}

	if l.Encode == nil {
		l.Encode = json.Marshal
	}

	l.events = make(chan Event, l.QueueSize)
}

func (l *Pusher) start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l.cancel = cancel
	l.wg.Add(1)

	ticker := time.NewTicker(l.FlushInterval)
	defer ticker.Stop()

	batch := l.newBatch()

	for {
		select {
		case <-ctx.Done():
			l.flush(batch)
			l.wg.Done()
			return

		case e := <-l.events:
			if e != nil {
				batch = append(batch, e)
			}
			if len(batch) == l.BatchSize {
				go l.flush(batch)
				batch = l.newBatch()
			}

		case <-ticker.C:
			go l.flush(batch)
			batch = l.newBatch()
		}
	}
}

func (l *Pusher) newBatch() []Event {
	return make([]Event, 0, l.BatchSize)
}

func (l *Pusher) flush(batch []Event) {
	count := len(batch)
	if count == 0 {
		return
	}

	logs.WithTag("size", count).Debug("flushing event batch")

	if err := l.postEvents(batch); err != nil {
		logs.Error(errors.New("flushing event batch failed").
			WithTag("size", count).
			Wrap(err))
	}
}

func (l *Pusher) postEvents(batch []Event) error {
	if l.Endpoint == "" {
		return nil
	}

	body, err := l.Encode(eventPayload{Events: batch})
	if err != nil {
		return errors.New("encoding events failed").Wrap(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), PushTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, l.Endpoint, bytes.NewReader(body))
	if err != nil {
		return errors.New("creating request failed").Wrap(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := l.Transport.RoundTrip(req)
	if err != nil {
		return errors.New("request failed").Wrap(err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return errors.New("request failed").
			WithTag("status", res.Status).
			Wrap(err)
	}

	return nil
}

type eventPayload struct {
	Events []Event `json:"events"`
}
