package types

import (
	"github.com/SecurityNeo/NexusbackupTool/common"
	"sync"
)

type SrcConfig struct {
	Platform    PlatformConfig `json:"srcPlatform" required:"true"`
	NexusConfig NexusConfig    `json:"srcNexus" required:"true"`
	Backup      BackConfig     `json:"backup" required:"true"`
}

type TargetConfig struct {
	Platform    PlatformConfig `json:"targetPlatform" required:"true"`
	NexusConfig NexusConfig    `json:"targetNexus" required:"true"`
	Backup      BackConfig     `json:"backup" required:"true"`
}

type PlatformConfig struct {
	Ip   string `json:"ip" required:"true"`
	Port string `json:"port" required:"true"`
}

type NexusConfig struct {
	Repository string `json:"repository" required:"true"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type BackConfig struct {
	Directory string `json:"directory" required:"true"`
}

type Components struct {
	Items             []Item `json:"items"`
	ContinuationToken string `json:"continuationToken"`
}

type Item struct {
	Assets []Asset `json:"assets"`
}

type Asset struct {
	DownloadUrl string   `json:"downloadUrl"`
	Path        string   `json:"path"`
	Check       Checksum `json:"checksum"`
}

type Checksum struct {
	MD5 string `json:"md5"`
}

type FileData []File

type DownloadInfo struct {
	Files FileData
	Wg    sync.WaitGroup
	Chs   chan int
	Ans   chan bool
}

type UploadInfo struct {
	Files FileData
	Wg    sync.WaitGroup
	Chs   chan int
	Ans   chan bool
}

func (di *DownloadInfo) Work(downloadUrl string, targetPath string, targetMd5 string) {
	defer func() {
		<-di.Chs
		di.Wg.Done()
	}()
	common.Download(downloadUrl, targetPath, targetMd5)
	di.Ans <- true
}

func (ui *UploadInfo) Work(uploadUrl, filePath, abFilePath, username, password, format string) {
	defer func() {
		<-ui.Chs
		ui.Wg.Done()
	}()
	common.Upload(uploadUrl, filePath, abFilePath, username, password, format)
	ui.Ans <- true
}

type File struct {
	DownloadUrl string `json:"downloadUrl"`
	Path        string `json:"path"`
	MD5         string `json:"md5"`
}

type Repos []Repo

type Repo struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format"`
}
