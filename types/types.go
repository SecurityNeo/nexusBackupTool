package types

type SrcConfig struct {
	Platform    PlatformConfig `json:"srcPlatform" required:"true"`
	NexusConfig NexusConfig    `json:"srcNexus" required:"true"`
	Backup      BackConfig     `json:"backup" required:"true"`
}

type TargetConfig struct {
	Platform    PlatformConfig `json:"srcPlatform" required:"true"`
	NexusConfig NexusConfig    `json:"srcNexus" required:"true"`
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
