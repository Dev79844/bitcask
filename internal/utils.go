package internal

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func getFiles(dir string) ([]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.db", dir))
	if err != nil {
		return nil, err
	}
	return files, nil
}

func getIDs(files []string) ([]int, error) {
	ids := make([]int, 0)

	for _, f := range files {
		id, err := strconv.ParseInt(strings.TrimPrefix(strings.TrimSuffix(filepath.Base(f), ".db"), "bitcask_"), 10, 32)
		if err!=nil{
			return nil, err
		}
		ids = append(ids, int(id))
	}

	sort.Ints(ids)

	return ids, nil
}
