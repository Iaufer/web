package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

// &&

func NewCasbin() (*casbin.Enforcer, error) {
	e, err := casbin.NewEnforcer("C:/Users/laufe/OneDrive/Рабочий стол/web/casbin/configs/model.conf", "C:/Users/laufe/OneDrive/Рабочий стол/web/casbin/configs/policy.csv")
	if err != nil {
		fmt.Println("err casbdi: ", err)
		return nil, err
	}

	if e == nil {
		fmt.Println("Casbin nil")
	} else {
		fmt.Println("Casbin не нил")

	}

	// e.AddFunction("timeFunc", func(args ...interface{}) (interface{}, error) {
	// 	fmt.Println("AddFunsadction")

	// 	u := args[0].(string)

	// 	fmt.Println("UUUU: ", u)

	// 	if u == "1" {
	// 		return true, nil
	// 	} else {
	// 		return false, nil
	// 	}
	// 	// fmt.Println("AddFunction", u)

	// 	// loc, err := time.LoadLocation("Local")

	// 	// if err != nil {
	// 	// 	fmt.Println("Locat: ", err)
	// 	// }

	// 	// lastUpdated, err := time.ParseInLocation(l, u, loc)

	// 	// fmt.Println("AddFunction", lastUpdated)

	// 	// if err != nil {
	// 	// 	fmt.Println("ADd err: ", err)
	// 	// }
	// 	// fmt.Println("ДАШЕЛ")
	// 	// fmt.Println("lastUpdated) >= 2*time.Minute: ", time.Since(lastUpdated) >= 2*time.Minute)
	// 	// fmt.Println("time.Since(lastUpdated): ", time.Since(lastUpdated))
	// 	// fmt.Println("2*time.Minute: ", 2*time.Minute)
	// 	// return time.Since(lastUpdated) >= 2*time.Minute, nil
	// })

	return e, nil

}

// func main() {
// 	_, err := casbin.NewEnforcer("configs/model.conf", "configs/policy.csv")

// 	l := "2006-01-02 15:04:05.999999"
// 	loc, err := time.LoadLocation("Local")
// 	if err != nil {
// 		fmt.Println("Locat: ", err)
// 	}
// 	lastUpdated, err := time.ParseInLocation(l, "2024-12-09 04:39:37.392205", loc)
// 	b := time.Since(lastUpdated) >= 2*time.Minute
// 	_b := ""
// 	if b {
// 		_b = "1"
// 	} else {
// 		_b = "0"
// 	}

// 	fmt.Println("BBBB: ", _b)

// }
