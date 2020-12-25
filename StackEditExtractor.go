package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//TODO make it look more like a composite pattern
type Node struct {
	fileName    string
	fileContent string
	parentId    string
}

//TODO would be better to have a tree struct
var filesMap = make(map[string]Node)

// output folder where your .md will be written
const PATH_OUTPUT_FOLDER = "target/"

// your StackEdit backup file name goes here
const JSON_INPUT_FILE = "source/StackEdit workspace.json"

//TODO make input and output const as arguments
func main() {
	marshaledData, err := readAndMarshalJsonFile()
	if err {
		return
	}

	buildTreeStruct(marshaledData)

	serializeTree()

	fmt.Println("num of files written:", len(filesMap))
}

func serializeTree() {
	// TODO i'm sure it can be done in 1 loop no need for this one
	for _, file := range filesMap {
		finalName := addParentToTheFileName(file)
		// if it's a file, not a folder
		if len(file.fileContent) > 0 {
			filenameAndPathFinal := filepath.Join(PATH_OUTPUT_FOLDER, finalName)
			//we create the folder if it does not exists
			folderPath := filepath.Dir(filenameAndPathFinal)
			_ = os.MkdirAll(folderPath, os.ModePerm)
			//we write the file
			err := ioutil.WriteFile(filenameAndPathFinal, []byte(file.fileContent), 0644)
			check(err)
		}
	}
}

func buildTreeStruct(marshaledData map[string]interface{}) {
	for key := range marshaledData {
		block := marshaledData[key].(map[string]interface{})
		id := block["id"].(string)
		if block["type"] == "content" {
			splitRes := strings.Split(key, "/")[0]
			foundFile, _ := filesMap[splitRes]
			text := block["text"].(string)
			foundFile.fileContent = text
			filesMap[splitRes] = foundFile
		}
		if block["type"] == "file" {
			var descNode = filesMap[id]
			folderName := block["name"]
			descNode.fileName = folderName.(string)
			parentId := block["parentId"]
			if parentId != nil && len(parentId.(string)) > 0 {
				descNode.parentId = parentId.(string)
			}
			filesMap[id] = descNode
		}
		if block["type"] == "folder" {
			var descNode = filesMap[id]
			folderName := block["name"]
			descNode.fileName = folderName.(string)
			parentId := block["parentId"]
			if parentId != nil && len(parentId.(string)) > 0 {
				descNode.parentId = parentId.(string)
			}
			filesMap[id] = descNode
		}
	}
}

func readAndMarshalJsonFile() (map[string]interface{}, bool) {
	data, err := ioutil.ReadFile(JSON_INPUT_FILE)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, true
	}
	//fmt.Println("Contents of file:", string(data))
	reader := strings.NewReader(string(data))
	dec := json.NewDecoder(reader)
	var marshaledData map[string]interface{}
	if err := dec.Decode(&marshaledData); err == io.EOF {
		//TODO this is never printed, see why
		fmt.Println("done")
	} else if err != nil {
		log.Fatal(err)
	}
	_ = json.Unmarshal(data, &marshaledData)
	return marshaledData, false
}

// TODO recursive fx
func addParentToTheFileName(file Node) string {
	for len(file.parentId) > 0 {
		parent := findParentNode(file.parentId)
		file.fileName = filepath.Join(parent.fileName, file.fileName)
		file.parentId = parent.parentId
	}
	return file.fileName
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// TODO not sure this is useful
func findParentNode(node string) Node {
	descNode := filesMap[node]
	return descNode
}
