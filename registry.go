package pocRegistry

import (
	"fmt"
	"github.com/12end/requests"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/url"
	"strings"
)

type Registry struct {
	// map[productName][vulName]
	pocs     map[string]map[string]POC
	logLevel *zap.AtomicLevel
	Logger   *zap.Logger
}

func NewRegistry() Registry {
	level := zap.NewAtomicLevelAt(zap.ErrorLevel)
	config := zap.Config{
		Level:            level,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			TimeKey:     "time",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			EncodeTime:  zapcore.TimeEncoderOfLayout("15:04:05"),
		},
	}
	log, err := config.Build()
	if err != nil {
		fmt.Println(err)
	}
	return Registry{pocs: map[string]map[string]POC{}, logLevel: &level, Logger: log}
}

func (r Registry) SetLogLevel(level zapcore.Level) {
	r.logLevel.SetLevel(level)
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

// poc log
func (r Registry) check(poc POC, target *url.URL, trace bool) (vulnerable bool, traceInfo []requests.TraceInfo) {
	r.Logger.Info(fmt.Sprintf("Checking %s/%s(%s) for %s", poc.ProductName, poc.Name, poc.Alias, target.String()))
	vulnerable, traceInfo = poc.Check(target, trace)
	if vulnerable {
		r.Logger.Info(fmt.Sprintf("%s has vulnerability for %s/%s!", target.String(), poc.ProductName, poc.Name))
	}
	return
}

func (r Registry) ExecutePOC(target *url.URL, productName string, pocName string) (vulnerable bool) {
	if _, ok := r.pocs[productName]; ok {
		if poc, ok := r.pocs[productName][pocName]; ok {
			vulnerable, _ = r.check(poc, target, false)
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

func (r Registry) ExecutePOCWithTrace(target *url.URL, productName string, pocName string) (vulnerable bool, trace []requests.TraceInfo) {
	if _, ok := r.pocs[productName]; ok {
		if poc, ok := r.pocs[productName][pocName]; ok {
			vulnerable, trace = r.check(poc, target, true)
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

func (r Registry) ExecutePOCs(target *url.URL, productName string) (result map[string]bool) {
	if _, ok := r.pocs[productName]; ok {
		for pocName, poc := range r.pocs[productName] {
			if vulnerable, _ := r.check(poc, target, false); vulnerable {
				result[pocName] = true
			}
		}
		return
	} else {
		r.Logger.Error(fmt.Sprintf("No such product: %s", productName))
		return
	}
}

func (r Registry) ExecutePOCsWithTrace(target *url.URL, productName string) (result map[string][]requests.TraceInfo) {
	if _, ok := r.pocs[productName]; ok {
		for pocName, poc := range r.pocs[productName] {
			vulnerable, trace := r.check(poc, target, true)
			if vulnerable {
				result[pocName] = trace
			}
		}
		return
	} else {
		r.Logger.Error(fmt.Sprintf("No such product: %s", productName))
		return
	}
}
