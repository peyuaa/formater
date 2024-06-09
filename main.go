package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

const (
	input = "jokes.json"

	output = "jokes_formated.json"
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
		fj := formatedJoke{
			Input:  joke.Context,
			Output: joke.Utterance,
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

	outBytes, err := json.Marshal(formattedJokes)
	if err != nil {
		log.Fatal("couldn't marshal json", err)
	}

	_, err = outFile.Write(outBytes)
	if err != nil {
		log.Fatal("couldn't write json", err)
	}

	err = outFile.Sync()
	if err != nil {
		log.Fatal("couldn't sync output file", err)
	}
}