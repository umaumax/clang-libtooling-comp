package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
)

var readTest bool

func init() {
	flag.BoolVar(&readTest, "read-test", false, "read test flag")
}

func main() {
	flag.Parse()

	dataPath, err := homedir.Expand("~/.config/clang-libtooling-comp/data.json")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("datpath:", dataPath)
	err = os.MkdirAll(filepath.Dir(dataPath), 0755)
	if err != nil {
		log.Fatalln(err)
	}
	datas := make(map[string][]string)

	if readTest {
		// NOTE: read test
		bytes, err := ioutil.ReadFile(dataPath)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(bytes, &datas); err != nil {
			log.Fatal(err)
		}
		fmt.Println(datas)
		return
	}

	targetDir := "results/"
	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".comp" {
			return nil
		}
		filename := filepath.Base(path)
		namespace := strings.TrimSuffix(strings.Replace(filename, "--", "::", -1), ".cpp.comp")
		// 		fmt.Println(namespace)
		var bytes []byte
		if bytes, err = ioutil.ReadFile(path); err != nil {
			return err
		}

		for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(bytes), -1) {
			if v != "" {
				datas[namespace] = append(datas[namespace], v)
			}
			// 			fmt.Println(i+1, ":", v)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	// 	return

	// NOTE: create dummy data
	// 	for i := 0; i < 10; i++ {
	// 		key := fmt.Sprint(i)
	// 		datas[key] = []string{}
	// 		for j := 0; j < i; j++ {
	// 			datas[key] = append(datas[key], key)
	//
	// 		}
	// 	}
	// 	fmt.Println(datas)
	j, err := json.Marshal(datas)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(dataPath, j, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
