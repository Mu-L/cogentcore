+++
Categories = ["Widgets"]
+++

A **switch** is a [[widget]] that allows users to control a single bool value. See [[switches]] to handle multiple items at once.

## Properties

You can make a switch with no text:

```Go
core.NewSwitch(b)
```

You can add [[text]] to a switch:

```Go
core.NewSwitch(b).SetText("Remember me")
```

## Types

You can make a switch render as a checkbox:

```Go
core.NewSwitch(b).SetType(core.SwitchCheckbox).SetText("Remember me")
```

You can make a switch render as a radio button:

```Go
core.NewSwitch(b).SetType(core.SwitchRadioButton).SetText("Remember me")
```

## Events

You can detect when the user [[events#change]]s whether the switch is checked:

```Go
sw := core.NewSwitch(b).SetText("Remember me")
sw.OnChange(func(e events.Event) {
    core.MessageSnackbar(sw, fmt.Sprintf("Switch is %v", sw.IsChecked()))
})
```

You can make a switch read-only to prevent users from changing it:

```Go
core.NewSwitch(b).SetReadOnly(true)
```

## Styles

You can change the [[styles#font size]] of a switch:

```Go
sw := core.NewSwitch(b).SetText("Remember me")
sw.Styler(func(s *styles.Style) {
    s.Font.Size.Dp(25)
})
```

You can change the [[icon#icon size]] of a switch:

```Go
sw := core.NewSwitch(b).SetType(core.SwitchCheckbox).SetText("Remember me")
sw.Styler(func(s *styles.Style) {
    s.IconSize.Set(units.Dp(40))
})
```
