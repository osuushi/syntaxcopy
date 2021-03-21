package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os/exec"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/go-enry/go-enry/v2"
	"github.com/go-enry/go-enry/v2/data"
)

//go:embed copyHtml.applescript
var copyHtml string

// Rarer languages that cause misfires
var languageBlockList = map[string]bool{
	"viml":      true,
	"squidconf": true,
}

func detectLanguage(snippet string) string {
	// Get all languages supported by chroma
	allLanguages := []string{}
	chromaLanguagesMap := make(map[string]bool)
	for _, lexer := range lexers.Registry.Lexers {
		name := strings.ToLower(lexer.Config().Name)
		if languageBlockList[name] {
			continue
		}
		allLanguages = append(allLanguages, name)
		chromaLanguagesMap[name] = true
	}

	langName, _ := enry.GetLanguageByClassifier([]byte(snippet), allLanguages)
	// Convert back to alias from candidate because long names don't match between
	// chroma and enry
	for alias, name := range data.LanguageByAliasMap {
		if name == langName && chromaLanguagesMap[alias] {
			return alias
		}
	}

	// Fall back to return the language name
	return langName
}

func main() {
	input, err := exec.Command("pbpaste").Output()
	if err != nil {
		panic(err.Error())
	}

	println(string(input))

	// Unfortunately, Google Docs won't respect CSS tab-size, so replace all tabs
	// with two spaces
	spacified := strings.ReplaceAll(string(input), "\t", "  ")

	// Google Docs will eat the first space for some reason. Give it a little
	// morsel to satisfy its hunger.
	spacified = " " + spacified

	// Check if there's a "shebang" to force the language
	var lang string
	parts := strings.SplitN(spacified, "\n", 2)
	possibleShebangLine := strings.TrimSpace(parts[0])
	if withoutShebang := strings.TrimPrefix(possibleShebangLine, "#!"); withoutShebang != possibleShebangLine {
		lang = withoutShebang
		spacified = parts[1] // remove shebang
	} else {
		lang = detectLanguage(string(input))
	}
	println("Language:", lang)

	formatter := html.New(html.WithClasses(true))
	style := styles.Get("github")
	lexer := lexers.Get(lang)

	iterator, err := lexer.Tokenise(nil, spacified)
	if err != nil {
		panic(err.Error())
	}

	// Note the redirect to the file descriptor which allows us to get around
	// AppleScript's jank with stdin
	osascript := exec.Command("bash", "-c", `osascript -e "$0" 3<&0`, copyHtml)

	inPipe, err := osascript.StdinPipe()
	if err != nil {
		panic(err.Error())
	}

	stdin := bufio.NewWriter(inPipe)
	if err := osascript.Start(); err != nil {
		panic(err.Error())
	}

	stdin.WriteString("<style>")
	stdin.WriteString(`
		.chroma, .lang {
			font-family: Inconsolata;
			font-size: 11pt;
		}
		.lang {
			font-size: 9pt;
			color: #888;
		}
	`)
	formatter.WriteCSS(stdin, style)
	stdin.WriteString("</style>")
	stdin.WriteString(fmt.Sprintf(`<div class="lang">#!%s</div>`, lang))
	formatter.Format(stdin, style, iterator)
	stdin.Flush()
	inPipe.Close()

	if err := osascript.Wait(); err != nil {
		panic(err.Error())
	}
}
