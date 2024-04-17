package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
)

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

func main() {
	config := loadConfig("config.txt")

	fmt.Println(randRange(1, 5))
	fmt.Println((randRange(1, 5)))
	fmt.Println((randRange(1, 5)))

	save(config, "results.txt")

	fmt.Println("done")
}
