package main

import (
	"fmt"
	"testing"
)

func TestReWriteResult(t *testing.T) {
	result := `{
  "data": {
    "queryProduct": []
  },
  "extensions": {
    "tracing": {
      "version": 1,
      "startTime": "2023-04-11T02:53:54.820966972Z",
      "endTime": "2023-04-11T02:53:54.84397305Z",
      "duration": 23006056,
      "execution": {
        "resolvers": [
          {
            "path": [
              "queryProduct"
            ],
            "parentType": "Query",
            "fieldName": "queryProduct",
            "returnType": "[Product]",
            "startOffset": 176707,
            "duration": 22816420,
            "dgraph": [
              {
                "label": "query",
                "startOffset": 273462,
                "duration": 22716903
              }
            ]
          }
        ]
      }
    }
  }
}`
	writeResult, err := ReWriteResult([]byte(result))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(writeResult)
}
