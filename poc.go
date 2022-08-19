package pocRegistry

import (
	"fmt"
	"github.com/12end/pocRegistry/cli"
	"github.com/12end/requests"
	"net/url"
	"strings"
)

type POC struct {
	Name        string
	ProductName string
	Alias       string
	Desc        string
	Help        string
	Effective   string
	Check       cli.Command
	Exploit     []cli.Command
}

// poc log
func (r Registry) check(poc POC, target *url.URL, trace bool, args []string) (vulnerable bool, traceInfo []requests.TraceInfo) {
	r.Logger.Info(fmt.Sprintf("Check %s/%s(%s) for %s", poc.ProductName, poc.Name, poc.Alias, target.String()))
	params := cli.NewFlagSet("check")
	for _, flag := range poc.Check.Flags {
		if flag.IsRequired() && args == nil {
			// not cli mode
			return
		}
		flag.Apply(params.Set)
	}
	//添加-h进行help打印
	err := params.Set.Parse(args)
	if err != nil {
		fmt.Println(err)
	}
	vulnerable, traceInfo = poc.Check.Action(&cli.Context{Trace: trace, Target: target, Logger: r.Logger, Params: params})
	if vulnerable {
		r.Logger.Warn(fmt.Sprintf("%s has vulnerability for %s/%s!", target.String(), poc.ProductName, poc.Name))
	}
	return
}

func (r Registry) ExecutePOC(target *url.URL, productName string, pocName string, args []string) (vulnerable bool) {
	productName_ := strings.ToLower(productName)
	pocName_ := strings.ToLower(pocName)
	if _, ok := r.Pocs[productName_]; ok {
		if poc, ok := r.Pocs[productName_][pocName_]; ok {
			vulnerable, _ = r.check(poc, target, false, args)
			return
		} else {
			r.Logger.Error(fmt.Sprintf("No such poc(%s) in product: %s", pocName, productName))
			return false
		}
	} else {
		r.Logger.Error(fmt.Sprintf("No such product: %s", productName))
		return false
	}
}

func (r Registry) ExecutePOCWithTrace(target *url.URL, productName string, pocName string, args []string) (vulnerable bool, trace []requests.TraceInfo) {
	productName_ := strings.ToLower(productName)
	pocName_ := strings.ToLower(pocName)
	if _, ok := r.Pocs[productName_]; ok {
		if poc, ok := r.Pocs[productName_][pocName_]; ok {
			vulnerable, trace = r.check(poc, target, true, args)
			return
		} else {
			r.Logger.Error(fmt.Sprintf("No such poc(%s) in product: %s", pocName, productName))
			return false, nil
		}
	} else {
		r.Logger.Error(fmt.Sprintf("No such product: %s", productName))
		return false, nil
	}
}

func (r Registry) ExecutePOCs(target *url.URL, productName string, args []string) (result map[string]string) {
	productName_ := strings.ToLower(productName)
	result = map[string]string{}
	if _, ok := r.Pocs[productName_]; ok {
		for pocName, poc := range r.Pocs[productName_] {
			if vulnerable, _ := r.check(poc, target, false, args); vulnerable {
				result[pocName] = poc.Alias
			}
		}
		return
	} else {
		r.Logger.Error(fmt.Sprintf("No such product: %s", productName))
		return
	}
}

func (r Registry) ExecutePOCsWithTrace(target *url.URL, productName string, args []string) (result map[string][]requests.TraceInfo, aliases map[string]string) {
	result = map[string][]requests.TraceInfo{}
	aliases = map[string]string{}
	productName_ := strings.ToLower(productName)
	if _, ok := r.Pocs[productName_]; ok {
		for pocName, poc := range r.Pocs[productName_] {
			vulnerable, trace := r.check(poc, target, true, args)
			if vulnerable {
				result[pocName] = trace
				aliases[pocName] = poc.Alias
			}
		}
		return
	} else {
		r.Logger.Error(fmt.Sprintf("No such product: %s", productName))
		return
	}
}
