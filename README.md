# syntaxcopy

This is a very hacky solution to syntax highlighting in Google Docs.

## Usage

First, copy a block of code onto your clipboard. By default, syntaxcopy will try to infer the language using [go-enry](https://github.com/go-enry/go-enry), but it isn't very accurate. You can force a specific language by adding a line like `#!go` before your block.

Then run `syntaxcopy` in your terminal. This will replace your clipboard with the syntax highlighted code in HTML format. You can now paste this into a doc.

A handy Hammerspoon config is:

```lua
hs.hotkey.bind('cmd alt', 'c', nil, function ()
  hs.eventtap.event.newKeyEvent({"cmd"}, "c", true):post()
  hs.task.new("/path/to/syntaxcopy", nil):start()
end, nil, nil)
```

This will bind cmd+alt+c to automatically apply syntaxcopy.
