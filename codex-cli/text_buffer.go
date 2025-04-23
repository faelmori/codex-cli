package main

import (
	"strings"
	"unicode/utf8"
)

type Direction string

const (
	Left      Direction = "left"
	Right     Direction = "right"
	Up        Direction = "up"
	Down      Direction = "down"
	WordLeft  Direction = "wordLeft"
	WordRight Direction = "wordRight"
	Home      Direction = "home"
	End       Direction = "end"
)

func isWordChar(ch rune) bool {
	return !strings.ContainsRune(" \t\n\r.,;!?()[]{}", ch)
}

type Viewport struct {
	Height int
	Width  int
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func toCodePoints(str string) []rune {
	return []rune(str)
}

func cpLen(str string) int {
	return utf8.RuneCountInString(str)
}

func cpSlice(str string, start, end int) string {
	return string([]rune(str)[start:end])
}

type TextBuffer struct {
	lines        []string
	cursorRow    int
	cursorCol    int
	scrollRow    int
	scrollCol    int
	preferredCol *int
	version      int
	undoStack    []snapshot
	redoStack    []snapshot
	historyLimit int
	clipboard    *string
}

type snapshot struct {
	lines []string
	row   int
	col   int
}

func NewTextBuffer(text string) *TextBuffer {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	return &TextBuffer{
		lines:        lines,
		historyLimit: 100,
	}
}

func (tb *TextBuffer) line(r int) string {
	if r < 0 || r >= len(tb.lines) {
		return ""
	}
	return tb.lines[r]
}

func (tb *TextBuffer) lineLen(r int) int {
	return cpLen(tb.line(r))
}

func (tb *TextBuffer) ensureCursorInRange() {
	tb.cursorRow = clamp(tb.cursorRow, 0, len(tb.lines)-1)
	tb.cursorCol = clamp(tb.cursorCol, 0, tb.lineLen(tb.cursorRow))
}

func (tb *TextBuffer) snapshot() snapshot {
	return snapshot{
		lines: append([]string{}, tb.lines...),
		row:   tb.cursorRow,
		col:   tb.cursorCol,
	}
}

func (tb *TextBuffer) pushUndo() {
	tb.undoStack = append(tb.undoStack, tb.snapshot())
	if len(tb.undoStack) > tb.historyLimit {
		tb.undoStack = tb.undoStack[1:]
	}
	tb.redoStack = nil
}

func (tb *TextBuffer) restore(state *snapshot) bool {
	if state == nil {
		return false
	}
	tb.lines = append([]string{}, state.lines...)
	tb.cursorRow = state.row
	tb.cursorCol = state.col
	tb.ensureCursorInRange()
	return true
}

func (tb *TextBuffer) ensureCursorVisible(vp Viewport) {
	if tb.cursorRow < tb.scrollRow {
		tb.scrollRow = tb.cursorRow
	} else if tb.cursorRow >= tb.scrollRow+vp.Height {
		tb.scrollRow = tb.cursorRow - vp.Height + 1
	}

	if tb.cursorCol < tb.scrollCol {
		tb.scrollCol = tb.cursorCol
	} else if tb.cursorCol >= tb.scrollCol+vp.Width {
		tb.scrollCol = tb.cursorCol - vp.Width + 1
	}
}

func (tb *TextBuffer) GetVersion() int {
	return tb.version
}

func (tb *TextBuffer) GetCursor() (int, int) {
	return tb.cursorRow, tb.cursorCol
}

func (tb *TextBuffer) GetVisibleLines(vp Viewport) []string {
	tb.ensureCursorVisible(vp)
	return tb.lines[tb.scrollRow:clamp(tb.scrollRow+vp.Height, 0, len(tb.lines))]
}

func (tb *TextBuffer) GetText() string {
	return strings.Join(tb.lines, "\n")
}

func (tb *TextBuffer) GetLines() []string {
	return append([]string{}, tb.lines...)
}

func (tb *TextBuffer) Undo() bool {
	if len(tb.undoStack) == 0 {
		return false
	}
	state := tb.undoStack[len(tb.undoStack)-1]
	tb.undoStack = tb.undoStack[:len(tb.undoStack)-1]
	tb.redoStack = append(tb.redoStack, tb.snapshot())
	tb.restore(&state)
	tb.version++
	return true
}

func (tb *TextBuffer) Redo() bool {
	if len(tb.redoStack) == 0 {
		return false
	}
	state := tb.redoStack[len(tb.redoStack)-1]
	tb.redoStack = tb.redoStack[:len(tb.redoStack)-1]
	tb.undoStack = append(tb.undoStack, tb.snapshot())
	tb.restore(&state)
	tb.version++
	return true
}

func (tb *TextBuffer) Insert(ch string) {
	if strings.ContainsAny(ch, "\n\r") {
		tb.InsertStr(ch)
		return
	}

	tb.pushUndo()

	line := tb.line(tb.cursorRow)
	tb.lines[tb.cursorRow] = cpSlice(line, 0, tb.cursorCol) + ch + cpSlice(line, tb.cursorCol, len(line))
	tb.cursorCol += cpLen(ch)
	tb.version++
}

func (tb *TextBuffer) Newline() {
	tb.pushUndo()

	line := tb.line(tb.cursorRow)
	before := cpSlice(line, 0, tb.cursorCol)
	after := cpSlice(line, tb.cursorCol, len(line))

	tb.lines[tb.cursorRow] = before
	tb.lines = append(tb.lines[:tb.cursorRow+1], append([]string{after}, tb.lines[tb.cursorRow+1:]...)...)

	tb.cursorRow++
	tb.cursorCol = 0
	tb.version++
}

func (tb *TextBuffer) Backspace() {
	if tb.cursorCol == 0 && tb.cursorRow == 0 {
		return
	}

	tb.pushUndo()

	if tb.cursorCol > 0 {
		line := tb.line(tb.cursorRow)
		tb.lines[tb.cursorRow] = cpSlice(line, 0, tb.cursorCol-1) + cpSlice(line, tb.cursorCol, len(line))
		tb.cursorCol--
	} else if tb.cursorRow > 0 {
		prev := tb.line(tb.cursorRow - 1)
		cur := tb.line(tb.cursorRow)
		tb.lines[tb.cursorRow-1] = prev + cur
		tb.lines = append(tb.lines[:tb.cursorRow], tb.lines[tb.cursorRow+1:]...)
		tb.cursorRow--
		tb.cursorCol = cpLen(prev)
	}
	tb.version++
}

func (tb *TextBuffer) Del() {
	line := tb.line(tb.cursorRow)
	if tb.cursorCol < tb.lineLen(tb.cursorRow) {
		tb.pushUndo()
		tb.lines[tb.cursorRow] = cpSlice(line, 0, tb.cursorCol) + cpSlice(line, tb.cursorCol+1, len(line))
	} else if tb.cursorRow < len(tb.lines)-1 {
		tb.pushUndo()
		next := tb.line(tb.cursorRow + 1)
		tb.lines[tb.cursorRow] = line + next
		tb.lines = append(tb.lines[:tb.cursorRow+1], tb.lines[tb.cursorRow+2:]...)
	}
	tb.version++
}

func (tb *TextBuffer) DeleteWordLeft() {
	if tb.cursorCol == 0 && tb.cursorRow == 0 {
		return
	}

	if tb.cursorCol == 0 {
		tb.Backspace()
		return
	}

	tb.pushUndo()

	line := tb.line(tb.cursorRow)
	arr := toCodePoints(line)

	start := tb.cursorCol
	onlySpaces := true
	for i := 0; i < start; i++ {
		if isWordChar(arr[i]) {
			onlySpaces = false
			break
		}
	}

	if onlySpaces && start > 0 {
		start--
	} else {
		for start > 0 && !isWordChar(arr[start-1]) {
			start--
		}
		for start > 0 && isWordChar(arr[start-1]) {
			start--
		}
	}

	tb.lines[tb.cursorRow] = cpSlice(line, 0, start) + cpSlice(line, tb.cursorCol, len(line))
	tb.cursorCol = start
	tb.version++
}

func (tb *TextBuffer) DeleteWordRight() {
	line := tb.line(tb.cursorRow)
	arr := toCodePoints(line)
	if tb.cursorCol >= len(arr) && tb.cursorRow == len(tb.lines)-1 {
		return
	}

	if tb.cursorCol >= len(arr) {
		tb.Del()
		return
	}

	tb.pushUndo()

	end := tb.cursorCol
	for end < len(arr) && !isWordChar(arr[end]) {
		end++
	}
	for end < len(arr) && isWordChar(arr[end]) {
		end++
	}

	tb.lines[tb.cursorRow] = cpSlice(line, 0, tb.cursorCol) + cpSlice(line, end, len(line))
	tb.version++
}

func (tb *TextBuffer) Move(dir Direction) {
	switch dir {
	case Left:
		tb.preferredCol = nil
		if tb.cursorCol > 0 {
			tb.cursorCol--
		} else if tb.cursorRow > 0 {
			tb.cursorRow--
			tb.cursorCol = tb.lineLen(tb.cursorRow)
		}
	case Right:
		tb.preferredCol = nil
		if tb.cursorCol < tb.lineLen(tb.cursorRow) {
			tb.cursorCol++
		} else if tb.cursorRow < len(tb.lines)-1 {
			tb.cursorRow++
			tb.cursorCol = 0
		}
	case Up:
		if tb.cursorRow > 0 {
			if tb.preferredCol == nil {
				tb.preferredCol = &tb.cursorCol
			}
			tb.cursorRow--
			tb.cursorCol = clamp(*tb.preferredCol, 0, tb.lineLen(tb.cursorRow))
		}
	case Down:
		if tb.cursorRow < len(tb.lines)-1 {
			if tb.preferredCol == nil {
				tb.preferredCol = &tb.cursorCol
			}
			tb.cursorRow++
			tb.cursorCol = clamp(*tb.preferredCol, 0, tb.lineLen(tb.cursorRow))
		}
	case Home:
		tb.preferredCol = nil
		tb.cursorCol = 0
	case End:
		tb.preferredCol = nil
		tb.cursorCol = tb.lineLen(tb.cursorRow)
	case WordLeft:
		tb.preferredCol = nil
		regex := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "", ".", "", ",", "", ";", "", "!", "", "?", "", "(", "", ")", "", "[", "", "]", "", "{", "", "}", "")
		slice := regex.Replace(cpSlice(tb.line(tb.cursorRow), 0, tb.cursorCol))
		lastIdx := 0
		for i, ch := range slice {
			if !isWordChar(ch) {
				lastIdx = i
			}
		}
		tb.cursorCol = lastIdx
	case WordRight:
		tb.preferredCol = nil
		regex := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "", ".", "", ",", "", ";", "", "!", "", "?", "", "(", "", "", ")", "", "[", "", "]", "", "{", "", "}", "")
		l := tb.line(tb.cursorRow)
		moved := false
		for i, ch := range l {
			if !isWordChar(ch) && i > tb.cursorCol {
				tb.cursorCol = i
				moved = true
				break
			}
		}
		if !moved {
			tb.cursorCol = tb.lineLen(tb.cursorRow)
		}
	}

	if dir != Up && dir != Down {
		tb.preferredCol = nil
	}
}

