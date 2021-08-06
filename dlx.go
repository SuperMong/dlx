// Package dlx implements Dancing Link X.
package dlx

type SolutionAccepter interface {
	AcceptSolution(exactCover [][]int) bool
}

type SolutionAccepterFunc func([][]int) bool

func (f SolutionAccepterFunc) AcceptSolution(exactCover [][]int) bool {
	return f(exactCover)
}

type Matrix interface {
	AddRow(row ...int)
	Solve(accepter SolutionAccepter) bool
}

type node struct {
	rowStart *node
	column   *colNode

	up, down    *node
	left, right *node
}

type colNode struct {
	node

	size  int
	index int
}

type matrix struct {
	head    *colNode
	columns []colNode
	answer  []*colNode
}

// NewMatrix create a Matrix with n columns.
func NewMatrix(n int) Matrix {
	columns := make([]colNode, n+1)
	head := &columns[0]
	head.column = head
	head.up = &head.node
	head.down = &head.node
	head.left = &columns[n].node
	head.index = -1

	preCol := head
	for i := 1; i != n+1; i++ {
		columns[i].size = 0
		columns[i].index = i
		columns[i].column = &columns[i]
		columns[i].up = &columns[i].node
		columns[i].down = &columns[i].node
		columns[i].left = &preCol.node
		preCol.right = &columns[i].node
		preCol = &columns[i]
	}

	return &matrix{head, columns, nil}
}

// AddRow adds a row to matrix.
func (m *matrix) AddRow(row ...int) {
	count := len(row)
	if count == 0 {
		return
	}

	// Add nodes in row.
	line := make([]node, count)
	start := &line[0]
	for i, v := range row {
		newNode := &line[i]
		column := &m.columns[v]
		newNode.rowStart = start
		newNode.column = column
		column.size++

		// Get neighbors.
		newNode.right = &line[(i+1)%count]
		newNode.left = &line[(i-1+count)%count]
		newNode.up = column.up
		newNode.down = column.up.down

		// Change neighbors.
		newNode.down.up, newNode.up.down = newNode, newNode
		newNode.left.right, newNode.right.left = newNode, newNode
	}
}

// Solve the problem.
func (m *matrix) Solve(accepter SolutionAccepter) bool {
	return m.solve(accepter)
}

func (m *matrix) solve(accepter SolutionAccepter) bool {
	head := m.head
	if head.right == &head.node && accepter.AcceptSolution(m.getExactCover()) {
		return true
	}

	// Randomly pick one column, we pick the smallest one.
	col := head.right.column
	min := col
	for col != head.column {
		if col.size < min.size {
			min = col
		}
		col = col.right.column
	}

	// Back Tracking
	m.cover(col)
	m.answer = append(m.answer, nil)
	for d := col.down; d != &col.node; d = d.down {
		m.answer[len(m.answer)-1] = d.column
		for r := d.right; r != d; r = r.right {
			m.cover(r.column)
		}
	}
	if m.solve(accepter) {
		return true
	}
	for u := col.up; u != &col.node; u = u.up {
		for l := u.left; l != u; l = l.left {
			m.uncover(l.column)
		}
	}
	m.uncover(col)
	m.answer = m.answer[:len(m.answer)-1]

	return false
}

// Translate the answer.
func (m *matrix) getExactCover() [][]int {
	length := len(m.answer)
	ec := make([][]int, length)
	for i, n := range m.answer {
		start := n.rowStart
		row := make([]int, 0)
		for rn := start.right; rn != start; rn = rn.right {
			row = append(row, rn.column.index)
		}
		ec[i] = row
	}
	return ec
}

// Cover the column and all rows intersect it.
func (m *matrix) cover(column *colNode) {
	for dn := column.down; dn != &column.node; dn = dn.down {
		for rn := dn.right; rn != dn; rn = rn.right {
			rn.up.down, rn.down.up = rn.down, rn.up
			rn.left.right, rn.right.left = rn.right, rn.left
		}
	}
}

// Uncover the column and all rows intersect it.
func (m *matrix) uncover(column *colNode) {
	for un := column.up; un != &column.node; un = un.up {
		for ln := un.left; ln != un; ln = ln.left {
			ln.down.up, ln.up.down = ln, ln
			ln.left.right, ln.right.left = ln, ln
		}
	}
}
