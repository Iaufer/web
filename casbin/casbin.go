package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

func NewCasbin() (*casbin.Enforcer, error) {
	e, err := casbin.NewEnforcer("C:/Users/laufe/OneDrive/Рабочий стол/web/casbin/configs/model.conf", "C:/Users/laufe/OneDrive/Рабочий стол/web/casbin/configs/policy.csv")

	if err != nil {
		fmt.Println("err casbdi")
		return nil, err
	}

	if e == nil {
		fmt.Println("Casbin nil")
	} else {
		fmt.Println("Casbin не нил")

	}

	return e, nil

}

// func main() {
// 	e, err := casbin.NewEnforcer("configs/model.conf", "configs/policy.csv")

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	sub := "Bob"
// 	role := "reader"

// 	_, err = e.AddGroupingPolicy(sub, role)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	obj := "topic"
// 	act := "create"
// 	creator := "Bob"

// 	ok, err := e.Enforce(sub, obj, act)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if ok {
// 		fmt.Println(sub, " доступ есть", act, obj)
// 	} else {
// 		fmt.Println(sub, "доступа нет", act, obj)
// 	}

// }
