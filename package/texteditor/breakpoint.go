package texteditor

type BrkType byte

const (
	// BrkTypeNone means no breakpoint.
	BrkTypeNone BrkType = ' '
	// BrkType0 means no breakpoint enabled.
	BrkType0 BrkType = '0'
	// BrkType1 means one breakpoint enabled.
	BrkType1 BrkType = '1'
	// BrkType2 means two breakpoints enabled.
	BrkType2 BrkType = '2'
	BrkType3 BrkType = '3'
	BrkType4 BrkType = '4'
	BrkType5 BrkType = '5'
	BrkType6 BrkType = '6'
	BrkType7 BrkType = '7'
	BrkType8 BrkType = '8'
	BrkType9 BrkType = '9'
	// BrkTypeN means more than 9 breakpoints enabled.
	BrkTypeN BrkType = 'N'
)

func brkEnabledCountToBrkType(count int) BrkType {
	switch count {
	case 0:
		return BrkType0
	case 1:
		return BrkType1
	case 2:
		return BrkType2
	case 3:
		return BrkType3
	case 4:
		return BrkType4
	case 5:
		return BrkType5
	case 6:
		return BrkType6
	case 7:
		return BrkType7
	case 8:
		return BrkType8
	case 9:
		return BrkType9
	default:
		return BrkTypeN
	}
}

type location struct {
	filename string
	line     int
}

type idMapValueTy struct {
	loc     location
	enabled bool
}

type locMapValueTy map[int]bool

type BreakpointMap struct {
	idMap  map[int]*idMapValueTy
	locMap map[location]locMapValueTy
}

// NewBreakpointMap creates a new BreakpointMap.
func NewBreakpointMap() *BreakpointMap {
	return &BreakpointMap{
		idMap:  make(map[int]*idMapValueTy),
		locMap: make(map[location]locMapValueTy),
	}
}

// Add adds a breakpoint to the map.
func (b *BreakpointMap) Add(id int, filename string, line int) {
	idMapValue := &idMapValueTy{
		loc:     location{filename, line},
		enabled: true,
	}
	b.idMap[id] = idMapValue
	v, ok := b.locMap[idMapValue.loc]
	if !ok {
		v = make(locMapValueTy)
		b.locMap[idMapValue.loc] = v
	}
	v[id] = true
}

// RemoveByID removes a breakpoint from the map by ID.
func (b *BreakpointMap) RemoveByID(id int) {
	idMapValue, ok := b.idMap[id]
	if !ok {
		return
	}
	delete(b.idMap, id)
	delete(b.locMap[idMapValue.loc], id)
}

// RemoveByLocation removes a breakpoint from the map by location.
func (b *BreakpointMap) RemoveByLocation(filename string, line int) {
	loc := location{filename, line}
	for id := range b.locMap[loc] {
		delete(b.idMap, id)
	}
	delete(b.locMap, loc)
}

// RemoveAll removes all breakpoints from the map.
func (b *BreakpointMap) RemoveAll() {
	b.idMap = make(map[int]*idMapValueTy)
	b.locMap = make(map[location]locMapValueTy)
}

// GetByID returns the location of a breakpoint by ID.
func (b *BreakpointMap) GetByID(id int) (string, int, bool) {
	value, ok := b.idMap[id]
	if !ok {
		return "", 0, false
	}
	return value.loc.filename, value.loc.line, true
}

func (b *BreakpointMap) CalculateBrkType(filename string, line int) BrkType {
	loc := location{filename, line}
	ids, ok := b.locMap[loc]
	if !ok {
		return BrkTypeNone
	}

	enabledCount := 0
	for _, v := range ids {
		if v {
			enabledCount++
		}
	}
	return brkEnabledCountToBrkType(enabledCount)
}
