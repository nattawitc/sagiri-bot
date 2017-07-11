package emoji

type Name string

const (
	Zero  Name = "zero"
	One   Name = "one"
	Two   Name = "two"
	Three Name = "three"
	Four  Name = "four"
	Five  Name = "five"
	Six   Name = "six"
	Seven Name = "seven"
	Eight Name = "eight"
	Nine  Name = "nine"
)

var (
	GetName map[string]Name
	GetCode map[Name]string
)

func init() {
	GetName = make(map[string]Name)
	GetCode = make(map[Name]string)

	AddEmoji(Zero, "\x30\xE2\x83\xA3")
	AddEmoji(One, "\x31\xE2\x83\xA3")
	AddEmoji(Two, "\x32\xE2\x83\xA3")
	AddEmoji(Three, "\x33\xE2\x83\xA3")
	AddEmoji(Four, "\x34\xE2\x83\xA3")
	AddEmoji(Five, "\x35\xE2\x83\xA3")
	AddEmoji(Six, "\x36\xE2\x83\xA3")
	AddEmoji(Seven, "\x37\xE2\x83\xA3")
	AddEmoji(Eight, "\x38\xE2\x83\xA3")
	AddEmoji(Nine, "\x39\xE2\x83\xA3")
}

func AddEmoji(name Name, code string) {
	GetName[code] = name
	GetCode[name] = code
}