func (tb *TextBuffer) InsertStr(str string) bool {
	if str == "" {
		return false
	}

	normalised := strings.ReplaceAll(strings.ReplaceAll(str, "\r\n", "\n"), "\r", "\n")

	if !strings.Contains(normalised, "\n") {
		tb.Insert(normalised)
		return true
	}

	tb.pushUndo()

	parts := strings.Split(normalised, "\n")
	before := cpSlice(tb.line(tb.cursorRow), 0, tb.cursorCol)
	after := cpSlice(tb.line(tb.cursorRow), tb.cursorCol, len(tb.line(tb.cursorRow)))

	tb.lines[tb.cursorRow] = before + parts[0]

	if len(parts) > 2 {
		middle := parts[1 : len(parts)-1]
		tb.lines = append(tb.lines[:tb.cursorRow+1], append(middle, tb.lines[tb.cursorRow+1:]...)...)
	}

	last := parts[len(parts)-1] + after
	tb.lines = append(tb.lines[:tb.cursorRow+1], append([]string{last}, tb.lines[tb.cursorRow+1:]...)...)

	tb.cursorRow += len(parts) - 1
	tb.cursorCol = cpLen(parts[len(parts)-1])
	tb.version++
	return true
}

func (tb *TextBuffer) StartSelection() {
	tb.selectionAnchor = &snapshot{
		row: tb.cursorRow,
		col: tb.cursorCol,
	}
}

