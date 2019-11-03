package sample

import (
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/juju/errors"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
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

	act := Activity{
		ApiKey: strings.TrimSpace(settings.ApiKey),
	}

	if act.ApiKey == "" {
		return nil, errors.New("ApiKey is empty")
	}

	return &act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	ApiKey string
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

	if input.Content == "" {
		return false, errors.New("content is empty")
	}

	if input.ChatID <= 0 {
		return false, errors.New("invalid ChatID")
	}

	ctx.Logger().Debug("Connect Telegram API")
	bot, err := telegram.NewBotAPI(p.ApiKey)
	if err != nil {
		return false, errors.Annotate(err, "NewBotAPI")
	}

	ctx.Logger().Debugf("Authorized username -> %s", bot.Self.UserName)

	msgConfig := telegram.NewMessage(input.ChatID, input.Content)
	msg, err := bot.Send(msgConfig)
	if err != nil {
		return false, errors.Annotate(err, "Send")
	}

	msgIDString, err := coerce.ToString(msg.MessageID)
	if err != nil {
		return true, errors.Annotate(err, "ToString [messageID]")
	}

	output := Output{
		StatusCode: 200,
		Message:    msgIDString,
	}

	if err := ctx.SetOutputObject(&output); err != nil {
		return true, errors.Annotate(err, "SetOutputObject")
	}

	return true, nil
}
