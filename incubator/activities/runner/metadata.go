package sample

import (
	"github.com/juju/errors"
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	GoPath string `md:"goPath"`
}

type Input struct {
	Source     string      `md:"source,required"`
	SourceType string      `md:"sourceType,required"`
	Arguments  interface{} `md:"args"`
}

func (p *Input) FromMap(values map[string]interface{}) error {
	source, err := coerce.ToString(values["source"])
	if err != nil {
		return errors.Annotate(err, "ToString [source]")
	}

	p.Source = source

	sourceType, err := coerce.ToString(values["sourceType"])
	if err != nil {
		return errors.Annotate(err, "ToString [sourceType]")
	}

	p.SourceType = sourceType

	if args, ok := values["args"]; ok {
		p.Arguments = args
	}

	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"source":     r.Source,
		"sourceType": r.SourceType,
	}
}

type Output struct {
	Result interface{} `md:"result"`
	Error  string      `md:"error"`
}

func (p *Output) FromMap(values map[string]interface{}) error {
	if result, ok := values["result"]; ok {
		p.Result = result
	} else {
		return errors.New("FromMap: 'result' value not available")
	}

	if err, ok := values["error"]; ok {
		switch v := err.(type) {
		case string:
			p.Error = v
		case error:
			p.Error = v.Error()
		default:
			return errors.New("FromMap: can't assigning 'error' value")
		}
	} else {
		return errors.New("FromMap: 'error' value not available")
	}

	return nil
}

func (p *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": p.Result,
		"error":  p.Error,
	}
}
