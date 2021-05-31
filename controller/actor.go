package controller

import (
	"time"

	"n_communication/gateway/email"
	"n_communication/gateway/sms"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

const maxIdealTime = time.Minute * -1
const maxAliveTime = time.Minute * -5
const mailboxCapacity = 10

type ActorType string

type Action struct {
	Channel string
	To      string
	From    string
	Payload string
	Title   string
}

type Actor interface {
	IsFull() bool
	IsAlive() bool
	IsIdeal() bool
	Add(action Action) bool
	Name() string
	Shutdown()
}

type actor struct {
	lastHeartbeat   time.Time
	lastActionTime  time.Time
	mailbox         chan Action
	mailboxSize     int
	mailboxCapacity int
	name            string
	emailGateway    email.EmailGateway
	smsGateway      sms.SMSGateway
}

func NewActor() Actor {
	a := &actor{
		lastHeartbeat:   time.Now(),
		lastActionTime:  time.Now(),
		mailbox:         make(chan Action, mailboxCapacity),
		mailboxSize:     0,
		mailboxCapacity: mailboxCapacity,
		name:            uuid.NewString(),
		emailGateway:    email.New(),
		smsGateway:      sms.New(),
	}
	go a.Init()
	return a
}

func (a *actor) IsFull() bool {
	if a.mailboxCapacity == a.mailboxSize {
		return true
	}
	return false
}

func (a *actor) IsAlive() bool {
	return !a.lastHeartbeat.Before(time.Now().Add(maxAliveTime))
}

func (a *actor) IsIdeal() bool {
	return a.lastActionTime.Before(time.Now().Add(maxIdealTime))
}

func (a *actor) Name() string {
	return a.name
}

func (a *actor) Shutdown() {
	close(a.mailbox)
}

func (a *actor) Add(action Action) bool {
	a.mailbox <- action
	a.mailboxSize++
	return true
}

func (a *actor) Init() {
	for {
		select {
		case action := <-a.mailbox:
			{
				// handle email
				if action.Channel == actionTypeEmail {
					status, err := a.emailGateway.Send(action.To, action.From, action.Payload, action.Title)
					if err != nil {
						zap.L().Error("got error", zap.Error(err))
					} else {
						zap.L().Info("email sent",
							zap.Bool("status", status),
							zap.String("name", a.name))
					}
				}

				// handle sms
				if action.Channel == actionTypeSMS {
					status, err := a.smsGateway.Send(action.To, action.From, action.Payload)
					if err != nil {
						zap.L().Error("got error", zap.Error(err))
					} else {
						zap.L().Info("sms sent",
							zap.Bool("status", status),
							zap.String("name", a.name))
					}
				}

				if action.Channel == actionTypePush {
					zap.L().Info("sent push notification",
						zap.String("name", a.name))
				}

				a.mailboxSize--
				a.lastActionTime = time.Now()
			}
		case t := <-time.Tick(1 * time.Minute):
			{
				a.lastHeartbeat = t
			}
		}
	}
}