func (tb *TextBuffer) EndSelection() {
	tb.selectionAnchor = nil
}

func (tb *TextBuffer) getSelectedText() *string {
	if tb.selectionAnchor == nil {
		return nil
	}
	ar, ac := tb.selectionAnchor.row, tb.selectionAnchor.col
	br, bc := tb.cursorRow, tb.cursorCol

	if ar == br && ac == bc {
		return nil
	}

	topBefore := ar < br || (ar == br && ac < bc)
	sr, sc, er, ec := ar, ac, br, bc
	if !topBefore {
		sr, sc, er, ec = br, bc, ar, ac
	}

	if sr == er {
		text := cpSlice(tb.line(sr), sc, ec)
		return &text
	}

	parts := []string{cpSlice(tb.line(sr), sc, len(tb.line(sr)))}
	for r := sr + 1; r < er; r++ {
		parts = append(parts, tb.line(r))
	}
	parts = append(parts, cpSlice(tb.line(er), 0, ec))
	text := strings.Join(parts, "\n")
	return &text
}

func (tb *TextBuffer) Copy() *string {
	text := tb.getSelectedText()
	if text == nil {
		return nil
	}
	tb.clipboard = text
	return text
}

func (tb *TextBuffer) Paste() bool {
	if tb.clipboard == nil {
		return false
	}
	return tb.InsertStr(*tb.clipboard)
}

