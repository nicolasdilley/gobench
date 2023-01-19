package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	// filepath.Walk("./blocking", func(path string, info os.FileInfo, err error) error {

	// 	if strings.Contains(info.Name(), ".go") {

	// 		content, _ := ioutil.ReadFile(path)

	// 		new_content := ""

	// 		for _, line := range strings.Split(string(content), "\n") {
	// 			new_line := line

	// 			if strings.Contains(new_line, "package ") {
	// 				new_line = "package main"
	// 			} else if strings.Contains(new_line, "func main()") {
	// 				new_line = "func main() {"
	// 			} else if strings.Contains(new_line, "\"testing\"") {
	// 				new_line = ""
	// 			}

	// 			new_content += new_line + "\n"
	// 		}

	// 		os.Remove(path)

	// 		ioutil.WriteFile(path, []byte(new_content), 0644)

	// 	}

	// 	return nil
	// })
	pwd, _ := filepath.Abs(".")
	filepath.Walk("./blocking", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(path, ".go") {

			path, _ := filepath.Abs(path)

			fmt.Println(filepath.Dir(path))
			// run migoInfer
			run := exec.Command(
				"docker",
				"run",
				"-i",
				"--rm",
				"-v",
				filepath.Dir(path)+":/root",
				"jgabet/godel2:latest",
				"migoinfer", filepath.Base(path),
			)

			fl, _ := os.Create("example.cgo")

			writer := bufio.NewWriter(fl)
			run.Stdout = writer
			run.Dir = filepath.Dir(path)
			// run.Stderr = os.Stderr
			run.Run()
			writer.Flush()

			// run godel2
			godel2 := exec.Command("timeout", "20", "docker", "run", "-i", "--rm", "-v",
				pwd+":/root",
				"jgabet/godel2:latest",
				"/usr/bin/Godel", "example.cgo")
			// godel2.Dir = filepath.Dir(path)
			conten, _ := godel2.Output()

			godel2.Run()
			// parse result

			fmt.Println(string(conten))
			if strings.Contains(string(conten), "No global deadlock:") {
			}
			fl.Close()
			os.Remove("example.cgo")

		}

		return nil
	})

}
