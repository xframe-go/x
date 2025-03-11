package x

import "path/filepath"

func PublicPath(abs ...bool) string {
	path := filepath.Join("public")
	if len(abs) > 0 && abs[0] {
		dir, _ := filepath.Abs(path)
		return dir
	}
	return path
}
