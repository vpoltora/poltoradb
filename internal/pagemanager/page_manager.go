// Package pagemanager: This file defines the Page struct used in the database storage system.
package pagemanager

type PageManagerInterface interface {
	AllocatePage() (*Page, error)
}

type PageManager struct {
}

func New() PageManagerInterface {
	return &PageManager{}
}

func (pm *PageManager) AllocatePage() (*Page, error) {
	data := [PageSize]byte{}

	return &Page{
		Data: data,
	}, nil
}
