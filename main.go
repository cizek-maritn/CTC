package main

import (
	"bufio"
	"fmt"
	"log"
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

func main() {
	config := loadConfig("config.txt")

	save(config, "results.txt")

	fmt.Println("done")
}
