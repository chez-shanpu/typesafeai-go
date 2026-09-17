# typesafeai

A Go client for the [TypeSafe AI](https://typesafe.ai) API.
It is developed as a repository independent from the official SDKs (Python / JavaScript-TypeScript).

You can send three types of typed questions — Noul (yes/no probability), Choice (options), and Score (ordered levels) —
and receive answers in the corresponding types.

## Installation

```sh
go get github.com/chez-shanpu/typesafeai-go
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chez-shanpu/typesafeai-go"
)

func main() {
	apiKey := os.Getenv("TYPESAFE_API_KEY")
	if apiKey == "" {
		log.Fatal("TYPESAFE_API_KEY is not set")
	}

	trueDesc := "the number is divisible by 2 (even)"
	falseDesc := "the number is not divisible by 2 (odd)"

	const questionKey = "num"
	questions := map[string]typesafeai.Question{
		questionKey: &typesafeai.NoulQuestion{
			Instructions: "Is the following number even?: 42",
			Criteria: &typesafeai.NoulCriteria{
				True:  &trueDesc,
				False: &falseDesc,
			},
		},
	}

	client := typesafeai.NewSystemOneClient(apiKey, http.DefaultClient)

	resp, err := client.Do(context.Background(), &typesafeai.SystemOneRequest{
		State:     42,
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

	fmt.Printf("noul=%.4f\n", noul.Noul)
}
```

See [`examples/noul.go`](./examples/noul.go) for a more detailed example.
