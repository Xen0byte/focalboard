// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package notifymentions

import (
	"fmt"
	"sync"

	"github.com/mattermost/focalboard/server/services/notify"
	"github.com/wiggin77/merror"

	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

const (
	backendName = "notifyMentions"
)

type MentionListener interface {
	OnMention(userID string, evt notify.BlockChangeEvent)
}

// Backend provides the notification backend for @mentions.
type Backend struct {
	delivery Delivery
	logger   *mlog.Logger

	mux       sync.RWMutex
	listeners []MentionListener
}

func New(delivery Delivery, logger *mlog.Logger) *Backend {
	return &Backend{
		delivery: delivery,
		logger:   logger,
	}
}

func (b *Backend) Start() error {
	return nil
}

func (b *Backend) ShutDown() error {
	_ = b.logger.Flush()
	return nil
}

func (b *Backend) Name() string {
	return backendName
}

func (b *Backend) AddListener(l MentionListener) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.listeners = append(b.listeners, l)
}

func (b *Backend) RemoveListener(l MentionListener) {
	b.mux.Lock()
	defer b.mux.Unlock()
	list := make([]MentionListener, 0, len(b.listeners))
	for _, listener := range b.listeners {
		if listener != l {
			list = append(list, listener)
		}
	}
	b.listeners = list
}

func (b *Backend) BlockChanged(evt notify.BlockChangeEvent) error {
	if evt.Board == nil || evt.Card == nil {
		return nil
	}

	if evt.Action == notify.Delete {
		return nil
	}

	if evt.BlockChanged.Type != "text" && evt.BlockChanged.Type != "comment" {
		return nil
	}

	mentions := extractMentions(evt.BlockChanged)
	if len(mentions) == 0 {
		return nil
	}

	oldMentions := extractMentions(evt.BlockOld)
	merr := merror.New()

	b.mux.RLock()
	listeners := make([]MentionListener, 0, len(b.listeners))
	copy(listeners, b.listeners)
	b.mux.RUnlock()

	for username := range mentions {
		if _, exists := oldMentions[username]; exists {
			// the mention already existed; no need to notify again
			continue
		}

		extract := extractText(evt.BlockChanged.Title, username, newLimits())

		userID, err := b.delivery.Deliver(username, extract, evt)
		if err != nil {
			merr.Append(fmt.Errorf("cannot deliver notification for @%s: %w", username, err))
		}
		for _, listener := range listeners {
			safeCallListener(listener, userID, evt, b.logger)
		}
	}
	return merr.ErrorOrNil()
}

func safeCallListener(listener MentionListener, userID string, evt notify.BlockChangeEvent, logger *mlog.Logger) {
	// don't let panicky listeners stop notifications
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic calling @mention notification listener", mlog.Any("err", r))
		}
	}()
	listener.OnMention(userID, evt)
}