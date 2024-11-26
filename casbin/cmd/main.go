package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

func main() {

	e, err := casbin.NewEnforcer("../configs/model.conf", "../configs/policy.csv")

	if err != nil {
		log.Fatal(err)
	}

	sub := "bob"
	obj := "topic"
	act := "read"

	if res, err := e.Enforce(sub, obj, act); err != nil {
		log.Fatal(err)
	} else if res {
		fmt.Println("Succes")
	} else {
		fmt.Println("Unsucces")
	}

}
