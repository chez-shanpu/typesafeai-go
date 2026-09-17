// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of typesafeai-go

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/chez-shanpu/typesafeai-go"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <number>", os.Args[0])
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid number %q: %v", os.Args[1], err)
	}

	apiKey := os.Getenv("TYPESAFE_API_KEY")
	if apiKey == "" {
		log.Fatal("TYPESAFE_API_KEY is not set")
	}

	trueDesc := "the number is divisible by 2 (even)"
	falseDesc := "the number is not divisible by 2 (odd)"

	const questionKey = "num"
	questions := map[string]typesafeai.Question{
		questionKey: &typesafeai.NoulQuestion{
			Instructions: fmt.Sprintf("Is the following number even?: %d", n),
			Criteria: &typesafeai.NoulCriteria{
				True:  &trueDesc,
				False: &falseDesc,
			},
		},
	}

	client := typesafeai.NewSystemOneClient(apiKey, http.DefaultClient)

	resp, err := client.Do(context.Background(), &typesafeai.SystemOneRequest{
		State:     strconv.Itoa(n),
		Model:     typesafeai.ModelJevLatest,
		Questions: questions,
	})
	if err != nil {
		log.Fatalf("systemone request failed: %v", err)
	}

	ans, ok := resp.Answers[questionKey]
	if !ok {
		log.Fatalf("answer for %s not found", questionKey)
	}

	noul, ok := ans.(*typesafeai.NoulAnswer)
	if !ok {
		log.Fatalf("answer for %s is not a NoulAnswer: %T", questionKey, ans)
	}

	parity := "odd"
	if noul.Noul >= 0.5 {
		parity = "even"
	}
	fmt.Printf("%d is %s (noul=%.4f)\n", n, parity, noul.Noul)
}
