package agent

import (
	"strings"
)

type ApplyPatchCreateFileOp struct {
	Type    string
	Path    string
	Content string
}

type ApplyPatchDeleteFileOp struct {
	Type string
	Path string
}

type ApplyPatchUpdateFileOp struct {
	Type   string
	Path   string
	Update string
	Added  int
	Deleted int
}

type ApplyPatchOp interface{}

const (
	PATCH_PREFIX       = "*** Begin Patch\n"
	PATCH_SUFFIX       = "\n*** End Patch"
	ADD_FILE_PREFIX    = "*** Add File: "
	DELETE_FILE_PREFIX = "*** Delete File: "
	UPDATE_FILE_PREFIX = "*** Update File: "
	END_OF_FILE_PREFIX = "*** End of File"
	HUNK_ADD_LINE_PREFIX = "+"
)

func ParseApplyPatch(patch string) []ApplyPatchOp {
	if !strings.HasPrefix(patch, PATCH_PREFIX) {
		return nil
	} else if !strings.HasSuffix(patch, PATCH_SUFFIX) {
		return nil
	}

	patchBody := patch[len(PATCH_PREFIX) : len(patch)-len(PATCH_SUFFIX)]
	lines := strings.Split(patchBody, "\n")

	var ops []ApplyPatchOp

	for _, line := range lines {
		if strings.HasPrefix(line, END_OF_FILE_PREFIX) {
			continue
		} else if strings.HasPrefix(line, ADD_FILE_PREFIX) {
			ops = append(ops, ApplyPatchCreateFileOp{
				Type:    "create",
				Path:    strings.TrimPrefix(line, ADD_FILE_PREFIX),
				Content: "",
			})
			continue
		} else if strings.HasPrefix(line, DELETE_FILE_PREFIX) {
			ops = append(ops, ApplyPatchDeleteFileOp{
				Type: "delete",
				Path: strings.TrimPrefix(line, DELETE_FILE_PREFIX),
			})
			continue
		} else if strings.HasPrefix(line, UPDATE_FILE_PREFIX) {
			ops = append(ops, ApplyPatchUpdateFileOp{
				Type:   "update",
				Path:   strings.TrimPrefix(line, UPDATE_FILE_PREFIX),
				Update: "",
				Added:  0,
				Deleted: 0,
			})
			continue
		}

		lastOp := ops[len(ops)-1]

		switch op := lastOp.(type) {
		case ApplyPatchCreateFileOp:
			op.Content = appendLine(op.Content, line[len(HUNK_ADD_LINE_PREFIX):])
			ops[len(ops)-1] = op
		case ApplyPatchUpdateFileOp:
			if strings.HasPrefix(line, HUNK_ADD_LINE_PREFIX) {
				op.Added++
			} else if strings.HasPrefix(line, "-") {
				op.Deleted++
			}
			op.Update = appendLine(op.Update, line)
			ops[len(ops)-1] = op
		}
	}

	return ops
}

func appendLine(content, line string) string {
	if content == "" {
		return line
	}
	return content + "\n" + line
}
