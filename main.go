package main

import (
	"encoding/json"
	"fmt"
	"qles/qles"
)

// fix dot convertion shit
func main() {
	query := "SELECT a, b, c FROM Temp WHERE a IN (1, 2, 3) ORDER BY a desc, b asc"
	boolq, err := qles.BuildES(query, map[string]string{"a": "p1", "b": "path"})
	if err != nil {
		fmt.Println(err)
	}
	j, _ := json.Marshal(boolq)
	fmt.Println(string(j))
}
