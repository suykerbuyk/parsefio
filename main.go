package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/wI2L/jettison"
)

func describe(f interface{}) {
	val := reflect.TypeOf(f).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeF := val.Field(i)
		fieldName := typeF.Name
		jsonTag := typeF.Tag.Get("json")
		temp := val.Field(i)
		fmt.Printf("Field: %s  jsonTag: %s  ", fieldName, jsonTag)
		fmt.Println(typeF)
		if typeF.Type.Kind() == reflect.Struct {
			fmt.Println("Found embedded struct", fieldName)
			describe(temp)
		}
	}
}
func examiner(t reflect.Type, depth int) {
	fmt.Println(strings.Repeat("\t", depth), "Type is", t.Name(), "and kind is", t.Kind())
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		fmt.Println(strings.Repeat("\t", depth+1), "Contained type:")
		examiner(t.Elem(), depth+1)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Println(strings.Repeat("\t", depth+1), "Field", i+1, "name is", f.Name, "type is", f.Type.Name(), "and kind is", f.Type.Kind())
			if f.Tag != "" {
				fmt.Println(strings.Repeat("\t", depth+2), "Tag is", f.Tag)
			}
			if f.Type.Kind() == reflect.Struct {
				k := reflect.ValueOf(f)
				examiner(reflect.TypeOf(k), depth+1)
			}
		}
	}
}

func printValue(prefix string, v reflect.Value, visited map[interface{}]bool) {

	//fmt.Printf("%s: ", v.Type())

	// Drill down through pointers and interfaces to get a value we can print.
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.Kind() == reflect.Ptr {
			// Check for recursive data
			if visited[v.Interface()] {
				fmt.Println("visted")
				return
			}
			visited[v.Interface()] = true
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		//fmt.Printf("%d elements\n", v.Len())
		for i := 0; i < v.Len(); i++ {
			fmt.Printf("%s%d.", prefix, i)
			printValue(prefix+"++", v.Index(i), visited)
		}
	case reflect.Struct:
		t := v.Type() // use type to get number and names of fields
		//fmt.Printf("%d fields\n", t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fmt.Printf("%s%s.", prefix, t.Field(i).Name)
			//printValue(prefix+"   ", v.Field(i), visited)
			printValue(prefix+t.Field(i).Name+".", v.Field(i), visited)
		}
	case reflect.Invalid:
		fmt.Printf("nil\n")
	default:
		fmt.Printf("%v\n", v.Interface())
	}
}
func main() {
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
	printValue("FioDATA.", reflect.ValueOf(&result), make(map[interface{}]bool))
	fmt.Println("===============================")
	examiner(reflect.TypeOf(result), 1)
	fmt.Println("===============================")
	describe(&result)
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
