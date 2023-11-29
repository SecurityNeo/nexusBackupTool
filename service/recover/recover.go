package recover

import (
	"encoding/json"
	"fmt"
	"github.com/SecurityNeo/NexusbackupTool/service/backup"
	"github.com/SecurityNeo/NexusbackupTool/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const DefaultUserAgent = "NexusTool/1.0"
const applicationJSON = "application/json"
const ContentType = "application/json"

func StartRecover(configDirectory string, parallelism int) error {
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

	exists, format, err := checkRepo(checkRepoUrl, repository)
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

	infoFile := config.Backup.Directory + "info.json"
	fileInfo, err := readJSONFile(infoFile)
	if err != nil {
		log.Println(err)
		os.Exit(10)
	}

	ui := types.UploadInfo{
		Files: fileInfo,
		Chs:   make(chan int, parallelism),
		Ans:   make(chan bool),
	}

	go func() {
		for _, v := range ui.Files {
			ui.Chs <- 1
			ui.Wg.Add(1)
			go ui.Work(uploadUrl, config.Backup.Directory+v.Path, v.Path, config.NexusConfig.Username, config.NexusConfig.Password, format)
		}
		ui.Wg.Wait()
		close(ui.Ans)
	}()

	for _ = range ui.Ans {
	}
	log.Println("All files has been recover to ", url)

	return nil
}

func readJSONFile(filePath string) (types.FileData, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var data types.FileData
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

func checkRepo(Url string, repo string) (exists bool, format string, err error) {
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return false, format, err
	}
	req.Header.Set("Accept", applicationJSON)
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("User-Agent", DefaultUserAgent)

	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return false, format, err
	}

	var ok bool
	if resp.StatusCode == 200 {
		ok = true
	}

	if !ok {
		resp.Body.Close()
		return false, format, err
	}
	resBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var result types.Repos

	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.Println(err)
		os.Exit(100)
	}
	for _, v := range result {
		if v.Name == repo && v.Type == "hosted" {
			return true, v.Format, nil
		}
	}
	return false, format, nil
}
