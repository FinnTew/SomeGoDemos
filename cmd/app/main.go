package main

import (
	"fmt"
	"some-go-demos/pkgs/my"
	"strconv"
)

func main() {
	mySlice := my.NewSlice[string](2)
	for i := 0; i < 10; i++ {
		mySlice.Append(strconv.Itoa(i))
		fmt.Println("Len: ", mySlice.Len(), " Cap: ", mySlice.Cap())
	}
	fmt.Println(mySlice.Get(5))
	fmt.Println(mySlice.Slice(3, 6).Array())
	mySlice.Del(5)
	fmt.Println(mySlice.Slice(3, 6).Array())
	mySlice.Set(5, "hello")
	fmt.Println(mySlice.Slice(3, 6).Array())
	fmt.Println("Len: ", mySlice.Len(), " Cap: ", mySlice.Cap())
}
