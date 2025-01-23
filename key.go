package gui

type (
	KeyboardKey uint16
	MouseKey    uint16
	Modifier    uint8
)

const (
	KeyF1 KeyboardKey = 0xFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	key_min // see terminfo
)

const (
	MouseLeft MouseKey = 65512 - iota
	MouseMiddle
	MouseRight
	MouseRelease
	MouseWheelUp
	MouseWheelDown
)

const (
	KeyCtrlTilde      KeyboardKey = 0x00
	KeyCtrl2          KeyboardKey = 0x00
	KeyCtrlSpace      KeyboardKey = 0x00
	KeyCtrlA          KeyboardKey = 0x01
	KeyCtrlB          KeyboardKey = 0x02
	KeyCtrlC          KeyboardKey = 0x03
	KeyCtrlD          KeyboardKey = 0x04
	KeyCtrlE          KeyboardKey = 0x05
	KeyCtrlF          KeyboardKey = 0x06
	KeyCtrlG          KeyboardKey = 0x07
	KeyBackspace      KeyboardKey = 0x08
	KeyCtrlH          KeyboardKey = 0x08
	KeyTab            KeyboardKey = 0x09
	KeyCtrlI          KeyboardKey = 0x09
	KeyCtrlJ          KeyboardKey = 0x0A
	KeyCtrlK          KeyboardKey = 0x0B
	KeyCtrlL          KeyboardKey = 0x0C
	KeyEnter          KeyboardKey = 0x0D
	KeyCtrlM          KeyboardKey = 0x0D
	KeyCtrlN          KeyboardKey = 0x0E
	KeyCtrlO          KeyboardKey = 0x0F
	KeyCtrlP          KeyboardKey = 0x10
	KeyCtrlQ          KeyboardKey = 0x11
	KeyCtrlR          KeyboardKey = 0x12
	KeyCtrlS          KeyboardKey = 0x13
	KeyCtrlT          KeyboardKey = 0x14
	KeyCtrlU          KeyboardKey = 0x15
	KeyCtrlV          KeyboardKey = 0x16
	KeyCtrlW          KeyboardKey = 0x17
	KeyCtrlX          KeyboardKey = 0x18
	KeyCtrlY          KeyboardKey = 0x19
	KeyCtrlZ          KeyboardKey = 0x1A
	KeyEsc            KeyboardKey = 0x1B
	KeyCtrlLsqBracket KeyboardKey = 0x1B
	KeyCtrl3          KeyboardKey = 0x1B
	KeyCtrl4          KeyboardKey = 0x1C
	KeyCtrlBackslash  KeyboardKey = 0x1C
	KeyCtrl5          KeyboardKey = 0x1D
	KeyCtrlRsqBracket KeyboardKey = 0x1D
	KeyCtrl6          KeyboardKey = 0x1E
	KeyCtrl7          KeyboardKey = 0x1F
	KeyCtrlSlash      KeyboardKey = 0x1F
	KeyCtrlUnderscore KeyboardKey = 0x1F
	KeySpace          KeyboardKey = 0x20
	KeyBackspace2     KeyboardKey = 0x7F
	KeyCtrl8          KeyboardKey = 0x7F
)

const (
	ModAlt Modifier = 1 << iota
	ModMotion
)
