package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-runewidth"
)

func main() {
	flag.Parse()
	os.Exit(run())
}

func run() int {
	// コマンドライン引数をファイルパスとして扱う
	files := os.Args[1:]

	// 出力ディレクトリ
	if _, err := os.Stat("out"); os.IsNotExist(err) {
		if err := os.MkdirAll("out", 0755); err != nil {
			fmt.Println("Error creating directory:", err)
			return 1
		}
	}

	// 出力ファイル
	outfiles := make(map[string]*bufio.Writer)

	for k := range outfiles {
		f, err := os.Create(k)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return 1
		}
		defer f.Close()
		outfiles[k] = bufio.NewWriter(f)
		defer outfiles[k].Flush()
	}

	for _, f := range files {
		// 入力ファイル
		infile, err := os.Open(f)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return 1
		}
		defer infile.Close()

		reader := bufio.NewReader(infile)
		for {
			r, size, err := reader.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("Invalid unicode: [%s] %v %s\n", string(r), r, err)
				return 1
			}

			width := runewidth.RuneWidth(r)

			fname := fmt.Sprintf("out/%db%dw.txt", size, width)
			if _, ok := outfiles[fname]; !ok {

				f, err := os.Create(fname)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return 1
				}
				defer f.Close()
				outfiles[fname] = bufio.NewWriter(f)
				defer outfiles[fname].Flush()
			}

			size, err = outfiles[fname].WriteRune(r)
			if err != nil {
				fmt.Printf("Write error: %s\n", err)
				return 1
			}

		}

	}
	return 0
}
