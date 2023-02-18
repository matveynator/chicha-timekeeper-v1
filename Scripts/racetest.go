package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"chicha/Packages/Config"
)

const (
	competitors          = 10 // riders
	results              = 10 // results from one rider
	laps                 = 10 // laps
	minimal_lap_time_sec = 120
)

var appAddr = Config.APP_ANTENNA_LISTENER_IP

func main() {
	var racers []Racer
	for i := 0; i < competitors; i++ {
		racers = append(racers, newRacer(i))
	}

	// вейтгруппа для всех участников соревнования
	wg := &sync.WaitGroup{}

	// запускаем всех участников в своей собственной горутине
	for _, racer := range racers {
		wg.Add(1)
		go racer.Run(wg)
	}

	// ожидаем заверешение гонки от всех гонщиков
	wg.Wait()
}

func newRacer(number int) Racer {
	return Racer{
		tag:   fmt.Sprintf("TESTRIDER000%d%d00", rand.Intn(10000), rand.Intn(10000)),
		racer: number,
	}
}

type Racer struct {
	tag      string
	racer    int
	unixtime string
	antenna  string
}

func (r Racer) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= laps; i++ {
		// любой участник кроме первого будет стартовать со случайной задержкой
		if r.racer != 0 {
			time.Sleep(time.Second * time.Duration(rand.Intn(11)))
		}

		iterations := rand.Intn(results)
		for j := 0; j < iterations; j++ {
			r.unixtime = fmt.Sprintf("%d", time.Now().UnixMilli())
			r.antenna = fmt.Sprintf("%d", rand.Intn(5))
			r.writeRfid()
			// задержка между разными источниками единовременными, как я понял
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
		}

		// после итераций ждем пока участник доедет снова до точки с антеннами
		time.Sleep(time.Second * minimal_lap_time_sec)
	}
}

func (r Racer) writeRfid() {
	dial, err := net.Dial("tcp", appAddr)
	if err != nil {
		log.Fatal("failed to dial to chicha:", err)
		return
	}
	defer dial.Close()

	_, err = fmt.Fprintf(os.Stdout, "%s%d, %s, %s\n", r.tag, r.racer, r.unixtime, r.antenna)
	_, err = fmt.Fprintf(dial, "%s%d, %s, %s\n", r.tag, r.racer, r.unixtime, r.antenna)
	if err != nil {
		log.Println("failed to write rfid:", err)
	}
}
