diff --git a/src/handlers/callbacks.go b/src/handlers/callbacks.go
index 18fd7b7..0000000
--- a/src/handlers/callbacks.go
+++ b/src/handlers/callbacks.go
@@
-    case strings.Contains(data, "play_add_to_list"):
-        playlists, err := db.Instance.GetUserPlaylists(cb.SenderUserId)
-        if err != nil {
-            _ = cb.Answer(c, 0, false, "Unable to fetch playlists.", "")
-            return nil
-        }
+    case strings.Contains(data, "play_add_to_list"):
+        if db.Instance == nil {
+            _ = cb.Answer(c, 0, false, "Database not ready.", "")
+            return nil
+        }
+        playlists, err := db.Instance.GetUserPlaylists(cb.SenderUserId)
+        if err != nil {
+            _ = cb.Answer(c, 0, false, "Unable to fetch playlists.", "")
+            return nil
+        }
@@
-        err = db.Instance.AddSongToPlaylist(playlistID, song)
-        if err != nil {
-            _ = cb.Answer(c, 0, false, "Unable to add track to playlist.", "")
-            return nil
-        }
+        if db.Instance == nil {
+            _ = cb.Answer(c, 0, false, "Database not ready.", "")
+            return nil
+        }
+        err = db.Instance.AddSongToPlaylist(playlistID, song)
+        if err != nil {
+            _ = cb.Answer(c, 0, false, "Unable to add track to playlist.", "")
+            return nil
+        }
@@
-        playlist, err := db.Instance.GetPlaylist(playlistID)
+        playlist, err := db.Instance.GetPlaylist(playlistID)
         if err != nil {
             _ = cb.Answer(c, 0, false, "Playlist not found.", "")
             return nil
         }
@@
-    case strings.HasPrefix(data, "play_now_"):
+    case strings.HasPrefix(data, "play_now_"):
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
