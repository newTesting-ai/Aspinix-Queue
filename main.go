package main

import (
	queuemanager "aspinix-queue/core/internals/queue_manager"
	"fmt"
)

func main() {

	heap := queuemanager.InitHeap()
	data := map[string]string{
		"data": "value",
	}
	heap.AddData(data, false)
	heap.AddData(data, false)
	heap.AddData(data, false)
	usages, percentage := heap.GetUsages()
	fmt.Println("Usages:", usages)
	fmt.Println("Percentage:", percentage)
	heap.AddData(data, false)
	heap.AddData(data, true)
	usages, percentage = heap.GetUsages()
	fmt.Println("Usages:", usages)
	fmt.Println("Percentage:", percentage)
	heap.AddData(data, true)
	heap.AddData(data, false)
	heap.AddData(data, true)
	usages, percentage = heap.GetUsages()
	fmt.Println("Usages:", usages)
	fmt.Println("Percentage:", percentage)
	heap.AddData(data, false)
	usages, percentage = heap.GetUsages()
	fmt.Println("Usages:", usages)
	fmt.Println("Percentage:", percentage)
}
