package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
)

func main() {
	var numbers []int
	for len(numbers) < 100 {
		number := rand.Intn(99999-10000) + 10000
		if Valid(number) {
			numbers = append(numbers, number)
		}
	}
	dataString2 := `{
    "match": "Bork",
    "reward": 10,
    "reward_type": "%"
} `
	URLGoods := "http://localhost:10001/api/goods"
	data2 := []byte(dataString2)
	buf2 := bytes.NewBuffer(data2)
	resp, err := http.Post(URLGoods, "application/json", buf2)
	if err := resp.Body.Close(); err != nil {
		return
	}
	if err != nil {
		return
	}

	URL := "http://localhost:10001/api/orders"

	for k, i := range numbers {
		str := fmt.Sprintf(`{"order": "%d", "goods": [{"description": "Чайник Bork", "price": %d}]}`, i, k*100)

		data := []byte(str)
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(data))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err := resp.Body.Close(); err != nil {
			return
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(i)

	}
}

func Valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
