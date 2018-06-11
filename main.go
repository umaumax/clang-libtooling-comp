package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/mitchellh/go-homedir"
)

func init() {
	_ = pp.Black
}

// TODO: error handling
func readInfo(r io.Reader) (line string, col int, err error) {
	scanner := bufio.NewScanner(r)
	n := 0
	for scanner.Scan() {
		n++
		if n == 1 {
			line = scanner.Text()
		}
		if n == 2 {
			col, _ = strconv.Atoi(scanner.Text())
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	return
}

func loadCompData(filename string) (datas map[string][]string, err error) {
	datas = make(map[string][]string)

	var bytes []byte
	bytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &datas); err != nil {
		return
	}
	return
}

func main() {
	flag.Parse()
	var fp *os.File
	var err error
	if flag.NArg() == 0 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}
	line, col, _ := readInfo(fp)

	// 	clangOutputMap := map[string][]string{
	// 		"clang::ast_matchers::MatchFinder::MatchResult": {
	// 			"COMPLETION: Context : [#clang::ASTContext *const#]Context",
	// 			"COMPLETION: MatchResult : MatchResult::",
	// 			"COMPLETION: Nodes : [#const BoundNodes#]Nodes",
	// 			"COMPLETION: SourceManager : [#clang::SourceManager *const#]SourceManager",
	// 		},
	// 	}
	dataPath, err := homedir.Expand("~/.config/clang-libtooling-comp/data.json")
	if err != nil {
		log.Fatalln(err)
	}
	clangOutputMap, err := loadCompData(dataPath)
	if err != nil {
		log.Fatalln(err)
	}
	fuzzyMap := make(map[string]string)
	fuzzyMap[""] = "clang"
	for clangType, _ := range clangOutputMap {
		keys := strings.Split(clangType, "::")
		for i, key := range keys {
			fuzzyMap[key] = strings.Join(keys[:i+1], "::")
		}
		// 		for i, key := range keys {
		// 			fuzzyMap[key] = clangType
		// 		}
	}
	n := len(line)
	if n < col-1 {
		// TODO: error handing
	}
	// 	pp.Println(fuzzyMap)

	line = line[:col-1]
	line = strings.TrimSpace(line)
	re := regexp.MustCompile(`(\w*)\W*$`)
	matches := re.FindStringSubmatch(line)
	keys := matches[1:]
	// 	pp.Println(clangOutputMap)
	// 	pp.Println(keys)
	// 	keys := strings.Split(line, "::")
	for _, key := range keys {
		if list, ok := clangOutputMap[fuzzyMap[key]]; ok {
			for i, v := range list {
				_ = i
				lists := strings.Split(v, ": ")
				// NOTE: 例外の理由を確認する
				// 				fmt.Println(i, v)
				// 				fmt.Println(lists)
				// 				fmt.Println(i, lists[1:])
				if len(lists) == 3 {
					fmt.Printf("%s____%s\n", strings.TrimSpace(lists[1]), lists[2])
				} else {
					// 					fmt.Println(i, v)
				}
				// 				fmt.Println()
			}
			break
		}
	}

	// 	fmt.Printf("%s %d", line, col)
}
