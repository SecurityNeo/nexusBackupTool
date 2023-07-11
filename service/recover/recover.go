package recover

/*
import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/SecurityNeo/NexusbackupTool/common"
	"github.com/SecurityNeo/NexusbackupTool/service/backup"
	"github.com/SecurityNeo/NexusbackupTool/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const DefaultUserAgent = "NexusTool/1.0"
const applicationJSON = "application/json"
const ContentType = "application/json"

func StartBackup(configDirectory string, parallelism int) error {
	var confirm string

	JsonParse := backup.NewJsonStruct()
	config := types.TargetConfig{}
	JsonParse.Load(configDirectory, &config)

	fmt.Printf("Please Confirm Your Configuration:\n")
	fmt.Printf("--------------------------------------\n")
	fmt.Printf("Configuration File: %v\n", configDirectory)
	fmt.Printf("Target Platform Configuration:\n")
	fmt.Printf("    IP:        %v\n", config.Platform.Ip)
	fmt.Printf("    Port:      %v\n", config.Platform.Port)
	fmt.Printf("target Nexus Configuration:\n")
	fmt.Printf("    Repo:      %v\n", config.NexusConfig.Repository)
	fmt.Printf("    Username:  %v\n", config.NexusConfig.Username)
	fmt.Printf("    password:  %v\n", config.NexusConfig.Password)
	fmt.Printf("Datasource Directory: %v\n", config.Backup.Directory)
	fmt.Printf("--------------------------------------\n")
	fmt.Printf("Please Enter yes To Continue:")
	fmt.Scanln(&confirm)
	if confirm != "yes" {
		log.Println("Cancel!")
		os.Exit(10)
	}

	url := "http://" + config.Platform.Ip + ":" + config.Platform.Port
	repository := config.NexusConfig.Repository
	log.Println("Check if the repository ", repository, "exists.")
	checkRepoUrl := getRepoURL(url)

	exists, err := checkRepo(checkRepoUrl, repository)
	if err != nil {
		log.Fatalln("Check repo error: ", err)
		os.Exit(10)
	}
	if exists {
		log.Println("Repository ", repository, "exists in the nexus, no need to create it.")
	} else {
		log.Println("Repository ", repository, "not exists in the nexus,please create it.")
		os.Exit(10)
	}

	uploadUrl := uploadURL(url, repository)

	return nil
}

func checkRepo(Url string, Repo string) (exists bool, err error) {
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Accept", applicationJSON)
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("User-Agent", DefaultUserAgent)

	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	var ok bool
	if resp.StatusCode == 200 {
		ok = true
	}

	if !ok {
		resp.Body.Close()
		return false, err
	}
	resBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	result := Repos{}
	json.Unmarshal(resBody, result)
	for _, v := range result {
		if v.Name == Repo && v.Type == "hosted" && v.Format == "raw" {
			return true, nil
		}
	}
	return false, nil
}

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func upload(url, filePath string) {

	ok, err := common.PathExists(filePath)
	if err != nil {
		log.Fatalln("Check File Error:", err)
		return
	}
	if !ok {
		log.Fatalln("File ", filePath, "not exists,skip it!")
		return
	}

	body := common.NewCircleByteBuffer(1024 * 2)
	boundary := randomBoundary()
	boundaryBytes := []byte("\r\n--" + boundary + "\r\n")
	endBytes := []byte("\r\n--" + boundary + "--\r\n")

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "multipart/form-data; charset=utf-8; boundary="+boundary)
	go func() {
		//defer ruisRecovers("upload.run")
		f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
		if err != nil {
			panic(err)
		}
		stat, err := f.Stat()
		if err != nil {
			panic(err)
		}
		defer f.Close()
		path, _ := filepath.Abs(filepath.Dir(filePath))
		header := fmt.Sprintf("Content-Disposition: form-data; raw.directory=\"%s\"; raw.asset1.filename=\"%s\"\r\nContent-Type: application/octet-stream\r\n\r\n", path, stat.Name())
		body.Write(boundaryBytes)
		body.Write([]byte(header))

		fsz := float64(stat.Size())
		fupsz := float64(0)
		buf := make([]byte, 1024)
		for {
			time.Sleep(10 * time.Microsecond)
			n, err := f.Read(buf)
			if n > 0 {
				nz, _ := body.Write(buf[0:n])
				fupsz += float64(nz)
				progress := strconv.Itoa(int((fupsz/fsz)*100)) + "%"
				fmt.Println("upload:", progress, "|", strconv.FormatFloat(fupsz, 'f', 0, 64), "/", stat.Size())
			}
			if err == io.EOF {
				break
			}
		}
		body.Write(endBytes)
		body.Write(nil)
	}()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		fmt.Println("上传成功")
	} else {
		fmt.Println("上传失败,StatusCode:", resp.StatusCode)
	}
}
*/
