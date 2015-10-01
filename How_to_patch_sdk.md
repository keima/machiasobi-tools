How to patch sdk
-----

1. Download `file_watcher.patch` from here: [Issue 11345 - googleappengine - PATCH: Modify the file watchers to ignore any file specified in skip-files - Google App Engine - Google Project Hosting](https://code.google.com/p/googleappengine/issues/detail?id=11345#c8)
 - Or, use `file_watcher.patch` from this repo.(Optimize `mtime_file_watcher#set_skip_files_re()` for GAE/Go)
2. `cd cd ~/go_appengine/google/appengine/tools/devappserver2/; patch -u < ~/Download/file_watcher.patch`
3. Enjoy!
