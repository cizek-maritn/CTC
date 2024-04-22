package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

// Type - typ auta (gas, diesel, lpg, electric)
// Created - čas "vzniku", použito pro výpočet stráveného času ve frontě pumpy
// TimeStation - doba obsluhy na pumpě
// FinishedStationQ - čas dostání fronty u pumpy
// TimeReg - doba obsluhy u pokladny
// FinishedRegQ - čas dostání fronty u pokladny
type Car struct {
	Type    int
	Created int64

	TimeStation      time.Duration
	FinishedStationQ int64

	TimeReg      time.Duration
	FinishedRegQ int64
}

func loadConfig(fn string) []string {
	//nacitani config dat
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//zpracovani config dat
	scanner := bufio.NewScanner(file)
	var configData []string
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Split(line, ":")
		if len(temp[1]) > 0 {
			temp[1] = strings.TrimLeftFunc(temp[1], unicode.IsSpace)
			configData = append(configData, temp[1])
		}
	}

	return configData
}

func save(lines []string, path string) error {
	//ukladani dat
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func randRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func pumpAction(id int, ch <-chan Car, wg *sync.WaitGroup, config []string, r chan<- Car) {
	defer wg.Done()

	for car := range ch {
		confIndex := 3 * car.Type
		tmin, _ := time.ParseDuration(config[confIndex+1])
		tmax, _ := time.ParseDuration(config[confIndex+2])
		temp := time.Duration(randRange(int(tmin), int(tmax)))
		//obsluha
		car.TimeStation = temp
		car.FinishedStationQ = time.Now().UnixMilli()
		time.Sleep(temp)
		r <- car
	}
}

func registerAction(id int, ch chan Car, wg *sync.WaitGroup, config []string, r chan<- Car) {
	defer wg.Done()

	for car := range ch {
		tmin, _ := time.ParseDuration(config[16])
		tmax, _ := time.ParseDuration(config[17])
		temp := time.Duration(randRange(int(tmin), int(tmax)))
		//obsluha
		car.TimeReg = temp
		car.FinishedRegQ = time.Now().UnixMilli()
		time.Sleep(temp)
		r <- car
	}
	wg.Wait()
	close(ch)
}

func main() {
	config := loadConfig("config.txt")

	//repetitive, but im lazy
	gasStations, _ := strconv.Atoi(config[3])
	dieselStations, _ := strconv.Atoi(config[6])
	lpgStations, _ := strconv.Atoi(config[9])
	elecStations, _ := strconv.Atoi(config[12])
	registers, _ := strconv.Atoi(config[15])

	var wg sync.WaitGroup
	cars, _ := strconv.Atoi(config[0])

	//fronty aut
	gasChan := make(chan Car)
	dieselChan := make(chan Car)
	lpgChan := make(chan Car)
	elecChan := make(chan Car)
	regChan := make(chan Car)

	resultChan := make(chan Car, cars)

	for i := 0; i < gasStations; i++ {
		wg.Add(1)
		go pumpAction(i, gasChan, &wg, config, regChan)
	}

	for i := 0; i < dieselStations; i++ {
		wg.Add(1)
		go pumpAction(i, dieselChan, &wg, config, regChan)
	}

	for i := 0; i < lpgStations; i++ {
		wg.Add(1)
		go pumpAction(i, lpgChan, &wg, config, regChan)
	}

	for i := 0; i < elecStations; i++ {
		wg.Add(1)
		go pumpAction(i, elecChan, &wg, config, regChan)
	}

	for i := 0; i < registers; i++ {
		wg.Add(1)
		go registerAction(i, regChan, &wg, config, resultChan)
	}

	for c := 0; c < cars; c++ {
		//arrival time
		tmin, _ := time.ParseDuration(config[1])
		tmax, _ := time.ParseDuration(config[2])
		t := time.Now().UnixMilli()
		time.Sleep(time.Duration(randRange(int(tmin), int(tmax))))

		carType := randRange(1, 4)
		car := Car{carType, t, time.Second, t, time.Second, t}
		switch carType {
		case 1:
			gasChan <- car
		case 2:
			dieselChan <- car
		case 3:
			lpgChan <- car
		case 4:
			elecChan <- car
		}
	}

	close(gasChan)
	close(dieselChan)
	close(lpgChan)
	close(elecChan)

	//variable init
	//amount of cars of each type
	var counts [5]int
	//total time spent in each queue
	var Qtimes [5]time.Duration
	//avg queue time
	var AvgQ [5]time.Duration
	//max queue time
	var MaxQ [5]time.Duration

	//only way i made it work
	for i := 0; i < cars; i++ {
		select {
		case c := <-resultChan:
			index := c.Type - 1
			//increment counters
			counts[index]++
			//could be lazy and just say counts[4]=cars, but this can help catch possible errors
			counts[4]++

			//total queue times
			QtimeStation := (time.Duration((c.FinishedStationQ)-(c.Created)) * time.Millisecond)
			QtimeReg := (time.Duration((c.FinishedRegQ)-(c.FinishedStationQ)) * time.Millisecond)
			Qtimes[index] += QtimeStation
			Qtimes[4] += QtimeReg

			//maximum queue times
			if MaxQ[index] < QtimeStation {
				MaxQ[index] = QtimeStation
			}
			if MaxQ[4] < QtimeReg {
				MaxQ[4] = QtimeReg
			}
		}
	}

	strings := [5]string{"gas:", "diesel:", "lpg:", "electric:", "registers:"}
	var results []string
	results = append(results, "stations:")
	//average queue times
	for i := 0; i < 5; i++ {
		AvgQ[i] = Qtimes[i] / time.Duration(counts[i])
		results = append(results, strings[i])
		results = append(results, "total_cars: "+fmt.Sprint(counts[i]))
		results = append(results, "total_queue_time: "+time.Duration.String(Qtimes[i]))
		results = append(results, "average_queue_time: "+time.Duration.String(AvgQ[i]))
		results = append(results, "max_queue_time: "+time.Duration.String(MaxQ[i]))
	}

	save(results, "results.txt")

	fmt.Println("done")
}
