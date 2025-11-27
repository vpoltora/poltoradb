package main

import (
	"fmt"

	"github.com/vpoltora/poltoradb/internal/pagemanager"
)

func main() {
	pageManager := pagemanager.New()

	page, err := pageManager.AllocatePage()
	if err != nil {
		fmt.Println("error allocating page:", err)

		return
	}

	fmt.Printf("Allocated page: %+v\n", page)

	header := page.GetHeader()
	fmt.Printf("Page header: %+v\n", header)

	slots := page.Slots()
	fmt.Printf("slots: #%v", slots)
}
