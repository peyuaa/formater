package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"mohamed.attahri.com/jsonl"
)

const (
	input = "jokes.json"

	output = "jokes_formated.jsonl"
)

type jokes []joke

type joke struct {
	Context   string `json:"context"`
	Utterance string `json:"utterance"`
}

type formatedJokes []formatedJoke

type formatedJoke struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func main() {
	f, err := os.Open(input)
	if err != nil {
		log.Fatal("couldn't open file", err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal("couldn't close file", err)
		}
	}(f)

	j := make(jokes, 0, 114579)

	byteValue, err := io.ReadAll(f)
	if err != nil {
		log.Fatal("couldn't read bytes from file", err)
	}

	err = json.Unmarshal(byteValue, &j)
	if err != nil {
		log.Fatal("couldn't parse json", err)
	}

	formattedJokes := make(formatedJokes, 0, len(j))
	for _, joke := range j {
		// delete all \n from input and output

		inputStr := strings.Replace(strings.Replace(joke.Context, "\n", " ", -1), ",", ".", -1)
		outputStr := strings.Replace(strings.Replace(joke.Utterance, "\n", " ", -1), ",", ".", -1)

		if len(inputStr) > 80 || len(outputStr) > 80 {
			continue
		}

		fj := formatedJoke{
			Input:  inputStr,
			Output: outputStr,
		}

		formattedJokes = append(formattedJokes, fj)
	}

	outFile, err := os.Create(output)
	if err != nil {
		log.Fatal("couldn't create output file", err)
	}

	defer func(f *os.File) {
		err := outFile.Close()
		if err != nil {
			log.Fatal("couldn't close file", err)
		}
	}(outFile)

	writer := jsonl.NewWriter[formatedJoke](outFile)
	if _, err := writer.Write(formattedJokes...); err != nil {
		log.Fatal(err)
	}

	err = outFile.Sync()
	if err != nil {
		log.Fatal("couldn't sync file", err)
	}
}
