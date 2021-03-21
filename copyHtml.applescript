use framework "AppKit"
use framework "Foundation"
use scripting additions

-- AppleScript closes STDIN immediately, but not a file descriptor we give it.
-- We'll wrap the call in bash -c to redirect stdin
set input to do shell script "cat 0<&3"

set appl to (current application)

set pboard to appl's NSPasteboard's generalPasteboard()
set htmlType to appl's NSPasteboardTypeHTML
set strType to appl's NSPasteboardTypeString

set htmlString to appl's NSString's stringWithString:input
set htmlData to htmlString's dataUsingEncoding:(appl's NSUTF8StringEncoding)

-- We still want to be able to paste as plain text
set originalData to pboard's dataForType:strType

pboard's clearContents()
pboard's setData:htmlData forType:htmlType
pboard's setData:originalData forType:strType
