package sample

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/juju/errors"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {
	ctx.Logger().Debug("init settings")

	var settings Settings
	if err := metadata.MapToStruct(ctx.Settings(), &settings, true); err != nil {
		return nil, errors.Annotate(err, "MapToStruct [settings]")
	}

	var options interp.Options
	if settings.GoPath != "" {
		options.GoPath = settings.GoPath
	} else {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			return nil, errors.New("Can't get a valid GOPATH. Either define 'goPath' setting or ensure the 'GOPATH' env variable is set")
		}
		options.GoPath = gopath
	}

	ip := interp.New(options)
	ip.Use(stdlib.Symbols)
	ip.Use(interp.Symbols)

	act := Activity{
		Interpreter: ip,
	}

	return &act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	Interpreter *interp.Interpreter
}

// Metadata returns the activity's metadata
func (p *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval
func (p *Activity) Eval(ctx activity.Context) (bool, error) {
	var input Input
	if err := ctx.GetInputObject(&input); err != nil {
		return false, errors.Annotate(err, "GetInputObject")
	}

	if input.Source == "" {
		return false, errors.New("source is undefined")
	}

	if input.SourceType == "" {
		return false, errors.New("sourceType is undefined")
	}

	var sourceCode string
	switch input.SourceType {
	case "CODE":
		sourceCode = input.Source
	case "PATH":
		if !filepath.IsAbs(input.Source) {
			srcPath, err := filepath.Abs(input.Source)
			if err != nil {
				return false, errors.Annotate(err, "Abs [input source]")
			}

			src, err := ioutil.ReadFile(srcPath)
			if err != nil {
				return false, errors.Annotate(err, "ReadFile [source path]")
			}

			sourceCode = string(src)
		}

	}

	ctx.Logger().Debug("Generate interpreter context")
	_, err := p.Interpreter.Eval(sourceCode)
	if err != nil {
		return false, errors.Annotate(err, "Eval [sourceCode]")
	}

	ctx.Logger().Debug("Evaluate entrypoint")
	v, err := p.Interpreter.Eval("runner.Entrypoint")
	if err != nil {
		return false, errors.Annotate(err, "Eval [entrypoint]")
	}

	ep := v.Interface()
	if ep == nil {
		return false, errors.New("entrypoint not available")
	}

	entryPoint, ok := ep.(func(interface{}) (interface{}, error))
	if !ok {
		return false, errors.New("invalid entrypoint signature")
	}

	var output Output
	result, err := entryPoint(input.Arguments)
	output.Result = result
	if err != nil {
		output.Error = err.Error()
	}

	if err := ctx.SetOutputObject(&output); err != nil {
		return false, errors.Annotate(err, "SetOutputObject")
	}

	return true, nil
}
