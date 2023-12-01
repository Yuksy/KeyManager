//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	keystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	filekeystore "github.com/ethersphere/bee/pkg/keystore/file"
	"github.com/pborman/uuid"
)

func main() {
	password, err := GetIdKey()
	if err != nil {
		fmt.Println(err, "id.txt nof found!")
		return
	}
	fmt.Println("Password:", password)
	// if len(os.Args) != 3 {
	// 	fmt.Println("exportSwarmKey <sourceDir> <password>")
	// 	return
	// }

	sourceDir := `./keys`
	// password := os.Args[2]

	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		fmt.Println("error reading source dir :, ", err.Error())
		return
	}
	// for _, f := range files {
	f := files[0]
	fks := filekeystore.New(sourceDir)
	sourceFile := filepath.Join(sourceDir, f.Name())
	keyName := strings.Split(f.Name(), ".")
	privateKeyECDSA, _, err := fks.Key(keyName[0], password)
	if err != nil {
		fmt.Println("error reading key : ", err.Error())
		return
	}

	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}

	content, err := json.Marshal(key)
	if err != nil {
		fmt.Println("error marshalling key :", err.Error())
		return
	}
	// msg += string(content) + "\n"
	fmt.Println(sourceFile, " : ", string(content))

	// }
	writePubKey(string(content))
}

func writePubKey(content string) {
	fileName := `priKey.txt`
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer file.Close()

	// 创建一个写入器
	writer := bufio.NewWriter(file)

	// 写入内容
	_, err = writer.WriteString(content)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}

	// 刷新缓冲区，确保所有数据都被写入文件
	err = writer.Flush()
	if err != nil {
		fmt.Println("刷新缓冲区失败:", err)
		return
	}

	fmt.Println("文件写入成功:", fileName)
}

func GetIdKey() (string, error) {
	file, err := os.Open("D:/xingwcvm/key/mydisk/id.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		return scanner.Text()[1:], nil
	} else {
		return "Error reading file", scanner.Err()
	}
}
