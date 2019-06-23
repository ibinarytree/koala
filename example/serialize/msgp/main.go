package main

import (
	"fmt"
	"io/ioutil"
)

//go:generate msgp
type Person struct {
	Age    int     `json:"age"`
	Id     int64   `json:"id,string"`
	Name   string  `json:"name_xx,omitempty"`
	Salary float32 `json:"-"`
}

func test1() {

	var p = &Person{
		Age:    20,
		Id:     38888232322323222,
		Name:   "axx",
		Salary: 38822.2,
	}

	data, err := p.MarshalMsg(nil)
	if err != nil {
		fmt.Printf("marshal failed, err:%v\n", err)
		return
	}

	ioutil.WriteFile("./msg.txt", data, 0777)

	data2, err := ioutil.ReadFile("./msg.txt")
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return

	}

	var person2 Person
	person2.UnmarshalMsg(data2)
	fmt.Printf("person2:%#v\n", person2)

}

func main() {

	test1()

}
