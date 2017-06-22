# Build server
Lightweight build server that can be used to build your projects asynchronously while working with Vim.

Tired of switching panes in tmux or even moving to another window to build my project I decided to create this tool, now whenever I need to build my project, within Vim I type a shortcut and I see my project building in a different window (and different screen).

## Profiles
You use profiles to indicate what needs to be built, this way you can multiple instances of build server running at the same time.

Profile sample:

```javascript
{
    "host": "localhost",
    "port": ":8080",
    "directory": "/home/rendon/workspace/buildserver",
    "command": ["go", "build"]
}
```

## Run server
Run the server as follows:
```bash
$ buildserver /path/to/profile.json
```

## Send build request
```bash
$ buildclient /path/to/profile.json
```
