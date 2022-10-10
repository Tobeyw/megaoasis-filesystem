package api

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func ImagRename(src string, to string) error {

	files, err := ioutil.ReadDir(src)
	if err != nil {
		panic(err)
	} // 获取文件名组成的切片，并遍历打印每一个文件名
	for _, file := range files {
		name := file.Name()
		oldname := strings.Split(name, ".")[0]
		newfilename, err := BigIntBase64(oldname)
		if err != nil {
			return err
		}
		err = CopyAndRename(src+"/"+file.Name(), to+"/"+newfilename)
		if err != nil {
			return err
		}
		fmt.Println(name)
	}

	return nil
}

func BigIntBase64(num string) (string, error) {
	// string -> float64
	bigNum, err := strconv.Atoi(num)
	if err != nil {
		return "", err
	}
	var blen byte = 1
	if bigNum > 255 {
		blen = 2
	}
	bytenum, err := IntToBytesLittleEndian(bigNum, blen)
	if err != nil {
		return "", err
	}
	base64 := base64.StdEncoding.EncodeToString(bytenum)
	return base64, nil
}

func CopyAndRename(srcFilename string, distFilename string) error {
	//只读方式打开源文件
	sF, err1 := os.Open(srcFilename)
	if err1 != nil {
		fmt.Println("err1=", err1)
		return err1
	}
	defer sF.Close()
	out, err := os.Create(distFilename)
	if err != nil {
		return err
	}
	wt := bufio.NewWriter(out)
	defer out.Close()
	n, err := io.Copy(wt, sF)
	fmt.Println("copy write", n)
	if err != nil {
		panic(err)
	}
	wt.Flush()
	return nil
}

func IntToBytesLittleEndian(n int, bytesLength byte) ([]byte, error) {
	switch bytesLength {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 3:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		return bytesBuffer.Bytes()[0:3], nil
	case 4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	case 5:
		tmp := n
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		return bytesBuffer.Bytes()[0:5], nil
	case 6:
		tmp := n
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		return bytesBuffer.Bytes()[0:6], nil
	case 7:
		tmp := n
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		return bytesBuffer.Bytes()[0:7], nil
	case 8:
		tmp := n
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("IntToBytesLittleEndian b param is invaild")
}
