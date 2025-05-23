+++
Categories = ["Widgets"]
+++

A **text editor** is a [[widget]] for editing complex text. It supports advanced code editing features, like syntax highlighting, completion, undo and redo, copy and paste, rectangular selection, and word, line, and page based navigation, selection, and deletion.

Text editors should mainly be used for editing code and other multiline syntactic data like markdown and JSON. For simpler use cases, consider using a [[text field]] instead.

## Properties

You can make a text editor without any custom options:

```Go
textcore.NewEditor(b)
```

You can set the text of a text editor:

```Go
textcore.NewEditor(b).Lines.SetString("Hello, world!")
```

You can set the highlighting language of a text editor:

```Go
textcore.NewEditor(b).Lines.SetLanguage(fileinfo.Go).SetString(`package main

func main() {
    fmt.Println("Hello, world!")
}
`)
```

You can set the text of a text editor from an embedded file:

```go
//go:embed file.go
var myFile embed.FS
```

```Go
errors.Log(textcore.NewEditor(b).Lines.OpenFS(myFile, "file.go"))
```

You can also set the text of a text editor directly from the system filesystem, but this is not recommended for files built into your app, since they will end up in a different location on different platforms:

```go
errors.Log(textcore.NewEditor(b).Lines.Open("file.go"))
```

You can make multiple text editors that edit the same underlying text buffer:

```Go
ls := lines.NewLines().SetString("Hello, world!")
textcore.NewEditor(b).SetLines(ls)
textcore.NewEditor(b).SetLines(ls)
```

## Events

You can detect when the user [[events#change]]s the content of a text editor and then exits it:

```Go
ed := textcore.NewEditor(b)
ed.OnChange(func(e events.Event) {
    core.MessageSnackbar(b, "OnChange: "+ed.Lines.String())
})
```

You can detect when the user makes any changes to the content of a text editor as they type ([[events#input]]):

```Go
ed := textcore.NewEditor(b)
ed.OnInput(func(e events.Event) {
    core.MessageSnackbar(b, "OnInput: "+ed.Lines.String())
})
```
