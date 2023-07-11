package recover

import "bytes"

const (
	URI  			= "/service/rest/v1/components?"
	ParamRepoName 	= "repository"
	RepoURI			= "/service/rest/v1/repositories"
)

func getRepoURL(BaseUrl string) string {
	var buffer bytes.Buffer
	buffer.WriteString(BaseUrl)
	buffer.WriteString(RepoURI)
	url := buffer.String()
	return url
}

func uploadURL(BaseUrl string, repo string) string {
	var buffer bytes.Buffer
	buffer.WriteString(BaseUrl)
	buffer.WriteString(URI)
	buffer.WriteString(ParamRepoName)
	buffer.WriteString("=")
	buffer.WriteString(repo)
	url := buffer.String()
	return url
}
