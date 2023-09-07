package main

/*
accept command flag -key to capture string
accept command flag -o to capture string
*/

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wejick/gchain/callback"
	"github.com/wejick/gchain/eval"
	"github.com/wejick/gchain/model"
	_openai "github.com/wejick/gchain/model/openAI"
)

var (
	o string
)

func init() {
	flag.StringVar(&o, "o", "", "output file")
	flag.Parse()
}

type Test struct {
	Name        string
	Evaluator   string
	Input       string
	Expectation string
	Reason      string
	Result      bool
}

func main() {
	// open csv file from o
	// read csv file
	// for each row
	// put to array of Test
	tests, err := readCSV(o)
	if err != nil {
		fmt.Println(err)
		return
	}

	var authToken = os.Getenv("OPENAI_API_KEY")
	chatModel := _openai.NewOpenAIChatModel(authToken, "", "", _openai.GPT3Dot5Turbo0301, callback.NewManager(), false)

	testRunner(tests, chatModel)
}

func testRunner(test []Test, llmModel model.LLMModel) {
	jsonEvaluator := eval.NewValidJson()
	var testResult []Test
	for _, t := range test {
		var evaluator eval.Evaluator
		if t.Evaluator == "valid_json" {
			evaluator = jsonEvaluator
		} else if t.Evaluator == "correctness" {
			evaluator = eval.NewCorrectnessEval(llmModel, t.Expectation)
		}
		var errReason error
		t.Result, errReason = evaluator.Evaluate(t.Input)
		if errReason != nil {
			t.Reason = errReason.Error()
		}
		testResult = append(testResult, t)
	}

	fmt.Println("Test Result")
	for _, t := range testResult {
		fmt.Printf("Test Name: %s\n", t.Name)
		fmt.Printf("Test Input: %s\n", t.Input)
		fmt.Printf("Test Expectation: %s\n", t.Expectation)
		fmt.Printf("Test Result: %v\n", t.Result)
		fmt.Printf("Test Reason: %s\n", t.Reason)
		fmt.Println("========================================")
	}
}

func readCSV(filename string) ([]Test, error) {
	var tests []Test

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read the first line to skip it
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		test := Test{
			Name:        record[0],
			Evaluator:   record[1],
			Input:       record[2],
			Expectation: record[3],
		}
		tests = append(tests, test)
	}

	return tests, nil
}
