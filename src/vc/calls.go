diff --git a/src/vc/calls.go b/src/vc/calls.go
index 44a5b04..0000000
--- a/src/vc/calls.go
+++ b/src/vc/calls.go
@@
-    if db.Instance.GetLoggerStatus() {
-        go sendLogger(bot, chatID, cache.ChatCache.GetPlayingTrack(chatID))
-    }
+    // Guard against nil DB (during startup/shutdown) and only call logger when available
+    if db.Instance != nil {
+        if db.Instance.GetLoggerStatus() {
+            go sendLogger(bot, chatID, cache.ChatCache.GetPlayingTrack(chatID))
+        }
+    }
 
     return nil
 }
