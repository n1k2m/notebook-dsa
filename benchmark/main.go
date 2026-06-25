package main

import (
	"fmt"
	"notebook-dsa/avl"
	"notebook-dsa/bst"
	"notebook-dsa/event"
	"notebook-dsa/hashtable"
	"time"
)

func main() {
	sizes := []int{100, 1000, 10000}

	fmt.Println("Сравнение времени поиска")
	fmt.Println("----------------------------------------")
	fmt.Printf("%-12s | %8s | %s\n", "structure", "n", "time")
	fmt.Println("----------------------------------------")

	for _, n := range sizes {
		events := generateEvents(n)
		target := events[n-1]

		runBST(n, events, target)
		runAVL(n, events, target)
		runHashTable(n, events, target)
	}
}

func runBST(n int, events []event.Event, target event.Event) {
	tree := bst.New()
	for _, e := range events {
		tree.Add(e)
	}
	duration := measureSearch(func() {
		tree.FindByDatetime(target.Year, target.Month, target.Day, target.Hour, target.Minute)
	})
	fmt.Printf("%-12s | %8d | %12s\n", "bst", n, duration)
}

func runAVL(n int, events []event.Event, target event.Event) {
	tree := avl.New()
	for _, e := range events {
		tree.Add(e)
	}
	duration := measureSearch(func() {
		tree.FindByDatetime(target.Year, target.Month, target.Day, target.Hour, target.Minute)
	})
	fmt.Printf("%-12s | %8d | %12s\n", "avl", n, duration)
}

func runHashTable(n int, events []event.Event, target event.Event) {
	table := hashtable.New()
	for _, e := range events {
		table.Add(e)
	}
	duration := measureSearch(func() {
		table.FindByDatetime(target.Year, target.Month, target.Day, target.Hour, target.Minute)
	})
	fmt.Printf("%-12s | %8d | %12s\n", "hashtable", n, duration)
}

// measureSearch усредняет несколько чистых поисков без пересборки структуры
func measureSearch(search func()) time.Duration {
	const attempts = 10000
	var total time.Duration
	for i := 0; i < attempts; i++ {
		start := time.Now()
		search()
		total += time.Since(start)
	}
	return total / time.Duration(attempts)
}

func generateEvents(n int) []event.Event {
	events := make([]event.Event, 0, n)
	year, month, day, hour, minute := 2024, 1, 1, 0, 0
	for i := 0; i < n; i++ {
		events = append(events, event.NewEvent(
			year, month, day, hour, minute,
			fmt.Sprintf("Событие %d", i+1),
			fmt.Sprintf("Локация-%d", i%10),
		))
		minute++
		if minute >= 60 {
			minute = 0
			hour++
		}
		if hour >= 24 {
			hour = 0
			day++
		}
		if day > 28 {
			day = 1
			month++
		}
		if month > 12 {
			month = 1
			year++
		}
	}
	return events
}
