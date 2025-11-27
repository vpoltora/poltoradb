package pagemanager

type CellType uint8

const (
	CellTypeInt64  CellType = 0
	CellTypeString CellType = 1
	CellTypeBool   CellType = 2
	CellTypeNull   CellType = 3
)

// Value represents a single value in a Row.
type Cell struct {
	Type CellType

	Int64Value  int64
	StringValue string
	BoolValue   bool
}

// Row represents a single row in a Page.
type Row struct {
	Cells []Cell
}
