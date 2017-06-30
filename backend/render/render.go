package render

import (
	"html/template"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/bigpigeon/easy-takeout/backend/logger"
)

type DirSorted []string

type PathFile struct {
	Path string
	File string
}

func CreateExcludeSorted(list []string) DirSorted {
	new_l := []string{}
	for _, l := range list {

		if strings.HasSuffix(l, "/") {
			new_l = append(new_l, l[:len(l)-1])
		} else {
			new_l = append(new_l, l)
		}
	}
	sort.Strings(new_l)
	return new_l
}

func (e DirSorted) search(s string) bool {
	start, end := 0, len(e)
	for start < end {
		half := (end - start) / 2
		if e[half] == s {
			return true
		} else if e[half] > s {
			end = half
		} else {
			start = half + 1
		}
	}
	return false
}

/*
* split file with need rending and not
 */
func SplitRenderFile(root string, keepdir DirSorted) ([]PathFile, []PathFile, error) {
	return _splitRenderFile(root, "", keepdir)
}

/*
* Render render html file from source path to target path
* use html/template lib to parse template to the specified data object
* any have ".html" suffix file in source path will be load
* any without "_" prefix and have ".html" suffix file in source path will be render to target path with same path name
* all in keepdirs list path file/dir will be copy to target path
 */
func Render(source, target string, keepdirs []string, data interface{}) error {
	pathfiles, keepfiles, err := SplitRenderFile(source, CreateExcludeSorted(keepdirs))
	if err != nil {
		return err
	}

	// copy all keep files to target path
	for _, f := range keepfiles {
		fSource := path.Join(source, f.Path, f.File)
		fTarget := path.Join(target, f.Path, f.File)
		logger.Debug.Println(fSource, fTarget)
		dir, _ := path.Split(fTarget)
		err := mkdirIfNotExist(dir)
		if err != nil {
			return err
		}
		err = copyFile(fSource, fTarget)
		if err != nil {
			return err
		}
	}
	// collect *.html ,render html file that without "_" prefix to target path
	commonPages := []string{}
	renderPages := []string{}
	for _, pf := range pathfiles {
		if strings.HasSuffix(pf.File, ".html") {
			if strings.HasPrefix(pf.File, "_") {
				commonPages = append(commonPages, path.Join(source, pf.Path, pf.File))
			} else {
				renderPages = append(renderPages, path.Join(pf.Path, pf.File))
			}
		}
	}
	for _, r := range renderPages {
		rSource := path.Join(source, r)
		rTarget := path.Join(target, r)
		temp, err := template.ParseFiles(append([]string{rSource}, commonPages...)...)
		if err != nil {
			return err
		}
		// mkdir target path if not exist
		dir, _ := path.Split(rTarget)
		err = mkdirIfNotExist(dir)
		if err != nil {
			return err
		}
		// write data to target file
		output, err := os.Create(rTarget)
		defer output.Close()
		if err != nil {
			return err
		}
		err = temp.Execute(output, data)
		if err != nil {
			return err
		}
	}
	return nil
}
