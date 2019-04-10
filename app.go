package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	dir := "/Users/Derek/Desktop/tmp/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	regex := *regexp.MustCompile(`(?P<month>(January|February|March|April|May|June|July|August|September|October|November|December))\s+(?P<day>\d{1,2}),\s+(?P<year>\d{4})`)

	for _, f := range files {
 		if f.IsDir() {
			res := regex.FindAllStringSubmatch(f.Name(), -1)

			if len(res) == 0 {
				fmt.Println(fmt.Sprintf("[X] %s =====> skip", f.Name()))
				continue
			}

			of := f.Name()
			nf := fmtDirName(res[0][4], res[0][2], res[0][3])
			o := fmt.Sprintf("%s%s/", dir, of)
			n := fmt.Sprintf("%s%s/", dir, nf)

			if _, err := os.Stat(n); err == nil {
				children, err := ioutil.ReadDir(o)
				if err != nil {
					panic(err)
				}

				fmt.Println(fmt.Sprintf("[M] %s =====> %s", of, nf))
				for _, c := range children {
					cof := fmt.Sprintf("%s%s", o, c.Name())
					cnf := fmt.Sprintf("%s%s", n, c.Name())

					err := os.Rename(cof, cnf)

					if err != nil {
						panic(err)
					}
				}

				os.Rename(o, fmt.Sprintf("%s%s/", dir, "_" + of))
			} else {
				fmt.Println(fmt.Sprintf("[O] %s =====> %s", of, nf))
				err := os.Rename(o, n)

				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func fmtDirName(year string, month string, day string) string {
	var m string
	var d string

	switch month {
	case "January":
		m = "01"
		break
	case "February":
		m = "02"
		break
	case "March":
		m = "03"
		break
	case "April":
		m = "04"
		break
	case "May":
		m = "05"
		break
	case "June":
		m = "06"
		break
	case "July":
		m = "07"
		break
	case "August":
		m = "08"
		break
	case "September":
		m = "09"
		break
	case "October":
		m = "10"
		break
	case "November":
		m = "11"
		break
	case "December":
		m = "12"
		break
	}

	if len(day) == 1 {
		d = "0" + day
	} else {
		d = day
	}

	return fmt.Sprintf("%s%s%s", year, m, d)
}

func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}