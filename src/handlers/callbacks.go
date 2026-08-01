/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"ashokshau/tgmusic/src/core/cache"
	"ashokshau/tgmusic/src/core/db"
	"ashokshau/tgmusic/src/vc"
	"strings"

	td "github.com/AshokShau/gotdbot"
)

// callbackHandler handles inline button callbacks related to playback and queue.
func callbackHandler(c *td.Client, cb *td.CallbackQuery) error {
	data := cb.Data
	chatID := cb.ChatId

	switch {
	case strings.Contains(data, "play_add_to_list"):
		// Minimal safe handling: ensure DB is ready and acknowledge.
		if db.Instance == nil {
			_ = cb.Answer(c, 0, false, "Database not ready.", "")
			return nil
		}
		playlists, err := db.Instance.GetUserPlaylists(cb.SenderUserId)
		if err != nil {
			_ = cb.Answer(c, 0, false, "Unable to fetch playlists.", "")
			return nil
		}
		if len(playlists) == 0 {
			_ = cb.Answer(c, 0, false, "You have no playlists.", "")
			return nil
		}
		// For now, just acknowledge and instruct user to use playlist commands in chat.
		_ = cb.Answer(c, 0, false, "Open your playlists with /playlists and add the track from there.", "")
		return nil

	case strings.HasPrefix(data, "play_now_"):
		trackID := strings.TrimPrefix(data, "play_now_")
		if ok := cache.ChatCache.MoveTrackToFront(chatID, trackID); !ok {
			_ = cb.Answer(c, 0, false, "Track not found in queue.", "")
			return nil
		}
		if err := vc.Calls.PlayNext(c, chatID); err != nil {
			_ = cb.Answer(c, 0, false, "Unable to play the track.", "")
			return nil
		}
		_ = cb.Answer(c, 0, false, "Playing now.", "")
		_ = c.DeleteMessages(chatID, []int64{cb.MessageId}, &td.DeleteMessagesOpts{Revoke: true})
		return nil
	}

	return nil
}
