package controller

import (
	"errors"
	"time"

	"go.uber.org/zap"
)

const actionTypeEmail = "email"
const actionTypeSMS = "sms"
const actionTypePush = "push"

const supervisionDelay = time.Second * 5

var actorTypes = []ActorType{
	actionTypeEmail,
	actionTypeSMS,
	actionTypePush,
}

type Coordinator interface {
	Send(channel string, to string, from string, payload string, title string) (bool, error)
}

type coordinator struct {
	directory map[ActorType][]Actor
	reqQueue  chan Action
}

func New() Coordinator {
	d := map[ActorType][]Actor{}
	for _, v := range actorTypes {
		// TODO: increase and decrease number of actors based on the load
		d[v] = []Actor{NewActor()}
	}

	rq := make(chan Action, 1000)

	go superviseActors(d, rq)

	return &coordinator{directory: d, reqQueue: rq}
}

func (c *coordinator) Send(channel string, to string, from string, payload string, title string) (bool, error) {
	if len(c.reqQueue) < cap(c.reqQueue) {
		c.reqQueue <- Action{Channel: channel, To: to, From: from, Payload: payload, Title: title}
		return true, nil
	}

	return false, errors.New("node is loaded please retry after some time")
}

func superviseActors(directory map[ActorType][]Actor, reqQueue chan Action) {
	zap.L().Info("started running supervisor")
	var action Action
	for {
		action = <-reqQueue
		actors := directory[ActorType(action.Channel)]
		for _, actor := range actors {
			if actor.IsAlive() && !actor.IsFull() {
				actor.Add(action)
			}
		}
	}
}
