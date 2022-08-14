package pocRegistry

import (
	"errors"
	"fmt"
	"github.com/12end/requests"
	"net/url"
	"strings"
)

type Registry struct {
	// map[productName][vulName]
	pocs map[string]map[string]POC
}

func NewRegistry() Registry {
	return Registry{pocs: map[string]map[string]POC{}}
}

func (r Registry) Register(poc POC) {
	product := poc.ProductName
	if _, ok := r.pocs[product]; !ok {
		r.pocs[product] = map[string]POC{
			poc.Name: poc,
		}
	} else {
		r.pocs[product][poc.Name] = poc
	}
}

func (r Registry) Unset(pocName string) {
	// productName = arr[0]
	// pocName = arr[1]
	arr := strings.SplitN(pocName, "/", 2)
	if len(arr) == 1 {
		productName := arr[0]
		delete(r.pocs, productName)
	} else {
		productName := arr[0]
		pocName := arr[1]
		if _, ok := r.pocs[productName]; ok {
			delete(r.pocs[productName], pocName)
		}
	}
}

func (r Registry) ExecutePOC(target *url.URL, productName string, pocName string) (vulnerable bool, err error) {
	if _, ok := r.pocs[productName]; ok {
		if poc, ok := r.pocs[productName][pocName]; ok {
			vulnerable, _ = poc.Check(target, false)
			return
		} else {
			err = errors.New(fmt.Sprintf("No such poc(%s) in product: %s", pocName, productName))
			return false, err
		}
	} else {
		err = errors.New(fmt.Sprintf("No such product: %s", productName))
		return false, err
	}
}

func (r Registry) ExecutePOCWithTrace(target *url.URL, productName string, pocName string) (vulnerable bool, trace []requests.TraceInfo, err error) {
	if _, ok := r.pocs[productName]; ok {
		if poc, ok := r.pocs[productName][pocName]; ok {
			vulnerable, trace = poc.Check(target, true)
			return
		} else {
			err = errors.New(fmt.Sprintf("No such poc(%s) in product: %s", pocName, productName))
			return false, nil, err
		}
	} else {
		err = errors.New(fmt.Sprintf("No such product: %s", productName))
		return false, nil, err
	}
}

func (r Registry) ExecutePOCs(target *url.URL, productName string) (result map[string]bool, err error) {
	if _, ok := r.pocs[productName]; ok {
		for pocName, poc := range r.pocs[productName] {
			if vulnerable, _ := poc.Check(target, false); vulnerable {
				result[pocName] = true
			}
		}
		return
	} else {
		err = errors.New(fmt.Sprintf("No such product: %s", productName))
		return nil, err
	}
}

func (r Registry) ExecutePOCsWithTrace(target *url.URL, productName string) (result map[string][]requests.TraceInfo, err error) {
	if _, ok := r.pocs[productName]; ok {
		for pocName, poc := range r.pocs[productName] {
			vulnerable, trace := poc.Check(target, true)
			if vulnerable {
				result[pocName] = trace
			}
		}
		return
	} else {
		err = errors.New(fmt.Sprintf("No such product: %s", productName))
		return nil, err
	}
}
