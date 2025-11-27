package pagemanager

import (
	"encoding/binary"
	"fmt"
)

const PageSize = 4096 // 4KB

type PageType uint8

const (
	PageTypeTable PageType = 0
	PageTypeIndex PageType = 1
)

// Header represents the metadata of a Page.
type Header struct {
	// table page, index page, free page
	Type PageType

	// number of filled rows in this page
	// max 65536 rows
	NumberOfSlots uint16

	// offset to the beginning of free space
	// used to insert new rows
	// stores the next index in the Rows byte array
	FreeSpacePointer uint16
}

type Slot struct {
	Offset uint16
	Length uint16
}

// Page represents a single page in the database storage system.
// page has fixed size of 4KB
type Page struct {
	Data [PageSize]byte
}

// GetHeader extracts the Header information from the Page's raw data.
// The first 5 bytes of the page are used for the header.
//
// Page type = 1 byte.
//
// NumberOfSlots = 2 bytes in little-endian format (default for x86/arm).
// little-endian (default for x86/arm) format says that the least significant byte is stored first:
// 0x1234 = [0x34, 0x12] = 0x34 XOR (0x12 << 8)
// or just call binary.LittleEndian.Uint16
//
// FreeSpacePointer is calculated based on the number of slots.
func (p *Page) GetHeader() *Header {
	// uint8 + uint16 + uint16 = 5 bytes

	pageType := PageType(p.Data[0])
	fmt.Printf("page type: #%v", pageType)

	numberOfSlots := binary.LittleEndian.Uint16(p.Data[1:3])

	// caused by offset (uint16) + length (uint16) = 4 bytes per slot
	freeSpacePointer := numberOfSlots*4 + 5

	return &Header{
		Type:             pageType,
		NumberOfSlots:    numberOfSlots,
		FreeSpacePointer: freeSpacePointer,
	}
}

func (p *Page) Slots() []Slot {
	header := p.GetHeader()

	slots := make([]Slot, header.NumberOfSlots)

	for i := uint16(0); i < header.NumberOfSlots; i++ {
		// each slot is 4 bytes: 2 bytes for offset, 2 bytes for length
		slotOffset := 5 + i*4
		offset := binary.LittleEndian.Uint16(p.Data[slotOffset : slotOffset+2])
		lenght := binary.LittleEndian.Uint16(p.Data[slotOffset+2 : slotOffset+4])

		slots[i] = Slot{
			Offset: offset,
			Length: lenght,
		}
	}

	return slots
}

func (p *Page) AddRow(rowBytes []byte) (slotID, offset uint16) {
	return 0, 0
}
