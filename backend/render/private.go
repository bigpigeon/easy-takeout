package render

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func mkdirIfNotExist(dir string) error {
	// mkdir target path if not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	return nil
}

func getAllFiles(root, curr string) ([]PathFile, error) {
	absPath := path.Join(root, curr)
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return []PathFile{}, err
	}
	pathfiles := []PathFile{}
	for _, f := range files {
		if f.IsDir() {
			subpathfiles, err := getAllFiles(root, path.Join(curr, f.Name()))
			if err != nil {
				return []PathFile{}, err
			}

			pathfiles = append(pathfiles, subpathfiles...)

		} else {
			// set path with relative position
			pathfiles = append(pathfiles, PathFile{curr, f.Name()})
		}

	}
	return pathfiles, nil
}

func _splitRenderFile(root, curr string, keepdir DirSorted) ([]PathFile, []PathFile, error) {
	absPath := path.Join(root, curr)
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return []PathFile{}, []PathFile{}, err
	}
	renderFiles := []PathFile{}
	keepFiles := []PathFile{}
	for _, f := range files {
		if keepdir.search(f.Name()) == true {
			files, err := getAllFiles(root, path.Join(curr, f.Name()))
			if err != nil {
				return []PathFile{}, []PathFile{}, err
			}
			keepFiles = append(keepFiles, files...)
		} else {
			if f.IsDir() { // recusion director
				subNotRending := DirSorted{}
				subkeepdir := f.Name() + "/"
				for _, e := range keepdir {

					if strings.HasPrefix(e, subkeepdir) && len(e) > len(subkeepdir) {
						subNotRending = append(subNotRending, e[len(subkeepdir):])
					}
				}
				subRenderFiles, subKeepFiles, err := _splitRenderFile(root, path.Join(curr, f.Name()), subNotRending)
				if err != nil {
					return []PathFile{}, []PathFile{}, err
				}
				// add prefix path to sub html file
				renderFiles = append(renderFiles, subRenderFiles...)
				keepFiles = append(keepFiles, subKeepFiles...)
			} else {
				renderFiles = append(renderFiles, PathFile{curr, f.Name()})
			}
		}
	}
	return renderFiles, keepFiles, nil
}