func (tb *TextBuffer) HandleInput(input *string, key map[string]bool, vp Viewport) bool {
	beforeVer := tb.version
	beforeRow, beforeCol := tb.GetCursor()

	if key["escape"] {
		return false
	}

	if key["return"] || (input != nil && (*input == "\r" || *input == "\n")) {
		tb.Newline()
	} else if key["leftArrow"] && !key["meta"] && !key["ctrl"] && !key["alt"] {
		tb.Move(Left)
	} else if key["rightArrow"] && !key["meta"] && !key["ctrl"] && !key["alt"] {
		tb.Move(Right)
	} else if key["upArrow"] {
		tb.Move(Up)
	} else if key["downArrow"] {
		tb.Move(Down)
	} else if (key["meta"] || key["ctrl"] || key["alt"]) && key["leftArrow"] {
		tb.Move(WordLeft)
	} else if (key["meta"] || key["ctrl"] || key["alt"]) && key["rightArrow"] {
		tb.Move(WordRight)
	} else if key["home"] {
		tb.Move(Home)
	} else if key["end"] {
		tb.Move(End)
	} else if (key["meta"] || key["ctrl"] || key["alt"]) && key["backspace"] {
		tb.DeleteWordLeft()
	} else if (key["meta"] || key["ctrl"] || key["alt"]) && key["delete"] {
		tb.DeleteWordRight()
	} else if key["backspace"] || (input != nil && *input == "\x7f") || (key["delete"] && !key["shift"]) {
		tb.Backspace()
	} else if key["delete"] {
		tb.Del()
	} else if input != nil && !key["ctrl"] && !key["meta"] {
		tb.Insert(*input)
	}

	tb.ensureCursorInRange()
	tb.ensureCursorVisible(vp)

	cursorMoved := tb.cursorRow != beforeRow || tb.cursorCol != beforeCol

	return tb.version != beforeVer || cursorMoved
}
