package sample

import (
	"github.com/juju/errors"
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	ApiKey string `md:"apiKey,required"`
}

type Input struct {
	Content string `md:"content,required"`
	ChatID  int64  `md:"chatId,required"`
}

func (p *Input) FromMap(values map[string]interface{}) error {
	content, err := coerce.ToString(values["content"])
	if err != nil {
		return errors.Annotate(err, "ToString [content]")
	}

	p.Content = content

	chatID, err := coerce.ToInt64(values["chatId"])
	if err != nil {
		return errors.Annotate(err, "ToInt64 [chatId]")
	}

	p.ChatID = chatID

	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"content": r.Content,
		"chatId":  r.ChatID,
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
