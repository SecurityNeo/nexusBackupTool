package common

import (
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type CircleByteBuffer struct {
	io.Reader
	io.Writer
	io.Closer
	datas []byte

	start   int
	end     int
	size    int
	isClose bool
	isEnd   bool
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Md5Sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	r := bufio.NewReader(f)

	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func Download(downloadUrl string, targetPath string, md5target string) {

	ok, err := PathExists(targetPath)
	if err != nil {
		log.Fatalln("Check File Error:", err)
		return
	}

	if ok {
		md5sum, err := Md5Sum(targetPath)
		if err != nil {
			log.Fatalln("Md5Error: ", err)
			return
		}
		if md5sum == md5target {
			log.Println(targetPath, "exists,MD5: ", md5target, "Do not need to download it again!")
			return
		} else {
			log.Println("The MD5 of file", targetPath, "do not match the remote repository,need to update!")
			log.Println("Remote MD5: ", md5target, "Act MD5: ", md5sum)
			err := os.Remove(targetPath)
			if err != nil {
				log.Fatalln("Remove file", targetPath, "failed!skip!")
				return
			}
		}
	}

	path, err := filepath.Abs(filepath.Dir(targetPath))
	if err != nil {
		log.Fatalln("Get absolute path of file", targetPath, "failed!")
	}
	exists, err := PathExists(path)
	if err != nil {
		log.Fatalln("Check File Error: ", err)
		return
	}
	if !exists {
		os.MkdirAll(path, 644)
	}

	newFile, err := os.Create(targetPath)
	if err != nil {
		log.Fatalln("Create file", targetPath, "failed! Error Log: ", err)
		return
	}
	defer newFile.Close()

	log.Println("Downloading ", downloadUrl)

	client := http.Client{}
	resp, err := client.Get(downloadUrl)
	if err != nil {
		log.Fatalln("Download", downloadUrl, "Error: ", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func WriteConfig(cfg string, jsonByte []byte) {

	file, err := os.Create(cfg)
	defer file.Close()
	if err != nil {
		log.Fatalln("Create Info File Failed! Error: ", err)
		os.Exit(1)
	}
	_, errs := file.Write(jsonByte)

	if errs != nil {
		log.Fatalln("Write config file:", cfg, "fail:", err)
	}

	log.Println("Write config file:", cfg, "successfully")
}

func NewCircleByteBuffer(len int) *CircleByteBuffer {
	var e = new(CircleByteBuffer)
	e.datas = make([]byte, len)
	e.start = 0
	e.end = 0
	e.size = len
	e.isClose = false
	e.isEnd = false
	return e
}

func (e *CircleByteBuffer) getLen() int {
	if e.start == e.end {
		return 0
	} else if e.start < e.end {
		return e.end - e.start
	} else {
		return e.start - e.end
	}
}
func (e *CircleByteBuffer) getFree() int {
	return e.size - e.getLen()
}
func (e *CircleByteBuffer) putByte(b byte) error {
	if e.isClose {
		return io.EOF
	}
	e.datas[e.end] = b
	var pos = e.end + 1
	for pos == e.start {
		if e.isClose {
			return io.EOF
		}
		time.Sleep(time.Microsecond)
	}
	if pos == e.size {
		e.end = 0
	} else {
		e.end = pos
	}
	return nil
}

func (e *CircleByteBuffer) getByte() (byte, error) {
	if e.isClose {
		return 0, io.EOF
	}
	if e.isEnd && e.getLen() <= 0 {
		return 0, io.EOF
	}
	if e.getLen() <= 0 {
		return 0, errors.New("no datas")
	}
	var ret = e.datas[e.start]
	e.start++
	if e.start == e.size {
		e.start = 0
	}
	return ret, nil
}

func (e *CircleByteBuffer) geti(i int) byte {
	if i >= e.getLen() {
		panic("out buffer")
	}
	var pos = e.start + i
	if pos >= e.size {
		pos -= e.size
	}
	return e.datas[pos]
}

/*
	func (e*CircleByteBuffer)puts(bts []byte){
		for i:=0;i<len(bts);i++{
			e.put(bts[i])
		}
	}

	func (e*CircleByteBuffer)gets(bts []byte)int{
		if bts==nil {return 0}
		var ret=0
		for i:=0;i<len(bts);i++{
			if e.getLen()<=0{break}
			bts[i]=e.get()
			ret++
		}
		return ret
	}
*/
func (e *CircleByteBuffer) Close() error {
	e.isClose = true
	return nil
}
func (e *CircleByteBuffer) Read(bts []byte) (int, error) {
	if e.isClose {
		return 0, io.EOF
	}
	if bts == nil {
		return 0, errors.New("bts is nil")
	}
	var ret = 0
	for i := 0; i < len(bts); i++ {
		b, err := e.getByte()
		if err != nil {
			if err == io.EOF {
				return ret, err
			}
			return ret, nil
		}
		bts[i] = b
		ret++
	}
	if e.isClose {
		return ret, io.EOF
	}
	return ret, nil
}
func (e *CircleByteBuffer) Write(bts []byte) (int, error) {
	if e.isClose {
		return 0, io.EOF
	}
	if bts == nil {
		e.isEnd = true
		return 0, io.EOF
	}
	var ret = 0
	for i := 0; i < len(bts); i++ {
		err := e.putByte(bts[i])
		if err != nil {
			fmt.Println("Write bts err:", err)
			return ret, err
		}
		ret++
	}
	if e.isClose {
		return ret, io.EOF
	}
	return ret, nil
}
