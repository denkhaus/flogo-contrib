package sample

import (
	"github.com/juju/errors"
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	ConsumerKey       string `md:"consumerKey,required"`
	ConsumerSecret    string `md:"consumerSecret,required"`
	AccessToken       string `md:"accessToken,required"`
	AccessTokenSecret string `md:"accessTokenSecret,required"`
}

type Input struct {
	Content string `md:"content,required"`
}

func (p *Input) FromMap(values map[string]interface{}) error {
	content, err := coerce.ToString(values["content"])
	if err != nil {
		return errors.Annotate(err, "ToString [content]")
	}

	p.Content = content
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"content": r.Content,
	}
}

type Output struct {
	StatusCode int    `md:"statusCode"`
	Message    string `md:"message"`
}

func (p *Output) FromMap(values map[string]interface{}) error {
	statusCode, err := coerce.ToInt(values["statusCode"])
	if err != nil {
		return errors.Annotate(err, "ToInt [statusCode]")
	}

	p.StatusCode = statusCode

	msg, err := coerce.ToString(values["message"])
	if err != nil {
		return errors.Annotate(err, "ToString [message]")
	}

	p.Message = msg

	return nil
}

func (p *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": p.StatusCode,
		"message":    p.Message,
	}
}
