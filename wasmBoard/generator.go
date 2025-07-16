package wasmBoard

import (
	"errors"
	"gochess/chessBoard"
	"os"
	"os/exec"
	"syscall/js"
)

type WASMBoard struct {
	board *chessBoard.Board
}

func NewWASMBoard(fen string) *WASMBoard {
	board := chessBoard.New()
	board.LoadBoard(fen)
	return &WASMBoard{
		board: board,
	}
}

func (b *WASMBoard) getValidMovesForSquare(squareId string) []string {
	ri, ci := int(squareId[1]-'1'), int(squareId[0]-'a')
	square := b.board.GetSquare(ri, ci)
	squaresInNotation := []string{}
	for _, square := range square.LegalMoves {
		squaresInNotation = append(squaresInNotation, square.File+square.Rank)
	}
	return squaresInNotation
}

func generateWasmModule() error {

	chessNamespace := js.Global().Get("Object").New()
	chessNamespace.Set("generateBoardFromFen", js.FuncOf(NewWASMBoard))

	js.Global().Set("chessWASM", chessNamespace)

	cmd := exec.Command("go", "build", "-o", "chessBoard.wasm", "wasmBoard.go")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	return cmd.Run()
}

func CreateChessWASMBinary() ([]byte, error) {
	if err := generateWasmModule(); err != nil {
		return nil, errors.New("failed to generate wasm module: " + err.Error())
	}

	file, err := os.ReadFile("chessBoard.wasm")

	if err != nil {
		return nil, errors.New("failed to generate wasm module: " + err.Error())
	}

	return file, nil
}
