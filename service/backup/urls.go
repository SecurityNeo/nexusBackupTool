package backup

import "bytes"

const (
	URI  			= "/service/rest/v1/components?"
	ParamRepoName 	= "repository"
)

func getURL(BaseUrl string, repo string) string {
	var buffer bytes.Buffer
	buffer.WriteString(BaseUrl)
	buffer.WriteString(URI)
	buffer.WriteString(ParamRepoName)
	buffer.WriteString("=")
	buffer.WriteString(repo)
	url := buffer.String()
	return url
}
