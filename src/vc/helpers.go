diff --git a/src/vc/helpers.go b/src/vc/helpers.go
index 965dea3..0000000
--- a/src/vc/helpers.go
+++ b/src/vc/helpers.go
@@
 func (c *TelegramCalls) downloadAndPrepareSong(bot *td.Client, song *utils.CachedTrack, reply *td.Message) error {
     if song.FilePath != "" {
         return nil
     }
 
     dlPath, err := dl.DownloadCachedTrack(song, bot)
     song.FilePath = dlPath
     if err != nil || song.FilePath == "" {
         _, _ = reply.EditText(bot, "⚠️ Download failed. Skipping track...", nil)
         return err
     }
 
+    // If file doesn't exist after download, treat as failure
+    if song.FilePath == "" {
+        _, _ = reply.EditText(bot, "⚠️ Download produced empty path. Skipping track...", nil)
+        return fmt.Errorf("empty download path")
+    }
+
     return nil
 }
