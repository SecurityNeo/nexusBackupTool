package backup

import "sync"

type Components struct {
	Items				[]Item `json:"items"`
	ContinuationToken	string `json:"continuationToken"`
}

type Item struct {
	Assets				[]Asset `json:"assets"`
}

type Asset struct {
	DownloadUrl		string	`json:"downloadUrl"`
	Path 			string  `json:"path"`
	Check			Checksum  `json:"checksum"`
}

type Checksum struct {
	MD5				string `json:"md5"`
}

type FileData		[]File

type DownloadInfo struct {
	Files 			FileData
	Wg   			sync.WaitGroup
	Chs  			chan int
	Ans  			chan bool
}

type File struct {
	DownloadUrl		string			`json:"downloadUrl"`
	Path 			string			`json:"path"`
	MD5				string			`json:"md5"`
}
