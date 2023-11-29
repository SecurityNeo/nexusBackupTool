package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SecurityNeo/NexusbackupTool/common"
	"github.com/SecurityNeo/NexusbackupTool/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const DefaultUserAgent = "NexusTool/1.0"
const applicationJSON = "application/json"
const ContentType = "application/json"

type JsonStruct struct {
}

func StartBackup(configDirectory string, parallelism int) error {
	var confirm string

	JsonParse := NewJsonStruct()
	config := types.SrcConfig{}
	JsonParse.Load(configDirectory, &config)

	fmt.Printf("Please Confirm Your Configuration:\n")
	fmt.Printf("--------------------------------------\n")
	fmt.Printf("Configuration File: %v\n", configDirectory)
	fmt.Printf("Platform Configuration:\n")
	fmt.Printf("    IP:        %v\n", config.Platform.Ip)
	fmt.Printf("    Port:      %v\n", config.Platform.Port)
	fmt.Printf("Nexus Configuration:\n")
	fmt.Printf("    Repo:      %v\n", config.NexusConfig.Repository)
	fmt.Printf("    Username:  %v\n", config.NexusConfig.Username)
	fmt.Printf("    password:  %v\n", config.NexusConfig.Password)
	fmt.Printf("Backup Directory: %v\n", config.Backup.Directory)
	fmt.Printf("--------------------------------------\n")
	fmt.Printf("Please Enter yes To Continue:")
	fmt.Scanln(&confirm)
	if confirm != "yes" {
		log.Println("Cancel!")
		os.Exit(10)
	}
	url := "http://" + config.Platform.Ip + ":" + config.Platform.Port
	getUrl := getURL(url, config.NexusConfig.Repository)

	files, err := getFileData(getUrl)
	if err != nil {
		log.Fatalln("Get Files Error: ", err)
		return nil
	}

	DI := types.DownloadInfo{
		Files: files,
		Chs:   make(chan int, parallelism),
		Ans:   make(chan bool),
	}
	fmt.Printf("Remote File Info:\n")
	for _, v := range DI.Files {
		fmt.Printf("    Download URL: %v\n", v.DownloadUrl)
		fmt.Printf("    MD5: %v\n", v.MD5)
	}
	file := config.Backup.Directory + "info.json"

	files2json, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		log.Fatalln("Marshal", files, "failed! Error: ", err)
		os.Exit(10)
	}

	exists, err := common.PathExists(config.Backup.Directory)
	if err != nil {
		log.Fatalln("Check File Error: ", err)
		return nil
	}
	if !exists {
		os.MkdirAll(config.Backup.Directory, 644)
	}
	common.WriteConfig(file, files2json)
	go func() {
		for _, v := range DI.Files {
			DI.Chs <- 1
			DI.Wg.Add(1)
			go DI.Work(v.DownloadUrl, config.Backup.Directory+v.Path, v.MD5)
		}
		DI.Wg.Wait()
		close(DI.Ans)
	}()

	for _ = range DI.Ans {
	}
	log.Println("All files has been saved to", config.Backup.Directory)
	return nil
}

func getFileData(Url string) (types.FileData, error) {
	var fileData types.FileData
	var components *types.Components
	components, err := Request("GET", Url, "", "")
	if err != nil {
		return nil, err
	}
	var files []map[string]string
	for _, value := range components.Items {
		v := map[string]string{
			"downloadUrl": value.Assets[0].DownloadUrl,
			"path":        value.Assets[0].Path,
			"md5":         value.Assets[0].Check.MD5,
		}
		files = append(files, v)
	}

	token := components.ContinuationToken

	for {
		if token == "" {
			break
		} else {
			components, err := Request("GET", Url, token, "")
			if err != nil {
				return nil, err
			}
			if components == nil {
				break
			}
			token = components.ContinuationToken

			for _, value := range components.Items {
				v := map[string]string{
					"downloadUrl": value.Assets[0].DownloadUrl,
					"path":        value.Assets[0].Path,
					"md5":         value.Assets[0].Check.MD5,
				}
				files = append(files, v)
			}
		}
	}

	map2json, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		log.Fatalln("Marshal Error: ", err)
		os.Exit(10)
	}

	errs := json.Unmarshal([]byte(string(map2json)), &fileData)
	if errs != nil {
		log.Fatalln("Marshal Error: ", errs)
		os.Exit(10)
	}

	return fileData, nil

}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

func Request(method string, url string, continueToken string, RequestBody interface{}) (*types.Components, error) {
	var URL string
	body, err := json.Marshal(RequestBody)
	if err != nil {
		return nil, err
	}
	if continueToken != "" {
		URL = url + "&continuationToken=" + continueToken
	} else {
		URL = url
	}

	req, err := http.NewRequest(method, URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", applicationJSON)
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("User-Agent", DefaultUserAgent)

	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var ok bool
	if resp.StatusCode == 200 {
		ok = true
	}

	if !ok {
		resp.Body.Close()
		return nil, nil
	}
	resBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	result := &types.Components{}
	json.Unmarshal(resBody, result)
	return result, nil
}
