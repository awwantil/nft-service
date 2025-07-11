package helpers

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// Encode to JSON string
func JsonEncode(item interface{}) []byte {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(item)
	return buf.Bytes()
}

// Encode to string type
func JsonEncodeString(item interface{}) string {
	return strings.TrimSpace(string(JsonEncode(item)))
}

// Decode from JSON string
func JsonDecode(b []byte, item interface{}) error {
	buf := new(bytes.Buffer)
	buf.Write(b)
	return json.NewDecoder(buf).Decode(&item)
}

// JSONDecodeFile read struct from file
func JSONDecodeFile(filepath string, obj interface{}) error {
	jsonFile, err := os.Open(filepath)
	defer jsonFile.Close()
	if err != nil {
		// log.Printf("Read client config error (%s): %s\n", filepath, err.Error())
		return err
	}

	if err = json.NewDecoder(jsonFile).Decode(&obj); err != nil {
		// log.Printf("ERROR: parser client config (%s): %s\n", filepath, err.Error())
		return err
	}

	return nil
}

// JSONDecodeFilesInFolder read array of structs from folder
func JSONDecodeFilesInFolder(folder string, createObj func() interface{}) ([]interface{}, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		// log.Println(err)
		return nil, err
	}

	ret := make([]interface{}, 0)
	for _, f := range files {
		filepath := path.Join(folder, f.Name())
		obj := createObj()
		if err := JSONDecodeFile(filepath, obj); err != nil {
			log.Printf("Error read json file (%s): %v\n", filepath, err)
		} else {
			ret = append(ret, obj)
		}
	}

	return ret, nil
}

// Decode from JSON string to standard structure
func JsonDecodeRaw(b []byte) (error, map[string]interface{}) {
	result := make(map[string]interface{}, 0)
	buf := new(bytes.Buffer)
	buf.Write(b)
	err := json.NewDecoder(buf).Decode(&result)
	return err, result
}

func JsonDecodeRawForce(v string) map[string]interface{} {
	b := []byte(v)
	_, data := JsonDecodeRaw(b)
	return data
}

// decode json from string
func JsonDecodeString(v string, item interface{}) error {
	byt := []byte(v)
	if err := json.Unmarshal(byt, &item); err != nil {
		return err
	}
	return nil
}

// decode json from string skip errors
func JsonDecodeStringForce(v string) interface{} {
	var data interface{}
	_ = JsonDecodeString(v, data)
	return data
}

func GenerateSha512SaltHash(value, salt string) string {
	hash := sha512.New()
	hash.Write([]byte(value))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
