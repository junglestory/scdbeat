package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/junglestory/scdbeat/config"
)

type Scdbeat struct {
	beatConfig *config.Config
	done       chan struct{}
	period     time.Duration
	client     publisher.Client
}

// Creates beater
func New() *Scdbeat {
	return &Scdbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Scdbeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := b.RawConfig.Unpack(&bt.beatConfig)
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

func (bt *Scdbeat) Setup(b *beat.Beat) error {

	// Setting default period if not set
	if bt.beatConfig.Scdbeat.Period == "" {
		bt.beatConfig.Scdbeat.Period = "1s"
	}

	bt.client = b.Publisher.Connect()

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Scdbeat.Period)
	if err != nil {
		return err
	}

	return nil
}

func (bt *Scdbeat) Run(b *beat.Beat) error {
	logp.Info("scdbeat is running! Hit CTRL-C to stop it.")

	ticker := time.NewTicker(bt.period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Scdbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Scdbeat) Stop() {
	close(bt.done)
}
