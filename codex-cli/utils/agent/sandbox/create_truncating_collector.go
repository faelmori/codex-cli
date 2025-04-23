package sandbox

import (
	"bufio"
	"bytes"
	"io"
)

const (
	MAX_OUTPUT_BYTES = 1024 * 10 // 10 KB
	MAX_OUTPUT_LINES = 256
)

type TruncatingCollector struct {
	buffer    bytes.Buffer
	byteCount int
	lineCount int
	hitLimit  bool
}

func NewTruncatingCollector() *TruncatingCollector {
	return &TruncatingCollector{}
}

func (tc *TruncatingCollector) Write(p []byte) (n int, err error) {
	if tc.hitLimit {
		return len(p), nil
	}

	reader := bufio.NewReader(bytes.NewReader(p))
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return 0, err
		}

		if tc.byteCount+len(line) > MAX_OUTPUT_BYTES || tc.lineCount+1 > MAX_OUTPUT_LINES {
			tc.hitLimit = true
			break
		}

		tc.buffer.Write(line)
		tc.byteCount += len(line)
		tc.lineCount++

		if err == io.EOF {
			break
		}
	}

	return len(p), nil
}

func (tc *TruncatingCollector) String() string {
	return tc.buffer.String()
}

func (tc *TruncatingCollector) HitLimit() bool {
	return tc.hitLimit
}
