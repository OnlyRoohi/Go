/*
 * TgMusicBot - helpers for VC
 */

package vc

import (
	"ashokshau/tgmusic/src/core/cache"
	"ashokshau/tgmusic/src/core/dl"
	"ashokshau/tgmusic/src/utils"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os/exec"
	"strconv"
	"strings"
	"time"

	dt "github.com/AshokShau/gotdbot"
	"github.com/amarnathcjd/gogram/telegram"
)

// downloadAndPrepareSong handles the download and preparation of a song for playback.
// It returns an error if the download or preparation fails.
func (c *TelegramCalls) downloadAndPrepareSong(bot *dt.Client, song *utils.CachedTrack, reply *dt.Message) error {
	if song.FilePath != "" {
		return nil
	}

	dlPath, err := dl.DownloadCachedTrack(song, bot)
	song.FilePath = dlPath
	if err != nil || song.FilePath == "" {
		_, _ = reply.EditText(bot, "⚠️ Download failed. Skipping track...", nil)
		if song.FilePath == "" {
			return fmt.Errorf("empty download path")
		}
		return err
	}

	// If file doesn't exist after download, treat as failure
	// (double-checking path validity)
	if song.FilePath == "" {
		_, _ = reply.EditText(bot, "⚠️ Download produced empty path. Skipping track...", nil)
		return fmt.Errorf("empty download path")
	}

	return nil
}
