package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Birthday time.Time

func (b *Birthday) MarshalJSON() (data []byte, err error) {

	//str :=fmt.Sprintf("2006/01/02", b)
	now := time.Time(*b)
	str := now.Format("2006-01-02")
	data, err = json.Marshal(str)
	fmt.Printf("data:%s\n", str)
	return
}

func (b *Birthday) UnmarshalJSON(data []byte) (err error) {

	str := string(data)
	now, _ := time.Parse("2006/01/02", str)
	*b = Birthday(now)
	return
}

type Person struct {
	Age      int       `json:"age"`
	Id       int64     `json:"id,string"`
	Name     string    `json:"name_xx,omitempty"`
	Salary   float32   `json:"-"`
	Birthday *Birthday `json:"birthday"`
}

func test1() {
	var p = &Person{
		Age:    20,
		Id:     38888232322323222,
		Name:   "axx",
		Salary: 38822.2,
		//Birthday: time.Now(),
	}

	var birthday Birthday
	birthday = Birthday(time.Now())
	p.Birthday = &birthday

	data, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("marshal failed, err:%v\n", err)
		return
	}

	ioutil.WriteFile("./json.txt", data, 0777)

	data2, err := ioutil.ReadFile("./json.txt")
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return

	}

	var person2 Person
	json.Unmarshal(data2, &person2)
	fmt.Printf("person2:%#v\n", person2)
}

func main() {

	test1()

}
