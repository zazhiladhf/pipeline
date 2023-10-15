package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type User struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Salary int    `json:"salary"`
}

func main() {
	now := time.Now()
	asynchornous()
	// synchronous()
	log.Println("Done in", time.Since(now).Seconds())
}

// func synchronous() {
// 	users, _ := readFile("./data.json")
// 	users = usdToidr(users)
// 	log.Println(users[0], cap(users))
// }

func asynchornous() {

	userCh, _ := readFileConcurrent("./data.json")

	_ = usdToidrConcurrent(userCh)
	_ = usdToidrConcurrent(userCh)
	_ = usdToidrConcurrent(userCh)
	_ = usdToidrConcurrent(userCh)
	_ = usdToidrConcurrent(userCh)
	_ = usdToidrConcurrent(userCh)

	// users := merge(
	// 	usdToidrConcurrent(userCh),
	// 	usdToidrConcurrent(userCh),
	// 	usdToidrConcurrent(userCh),
	// 	usdToidrConcurrent(userCh),
	// 	usdToidrConcurrent(userCh),
	// 	usdToidrConcurrent(userCh),
	// )

	done := make(chan bool)

	writeToFileConcurrent(userCh, done)

	if <-done {
		log.Println("Donee")
	}
}

// func readFile(filename string) (users []User, err error) {
// 	now := time.Now()
// 	dataByte, err := os.ReadFile(filename)
// 	if err != nil {
// 		return
// 	}

// 	users = []User{}

// 	err = json.Unmarshal(dataByte, &users)
// 	if err != nil {
// 		return
// 	}

// 	log.Println("success read data in", time.Since(now).Seconds(), "s")
// 	return
// }

func readFileConcurrent(filename string) (<-chan User, error) {
	userCh := make(chan User)
	now := time.Now()
	dataByte, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	users := []User{}

	err = json.Unmarshal(dataByte, &users)
	if err != nil {
		return nil, err
	}

	go func() {
		for _, user := range users {
			userCh <- user
		}

		close(userCh)
	}()

	log.Println("success read data in", time.Since(now).Seconds(), "s")
	return userCh, nil
}

// func usdToidr(users []User) (newUsers []User) {
// 	now := time.Now()
// 	for _, user := range users {
// 		time.Sleep(10 * time.Millisecond)
// 		user.Salary *= 100
// 		newUsers = append(newUsers, user)
// 	}
// 	log.Println("success change usd to idr in", time.Since(now).Seconds(), "s")

// 	return
// }

func usdToidrConcurrent(user <-chan User) <-chan User {
	now := time.Now()
	userCh := make(chan User)

	go func() {
		for u := range user {
			time.Sleep(10 * time.Millisecond)
			newData := u
			newData.Salary = newData.Salary * 15_000
			userCh <- newData
		}
		close(userCh)
	}()

	log.Println("success change usd to idr in", time.Since(now).Nanoseconds(), "ns")

	return userCh
}

// func merge(users ...<-chan User) <-chan User {
// 	usersCh := make(chan User)

// 	wg := sync.WaitGroup{}

// 	for _, user := range users {
// 		wg.Add(1)
// 		go func(user <-chan User) {
// 			for u := range user {
// 				usersCh <- u
// 			}
// 			wg.Done()
// 		}(user)
// 	}
// 	go func() {
// 		wg.Wait()
// 		close(usersCh)
// 	}()
// 	return usersCh
// }

func writeToFileConcurrent(dataCh <-chan User, done chan bool) {
	wg := sync.WaitGroup{}

	for data := range dataCh {
		wg.Add(1)
		go func(data User) {
			user, _ := json.Marshal(data)
			err := os.WriteFile("./users/"+data.Name+".json", user, 0666)
			if err != nil {
				log.Println("error when try to write file", err)
			}
			wg.Done()
		}(data)
	}

	go func() {
		wg.Wait()
		done <- true
		done <- true
		done <- true
		done <- true
		done <- true
		done <- true
		// close(done)
	}()
}
