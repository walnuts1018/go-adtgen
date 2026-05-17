package writer

import "os"

func WriteFile(path string, src string) error {
	return os.WriteFile(path, []byte(src), 0o644)
}
