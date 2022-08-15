package pocRegistry

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Registry struct {
	// map[productName][vulName]
	pocs     map[string]map[string]POC
	logLevel *zap.AtomicLevel
	Logger   *zap.Logger
}

func NewRegistry() Registry {
	level := zap.NewAtomicLevelAt(zap.WarnLevel)
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

func (r Registry) Register(p POC) {
	product := p.ProductName
	if _, ok := r.pocs[product]; !ok {
		r.pocs[product] = map[string]POC{
			p.Name: p,
		}
	} else {
		r.pocs[product][p.Name] = p
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
