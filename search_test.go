package main

import (
	"testing"
	"fmt"
)

func TestSumUp(t *testing.T) {
	store = append(store, LogLine{Id:1,Time:1,Body:"First"})
	store = append(store, LogLine{Id:2,Time:2,Body:"Second"})
	store = append(store, LogLine{Id:3,Time:2,Body:"Third"})
	store = append(store, LogLine{Id:4,Time:2,Body:"Fourth"})
	store = append(store, LogLine{Id:5,Time:5,Body:"Fifth"})
	ret,_ := search("", 3,0)
	for _, value := range ret {
		fmt.Println(value.Id, value.getTime(), value.Body)
	}
	fmt.Println(len(ret))

	ret,_ = search("", 3,1)

	for _, value := range ret {
		fmt.Println(value.Id, value.getTime(), value.Body)
	}
	fmt.Println(len(ret))


	ret,_ = search("", 3,2)
	for _, value := range ret {
		fmt.Println(value.Id, value.getTime(), value.Body)
	}
	fmt.Println(len(ret))
	if ret[9].Id != 88 {
		t.Fail()
	}

	ret,_ = search("", 10,-1)
	for _, value := range ret {
		fmt.Println(value.Id, value.getTime(), value.Body)
	}
	fmt.Println(len(ret))
	if ret[9].Id != 89 {
		t.Fail()
	}
}