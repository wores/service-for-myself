package input

import (
	"context"
	"encoding/json"

	"github.com/nlopes/slack/slackevents"
	"google.golang.org/appengine/log"
)

// InputSlackEvent SlackEvent
type InputSlackEvent struct {
	eventsAPIEvent slackevents.EventsAPIEvent
}

// GetEventsAPIEvent プロパティのeventsAPIEventを返す
func (inputSlackEvent *InputSlackEvent) GetEventsAPIEvent() slackevents.EventsAPIEvent {
	return inputSlackEvent.eventsAPIEvent
}

// ParseFromStringToInputSlackEvent 文字列からパースする
func ParseFromStringToInputSlackEvent(
	ctx context.Context,
	bodyString string,
	verificationToken string,
) (*InputSlackEvent, error) {
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(bodyString),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{verificationToken}),
	)
	if err != nil {
		log.Errorf(ctx, "errorの中身 = %#v", err)
		return nil, ErrParseSlackEvent
	}
	log.Debugf(ctx, "eventsAPIEvent = %#v", eventsAPIEvent)
	log.Debugf(ctx, "eventsAPIEvent.Type = %s", eventsAPIEvent.Type)

	return &InputSlackEvent{eventsAPIEvent}, nil
}
