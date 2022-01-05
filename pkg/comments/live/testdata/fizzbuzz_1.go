package main

import "fmt"

func main() {
	for i := 1; i <= 50; i++ {
		fizzbuzz(i)
	}
}

// this is a comment on the function!
func fizzbuzz(i int) {
	fizz := "fizz"
	buzz := "buzz"

	// below is the fizzbuzzing, it's pretty cool!
	// :-)
	if i%3 == 0 && i%5 == 0 {
		fmt.Println(i, fizz+buzz)
	} else if i%3 == 0 {
		fmt.Println(i, fizz)
	} else if i%5 == 0 {
		fmt.Println(i, buzz)
	} else {
		fmt.Println(i)
	}
}
