package recover

type Repos []Repo

type Repo struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format"`
}
