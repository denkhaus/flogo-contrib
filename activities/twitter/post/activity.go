package sample

import (
	"strings"

	"github.com/denkhaus/flogo/activities/twitter"
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

	act := Activity{
		ConsumerKey:       strings.TrimSpace(settings.ConsumerKey),
		ConsumerSecret:    strings.TrimSpace(settings.ConsumerSecret),
		AccessToken:       strings.TrimSpace(settings.AccessToken),
		AccessTokenSecret: strings.TrimSpace(settings.AccessTokenSecret),
	}

	if act.ConsumerKey == "" || act.ConsumerSecret == "" ||
		act.AccessToken == "" || act.AccessTokenSecret == "" {
		return nil, errors.New("settings incomplete")
	}

	return &act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// Metadata returns the activity's metadata
func (p *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval
func (p *Activity) Eval(ctx activity.Context) (done bool, err error) {
	var input Input
	if err = ctx.GetInputObject(&input); err != nil {
		return true, errors.Annotate(err, "GetInputObject")
	}

	if input.Content == "" {
		return true, errors.New("content is empty")
	}

	ctx.Logger().Debug("post tweet")
	statusCode, message := twitter.PostTweet(
		p.ConsumerKey,
		p.ConsumerSecret,
		p.AccessToken,
		p.AccessTokenSecret,
		input.Content,
	)

	output := Output{
		StatusCode: statusCode,
		Message:    message,
	}

	if err := ctx.SetOutputObject(&output); err != nil {
		return true, errors.Annotate(err, "SetOutputObject")
	}

	return true, nil
}
