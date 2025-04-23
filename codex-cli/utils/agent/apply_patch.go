package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	ADD_FILE_PREFIX    = "+++ "
	DELETE_FILE_PREFIX = "--- "
	UPDATE_FILE_PREFIX = "+++ "
	PATCH_SUFFIX       = "*** End Patch ***"
	HUNK_ADD_LINE_PREFIX = "+"
	PATCH_PREFIX       = "*** Begin Patch ***"
)

type ActionType string

const (
	ADD    ActionType = "add"
	DELETE ActionType = "delete"
	UPDATE ActionType = "update"
)

type FileChange struct {
	Type       ActionType
	OldContent *string
	NewContent *string
	MovePath   *string
}

type Commit struct {
	Changes map[string]FileChange
}

type Chunk struct {
	OrigIndex int
	DelLines  []string
	InsLines  []string
}

type PatchAction struct {
	Type     ActionType
	NewFile  *string
	Chunks   []Chunk
	MovePath *string
}

type Patch struct {
	Actions map[string]PatchAction
}

type DiffError struct {
	Message string
}

func (e *DiffError) Error() string {
	return e.Message
}

func assembleChanges(orig map[string]*string, updatedFiles map[string]*string) Commit {
	commit := Commit{Changes: make(map[string]FileChange)}
	for p, newContent := range updatedFiles {
		oldContent := orig[p]
		if oldContent == newContent {
			continue
		}
		if oldContent != nil && newContent != nil {
			commit.Changes[p] = FileChange{
				Type:       UPDATE,
				OldContent: oldContent,
				NewContent: newContent,
			}
		} else if newContent != nil {
			commit.Changes[p] = FileChange{
				Type:       ADD,
				NewContent: newContent,
			}
		} else if oldContent != nil {
			commit.Changes[p] = FileChange{
				Type:       DELETE,
				OldContent: oldContent,
			}
		} else {
			panic("Unexpected state in assembleChanges")
		}
	}
	return commit
}

func textToPatch(text string, orig map[string]string) (Patch, int) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 2 || !strings.HasPrefix(lines[0], PATCH_PREFIX) || lines[len(lines)-1] != PATCH_SUFFIX {
		panic("Invalid patch text")
	}
	parser := newParser(orig, lines)
	parser.index = 1
	parser.parse()
	return parser.patch, parser.fuzz
}

func identifyFilesNeeded(text string) []string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	result := make(map[string]struct{})
	for _, line := range lines {
		if strings.HasPrefix(line, UPDATE_FILE_PREFIX) {
			result[strings.TrimPrefix(line, UPDATE_FILE_PREFIX)] = struct{}{}
		}
		if strings.HasPrefix(line, DELETE_FILE_PREFIX) {
			result[strings.TrimPrefix(line, DELETE_FILE_PREFIX)] = struct{}{}
		}
	}
	files := make([]string, 0, len(result))
	for file := range result {
		files = append(files, file)
	}
	return files
}

func identifyFilesAdded(text string) []string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	result := make(map[string]struct{})
	for _, line := range lines {
		if strings.HasPrefix(line, ADD_FILE_PREFIX) {
			result[strings.TrimPrefix(line, ADD_FILE_PREFIX)] = struct{}{}
		}
	}
	files := make([]string, 0, len(result))
	for file := range result {
		files = append(files, file)
	}
	return files
}

func loadFiles(paths []string, openFn func(string) (string, error)) (map[string]string, error) {
	orig := make(map[string]string)
	for _, p := range paths {
		content, err := openFn(p)
		if err != nil {
			return nil, &DiffError{Message: fmt.Sprintf("File not found: %s", p)}
		}
		orig[p] = content
	}
	return orig, nil
}

func applyCommit(commit Commit, writeFn func(string, string) error, removeFn func(string) error) error {
	for p, change := range commit.Changes {
		switch change.Type {
		case DELETE:
			if err := removeFn(p); err != nil {
				return err
			}
		case ADD:
			if err := writeFn(p, *change.NewContent); err != nil {
				return err
			}
		case UPDATE:
			if change.MovePath != nil {
				if err := writeFn(*change.MovePath, *change.NewContent); err != nil {
					return err
				}
				if err := removeFn(p); err != nil {
					return err
				}
			} else {
				if err := writeFn(p, *change.NewContent); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func processPatch(text string, openFn func(string) (string, error), writeFn func(string, string) error, removeFn func(string) error) (string, error) {
	if !strings.HasPrefix(text, PATCH_PREFIX) {
		return "", &DiffError{Message: "Patch must start with *** Begin Patch\\n"}
	}
	paths := identifyFilesNeeded(text)
	orig, err := loadFiles(paths, openFn)
	if err != nil {
		return "", err
	}
	patch, _ := textToPatch(text, orig)
	commit := patchToCommit(patch, orig)
	if err := applyCommit(commit, writeFn, removeFn); err != nil {
		return "", err
	}
	return "Done!", nil
}

func openFile(p string) (string, error) {
	content, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func writeFile(p string, content string) error {
	if filepath.IsAbs(p) {
		return &DiffError{Message: "We do not support absolute paths."}
	}
	parent := filepath.Dir(p)
	if parent != "." {
		if err := os.MkdirAll(parent, 0755); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(p, []byte(content), 0644)
}

func removeFile(p string) error {
	return os.Remove(p)
}

func main() {
	var patchText string
	fmt.Scanln(&patchText)
	if patchText == "" {
		fmt.Println("Please pass patch text through stdin")
		os.Exit(1)
	}
	result, err := processPatch(patchText, openFile, writeFile, removeFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result)
}
