package pocRegistry

import (
	"github.com/12end/requests"
	"net/url"
)

type POC struct {
	Name        string
	ProductName string
	Alias       string
	Desc        string
	Help        string
	Effective   string
	//Check
	Check func(target *url.URL, trace bool) (vulnerable bool, traceInfo []requests.TraceInfo)
	//Exploit
	Exploit func(target *url.URL) (result Result)
}

type Result struct {
	Success bool
	Output  string
}
