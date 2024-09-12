package utils

import (
	"fmt"
	"os"
	"strings"
)

func AssertDir(file_path string) error {
	path_splited := strings.Split(file_path, "/")
	dirs := path_splited[0 : len(path_splited)-1]
	cur_path := ""
	for i := 0; i < len(dirs); i++ {
		if i == 0 {
			cur_path = dirs[0]
		} else {
			cur_path = fmt.Sprintf("%s/%s", cur_path, dirs[i])
		}
		_, err := os.Stat(cur_path)
		if os.IsNotExist(err) {
			if err := os.Mkdir(cur_path, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}
