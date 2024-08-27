package version

import (
	"fmt"
	"runtime"
	"time"
)

var (
	//x.y.z
	//x - major release
	//y - minor release
	//z - build number
	gitVersion = "v1.1.0-dirtree"
	gitCommit  = ""
	Time       = time.Now()
	buildDate  = Time.String()
)

type VersionInfo struct {
	GitVersion string `json:"gitVersion" yaml:"gitVersion"`
	GitCommit  string `json:"gitCommit" yaml:"gitCommit"`
	BuildDate  string `json:"buildDate" yaml:"buildDate"`
	GoVersion  string `json:"goVersion" yaml:"goVersion"`
	Compiler   string `json:"compiler" yaml:"compiler"`
	Program    string `json:"program" yaml:"program"`
}

func Get() *VersionInfo {
	return &VersionInfo{
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Program:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
