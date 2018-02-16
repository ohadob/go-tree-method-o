package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"GO-TREE-METHOD-O/RF"
)

var forest *RF.Forest

//buli

func main() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.HandleFunc("/updateTree", handler)
	mux.HandleFunc("/predicate", predicate)

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(":" + port, mux))
}

func predicate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	var arr []interface{}
	json.Unmarshal(body, &arr)
	output := forest.Predicate(arr)
	bytes, _ := json.Marshal(output)
	w.Write(bytes)
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	var arr [][]string
	json.Unmarshal(body, &arr)
	inputs := make([][]interface{}, 0)
	targets := make([]string, 0)
	for _, tup := range arr {
		pattern := tup[:len(tup)-1]
		target := tup[len(tup)-1]
		X := make([]interface{}, 0)
		for _, x := range pattern {
			X = append(X, x)
		}
		inputs = append(inputs, X)

		targets = append(targets, target)
	}
	train_inputs := make([][]interface{}, 0)

	train_targets := make([]string, 0)

	test_inputs := make([][]interface{}, 0)
	test_targets := make([]string, 0)

	for i, x := range inputs {
		if i%2 == 1 {
			test_inputs = append(test_inputs, x)
		} else {
			train_inputs = append(train_inputs, x)
		}
	}

	for i, y := range targets {
		if i%2 == 1 {
			test_targets = append(test_targets, y)
		} else {
			train_targets = append(train_targets, y)
		}
	}

	forest = RF.BuildForest(inputs, targets, 10, 500, len(train_inputs[0])) //100 trees
	RF.DumpForest(forest, "test")
	fjson, _ := json.Marshal(forest)
	log.Printf("forest - [%s]", string(fjson))
}
