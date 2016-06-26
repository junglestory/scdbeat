package beater

import (
	"fmt"
	"time"
	"io/ioutil"
	"bufio"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/junglestory/scdbeat/config"
	"strings"
	"os"
)

type Scdbeat struct {
	beatConfig *config.Config
	done       chan struct{}
	period     time.Duration
	client     publisher.Client
	path         string    //root 디렉토리
	lasIndexTime time.Time //가장 마지막에 검색한 시간
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

	// path 값을 lsbeat.yml 에 설정된 값으로 설정.
	if bt.beatConfig.Scdbeat.Path == "" {
		bt.beatConfig.Scdbeat.Path = "."
	}
	bt.path = bt.beatConfig.Scdbeat.Path

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

		readScd(bt.path, bt, b, counter) // scdbeat  데이터 송신
		bt.lasIndexTime = time.Now()     // 송신 후 타임스탬프 기록.

		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
/*
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
		}
		bt.client.PublishEvent(event)
		*/
		logp.Info("Event sent")
		counter++
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readScd(dirFile string, bt *Scdbeat, b *beat.Beat, counter int) {
	files, _ := ioutil.ReadDir(dirFile)

	for _, f := range files {
		pos := strings.Index(f.Name(), "-I-C.SCD")

		if pos > 0 {
			fmt.Println(f.Name())

			file, err := os.Open(dirFile + "/" + f.Name())
			check(err)

			defer file.Close()

			var scdMap map[string]string
			scdMap = make(map[string]string)

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				pos := strings.Index(scanner.Text(), ">")

				if pos > 0 {
					field := scanner.Text()[1:pos]
					text := scanner.Text()[pos+1:]

					if field == "DOCID" {
						if len(scdMap) > 0 {
							sendScd(scdMap, b, bt, counter);
							scdMap = make(map[string]string)
						}
					}
					scdMap[field] = text
				}
			}

			sendScd(scdMap, b, bt, counter);
		}
	}
}

func sendScd (scdMap map[string]string, b *beat.Beat, bt *Scdbeat, counter int) {
	event := common.MapStr{
		"DOCID":     scdMap["DOCID"],
		"date": scdMap["Date"],
		"title":    scdMap["TITLE"],
		"content":    scdMap["CONTENT"],
		"write": scdMap["WRITER"],
		"filelink": scdMap["FILELINK"],
		"link": scdMap["LINK"],
		"attachcon": scdMap["ATTACHCON"],
		"attachext": scdMap["ATTACHEXT"],
		"attachname": scdMap["ATTACHNAME"],
		"alias": scdMap["ALIAS"],
	}

	bt.client.PublishEvent(event) //elasticsearch 로 데이터 전송.
}

func (bt *Scdbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Scdbeat) Stop() {
	close(bt.done)
}
