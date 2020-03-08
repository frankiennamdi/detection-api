package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/frankiennamdi/detection-api/test"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/frankiennamdi/detection-api/models"
	"github.com/frankiennamdi/detection-api/support"
	"github.com/google/uuid"
)

func main() {
	numEvents := flag.Int("num", 3000, "number of events to generate")
	flag.Parse()
	log.Printf(support.Info, "Generating Events")

	users := []string{"bob", "mark", "johnny", "mary", "kevin", "mike", "case"}

	IPlist := []string{"206.81.252.6", "24.242.71.20", "91.207.175.104"}

	timeChanges := []int{-1, -2, -3, -4, -5, 1, 2, 3, 4, 5}

	startTime := int64(1514764800)

	wg := sync.WaitGroup{}

	for index := range users {
		wg.Add(1)

		go func(userForEvents string) {
			defer wg.Done()

			rand.Seed(time.Now().Unix())

			for i := 1; i < *numEvents; i++ {
				IP := IPlist[randomNum(0, len(IPlist)-1)]

				timeChange := timeChanges[randomNum(0, len(timeChanges)-1)]

				eventInfo := models.EventInfo{
					UUID:      uuid.New().String(),
					Username:  userForEvents,
					Timestamp: test.AddTime(startTime, timeChange*randomNum(1, 20), time.Hour),
					IP:        IP,
				}

				log.Printf(support.Info, eventInfo)
				body, err := json.Marshal(eventInfo)

				if err != nil {
					log.Panicf(support.Fatal, err)
				}

				resp, err := http.Post("http://localhost:3000/api/events", "application/json",
					bytes.NewBuffer(body))

				if err != nil {
					log.Fatalln(err)
				}

				var result models.SuspiciousTravelResult

				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					log.Panicf(support.Fatal, err)
				}

				data, err := json.MarshalIndent(&result, "", "    ")

				if err != nil {
					log.Panicf(support.Fatal, err)
				}

				log.Printf(support.Info, string(data))
			}
		}(users[index])
	}

	wg.Wait()
}

func randomNum(min, max int) int {
	return rand.Intn(max-min) + min
}
