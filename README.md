# rename

CLI app that renames audio files within a directory based on meta data that is ascertained with [dhowden/tag](https://github.com/dhowden/tag "dhowden/tag").

`go run main.go` runs the application in the working directory while `go run main.go ~/some/path` runs the application in the provided path. It only works on audio files and subdirectories are not effected. You will me asked to confirm the action, before anything happens. Be safe.