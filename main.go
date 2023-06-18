package main

import (
	"encoding/json"
	"fmt"
	"qles/qles"
)

// fix dot convertion shit
func main() {

	query := "SELECT * FROM Temp ORDER BY a ASC"
	ast, _ := qles.BuildSQL(query)
	reverse := qles.ReverseNot(ast)
	boolq := qles.BuildESQuery(reverse, map[string]string{"a": "p1", "b": "path"})
	j, _ := json.Marshal(boolq)
	fmt.Println(string(j))
}
