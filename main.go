package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wI2L/jettison"
	"io/ioutil"
	"os"
	"reflect"
)

func describe(f interface{}) {
	val := reflect.TypeOf(f).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeF := val.Field(i)
		fieldName := typeF.Name
		jsonTag := typeF.Tag.Get("json")
		fmt.Printf("Field: %s  jsonTag: %s\n", fieldName, jsonTag)
		fmt.Println(typeF)
	}
}

func main() {
	textPtr := flag.String("text", "my_text", "Something texty")
	flag.Parse()
	fmt.Println("test flag: ", *textPtr)
	fmt.Println(flag.Args())
	jsonFile, err := os.Open("librbd-lr02u27-J24-write-32k.fio.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	var result FioData
	json.Unmarshal([]byte(byteValue), &result)
	describe(&result)
	fmt.Println("===============================")
	describe(&result.GlobalOptions)
	fmt.Println("===============================")
	fmt.Printf("FioVersion %s\n", result.FioVersion)
	fmt.Printf("GlobalOptions.IoEngine %s\n", result.GlobalOptions.Ioengine)
	fmt.Printf("Job elements %d\n", len(result.Jobs))
	for idx, job := range result.Jobs {
		fmt.Printf("Idx: %d Jobname:%s\n", idx, job.Jobname)
		describe(&job.Read)
	}
	fmt.Println("===============================")
	jsonOut, err := jettison.Marshal(result)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(jsonOut)
}
