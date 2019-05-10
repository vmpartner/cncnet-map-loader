package main

import (
	"archive/zip"
	"fmt"
	brut "github.com/dieyushi/golang-brutedict"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const folderMaps = "yr_maps"
const folderTmp = "yr_tmp"

func main() {

	os.RemoveAll(folderTmp)
	os.MkdirAll(folderMaps, 0777)
	os.MkdirAll(folderTmp, 0777)

	bd := brut.New(false, true, false, 3, 3)
	defer bd.Close()
	for {

		id := bd.Id()
		fmt.Println(id)
		if id == "zzz" {
			break
		}

		urlSearch := "http://mapdb.cncnet.org/search.php?game=yr&age=0&search=" + id
		fmt.Println("SEARCH " + urlSearch)
		resp, err := http.Get(urlSearch)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		html := string(body)
		html = strings.ToLower(html)

		re, err := regexp.Compile(`<a href="\.(.*?)"`)
		if err != nil {
			panic(err)
		}
		res := re.FindAllStringSubmatch(html, -1)

		fmt.Println("FOUND " + strconv.Itoa(len(res)))

		for _, rs := range res {

			fileExist := strings.Replace(rs[1], ".zip", ".map", -1)
			fileExist = path.Base(fileExist)
			file, _ := os.Stat(folderMaps + "/" + fileExist)
			if file != nil {
				continue
			}

			url := "http://mapdb.cncnet.org" + rs[1]
			fmt.Println("MAP " + url)

			// Get the data
			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			fileZip := folderTmp + "/" + path.Base(url)
			out, err := os.Create(fileZip)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				panic(err)
			}
			out.Close()
			resp.Body.Close()

			files, err := Unzip(fileZip, folderMaps)
			if err != nil {
				panic(err)
			}

			for _, file := range files {
				fileNew := strings.Replace(file, ".mpr", ".map", -1)
				err := os.Rename(file, fileNew)
				if err != nil {
					panic(err)
				}
			}

			os.RemoveAll(fileZip)
		}
	}
}

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
