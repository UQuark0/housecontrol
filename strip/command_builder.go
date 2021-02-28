package strip

type Command byte
type Key byte
type Mode byte

const (
	ModeNoise     Mode = 0
	ModeRainbow   Mode = 1
	ModeEpileptic Mode = 2
	ModeTurnoff   Mode = 3
	ModeNight     Mode = 4
)

const (
	CommandGet     Command = 0
	CommandSet     Command = 1
	CommandSave    Command = 2
	CommandReset   Command = 3
	CommandPing    Command = 4
	CommandDefault Command = 5
)

const (
	KeyWidth      Key = 1 << 0
	KeySpeed      Key = 1 << 1
	KeyBrightness Key = 1 << 2
	KeyMode       Key = 1 << 3
)

type CommandBuilder struct {
	command    Command
	width      byte
	speed      byte
	brightness byte
	mode       Mode

	commandSet bool
	widthSet bool
	speedSet bool
	brightnessSet bool
	modeSet bool
}

func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{}
}

func (b *CommandBuilder) Set() *CommandBuilder {
	b.command = CommandSet
	b.commandSet = true
	return b
}

func (b *CommandBuilder) Get() *CommandBuilder {
	b.command = CommandGet
	b.commandSet = true
	return b
}

func (b *CommandBuilder) Ping() *CommandBuilder {
	b.command = CommandPing
	b.commandSet = true
	return b
}

func (b *CommandBuilder) Reset() *CommandBuilder {
	b.command = CommandReset
	b.commandSet = true
	return b
}

func (b *CommandBuilder) Default() *CommandBuilder {
	b.command = CommandDefault
	b.commandSet = true
	return b
}

func (b *CommandBuilder) Save() *CommandBuilder {
	b.command = CommandSave
	b.commandSet = true
	return b
}

func (b *CommandBuilder) allowValues() bool {
	return b.command == CommandSet || b.command == CommandGet || b.command == CommandSave || b.command == CommandDefault
}

func (b *CommandBuilder) Width(value ...byte) *CommandBuilder {
	if !b.allowValues() {
		panic("values are not allowed on this command")
	}
	if b.command == CommandSet {
		b.width = value[0]
	}
	b.widthSet = true
	return b
}

func (b *CommandBuilder) Speed(value ...byte) *CommandBuilder {
	if !b.allowValues() {
		panic("values are not allowed on this command")
	}
	if b.command == CommandSet {
		b.speed = value[0]
	}
	b.speedSet = true
	return b
}

func (b *CommandBuilder) Brightness(value ...byte) *CommandBuilder {
	if !b.allowValues() {
		panic("values are not allowed on this command")
	}
	if b.command == CommandSet {
		b.brightness = value[0]
	}
	b.brightnessSet = true
	return b
}

func (b *CommandBuilder) Mode(mode ...Mode) *CommandBuilder {
	if !b.allowValues() {
		panic("values are not allowed on this command")
	}
	if b.command == CommandSet {
		b.mode = mode[0]
	}
	b.modeSet = true
	return b
}

func (b *CommandBuilder) Build() []byte {
	var key Key
	values := make([]byte, 0)

	if b.widthSet {
		key |= KeyWidth
		if b.command == CommandSet {
			values = append(values, b.width)
		}
	}
	if b.speedSet {
		key |= KeySpeed
		if b.command == CommandSet {
			values = append(values, b.speed)
		}
	}
	if b.brightnessSet {
		key |= KeyBrightness
		if b.command == CommandSet {
			values = append(values, b.brightness)
		}
	}
	if b.modeSet {
		key |= KeyMode
		if b.command == CommandSet {
			values = append(values, byte(b.mode))
		}
	}

	command := byte(b.command << 4) | byte(key)

	return append([]byte{command}, values...)
}