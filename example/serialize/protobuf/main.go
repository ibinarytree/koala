package main

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/ibinarytree/micro-kit/serialize/protobuf/address"
)

func main() {
	var person address.Person
	person.Id = 3988222
	person.Name = "hua"

	var phone address.Phone
	phone.Number = "13872832832"
	person.Phones = append(person.Phones, &phone)

	data, err := proto.Marshal(&person)
	if err != nil {
		fmt.Printf("marshal failed, err:%v\n", err)
		return
	}

	ioutil.WriteFile("./proto.dat", data, 0777)

	data2, err := ioutil.ReadFile("./proto.dat")
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return
	}

	var person2 address.Person
	proto.Unmarshal(data2, &person2)

	fmt.Printf("unmarshal person:%#v\n", person2)
}
