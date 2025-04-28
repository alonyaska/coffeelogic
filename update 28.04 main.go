package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{
	"log in",
	"log out",
	"check balance",
	"pay balance",
	"buy coffe",
}

var products = []string{
	"coffe",
	"latte",
	"raf",
	"capuchino",
	"cake",
	"tea",
}

type logItem struct {
	action    string
	timestamp time.Time
	productID string
}

type User struct {
	id   int
	name string
	logs []logItem
}

func (u User) getActivity() string {
	out := fmt.Sprintf("ID: %d | Name: %s \n ActivityLog: \n ", u.id, u.name)
	for i, item := range u.logs {
		out += fmt.Sprintf("%d. user check: %s, [%s] at %s\n", i+1, item.productID, item.action, item.timestamp)
	}

	return out
}

func main() {

	rand.Seed(time.Now().Unix())

	users := make(chan User)
	go generateUsers(100, users) //  число юзеров

	wg := &sync.WaitGroup{}

	t := time.Now()

	for user := range users {
		wg.Add(1)
		go saveUserInfo(user, wg)
	}

	wg.Wait()
	analyzeLogs([]User{})
	fmt.Println("time elapsed:", time.Since(t).String())
}

func saveUserInfo(user User, wg *sync.WaitGroup) error {
	time.Sleep(time.Millisecond * 10)
	fmt.Printf("writing info for %d \n", user.id)

	filename := fmt.Sprintf("logs/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString(user.getActivity())
	if err != nil {
		return nil
	}

	wg.Done()

	return nil
}

func generateUsers(count int, users chan User) {

	for i := 0; i < count; i++ {
		users <- User{
			id:   i + 1,
			name: fmt.Sprintf("name%d", i+1),
			logs: generateLogs(rand.Intn(10)), // число логов
		}

	}
	close(users)
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			timestamp: time.Now(),
			action:    actions[rand.Intn(len(actions))],
			productID: products[rand.Intn(len(products)-1)],
		}

	}
	return logs
}

func analyzeLogs(users []User) {
	stats := make(map[string]struct {
		views    int
		purchase int
	})

	for _, user := range users {
		for _, log := range user.logs {
			s := stats[log.productID]
			if log.action == "buy coffe" {
				s.purchase++
			} else {
				s.views++
			}
			stats[log.productID] = s

		}

	}
	fmt.Println("\n=== Product Conversion Rate ===")
	for products, s := range stats {
		total := s.views + s.purchase
		if total == 0 {
			continue
		}
		percent := float64(s.purchase) / float64(total) * 100
		fmt.Printf("%s: %.2f%% конверсии (%d просмотров / %d покупок)\n",
			products, percent, s.views, s.purchase)
		saveAnalyzeInfo(stats)
	}
}

func saveAnalyzeInfo(stats map[string]struct {
	views    int
	purchase int
}) error {
	file, err := os.Create("logs/analyze.txt") // можно изменить путь
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("=== Product Conversion Rate ===\n")
	if err != nil {
		return err
	}

	for product, s := range stats {
		total := s.views + s.purchase
		if total == 0 {
			continue
		}
		percent := float64(s.purchase) / float64(total) * 100
		line := fmt.Sprintf("%s: %.2f%% конверсии (%d просмотров / %d покупок)\n",
			product, percent, s.views, s.purchase)
		_, err = file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}
