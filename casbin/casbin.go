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

// 	// sub := "bob"
// 	// obj := "topic"
// 	// act := "read"

// 	// if res, err := e.Enforce(sub, obj, act); err != nil {
// 	// 	log.Fatal(err)
// 	// } else if res {
// 	// 	fmt.Println("Succes")
// 	// } else {
// 	// 	fmt.Println("Unsucces")
// 	// }

// 	_, err = e.AddRoleForUser("Artyom", "reader")

// 	if err != nil {
// 		return
// 	}

// 	if res, err := e.Enforce("Artyom", "topic", "create"); err != nil {
// 		log.Fatal(err)
// 	} else if res {
// 		fmt.Println("Succes")
// 	} else {
// 		fmt.Println("Unsucces")
// 	}

// }
