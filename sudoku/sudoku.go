// Package sudoku solve the sudoku problem.
package sudoku

import (
	"io"

	"github.com/SuperMong/dlx"
)

type SolutionReceiver interface {
	ReceiveSolution(exactCover [][]int) bool
}

type SolutionReceiveFunc func([][]int) bool

func (f SolutionReceiveFunc) ReceiveSolution(exactCover [][]int) bool {
	return f(exactCover)
}

func ReadProblem(reader io.Reader) dlx.Matrix {
	m := dlx.NewMatrix(384)
}

func SolveProblem(m dlx.Matrix) [][]int {

}

func encode(problem string) dlx.Matrix {

}

func decode(exactCover [][]int) string {

}
