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

type Car struct {
	Type     int
	Created  int64
	Finished int64
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

func testGo(id int, ch <-chan Car, wg *sync.WaitGroup, config []string) {
	defer wg.Done()

	for car := range ch {
		fmt.Printf("%d: obsluhuje", id)
		fmt.Println()
		test, _ := time.ParseDuration(config[1])
		test2, _ := time.ParseDuration(config[2])
		temp := time.Duration(randRange(int(test), int(test2)))
		time.Sleep(temp)
		car.Finished = time.Now().UnixMilli()
		fmt.Printf("%d: skoncila s casem: %s", id, temp)
		fmt.Println()
	}
}

func main() {
	config := loadConfig("config.txt")

	//s = stations
	s, _ := strconv.Atoi(config[3])
	var wg sync.WaitGroup
	carChan := make(chan Car)

	for i := 0; i < s; i++ {
		wg.Add(1)
		go testGo(i, carChan, &wg, config)
	}

	for c := 0; c < 5; c++ {
		t := time.Now().UnixMilli()
		fmt.Println(t)
		car := Car{1, t, t}
		carChan <- car
	}

	close(carChan)
	wg.Wait()
	save(config, "results.txt")

	fmt.Println("done")
}
