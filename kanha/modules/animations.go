package modules

import (
"time"

td "github.com/Kanha/Meow"

)

// runCommandAnimation adds a reaction and a smooth loading animation.
// Animation failures never stop the original command.
func runCommandAnimation(c *td.Client, m *td.Message, command string) {
if c == nil || m == nil {
return
}

// Try adding a reaction to the command message.
reaction, err := td.UnmarshalReactionType(
	[]byte(`{"@type":"reactionTypeEmoji","emoji":"💜"}`),
)

if err != nil {
	c.Logger.Error(
		"Command reaction parse failed",
		"command", command,
		"error", err,
	)
} else {
	err = c.AddMessageReaction(
		m.ChatID(),
		m.Id,
		reaction,
		&td.AddMessageReactionOpts{
			IsBig: true,
		},
	)

	if err != nil {
		c.Logger.Error(
			"Command reaction failed",
			"command", command,
			"chat_id", m.ChatID(),
			"message_id", m.Id,
			"error", err,
		)
	}
}

// Temporary loading message.
loader, err := m.ReplyText(
	c,
	"💜 𓆩 𝙆𝘼𝙉𝙃𝘼 𝙈𝙐𝙎𝙄𝘾 𓆪 💜\n\n⏳ 𝙎𝙩𝙖𝙧𝙩𝙞𝙣𝙜 𝙪𝙥...",
	nil,
)
if err != nil {
	c.Logger.Error(
		"Command animation start failed",
		"command", command,
		"error", err,
	)
	return
}

// Smooth, stylish KanhaMusic animation frames.
frames := []string{
	"💜 𓆩 𝙆𝘼𝙉𝙃𝘼 𝙈𝙐𝙎𝙄𝘾 𓆪 💜\n\n✧ Welcome to the divine vibes ✧\n\n⏳ 𝙎𝙩𝙖𝙧𝙩𝙞𝙣𝙜 𝙪𝙥...",

	"🦚 𓆩 𝙆𝘼𝙉𝙃𝘼 𝙈𝙐𝙎𝙄𝘾 𓆪 🦚\n\n♫ Bringing your favourite vibes ♫\n\n🔄 𝙇𝙤𝙖𝙙𝙞𝙣𝙜 𝙢𝙪𝙨𝙞𝙘 𝙬𝙤𝙧𝙡𝙙...",

	"🎶 𓆩 𝙆𝘼𝙉𝙃𝘼 𝙈𝙐𝙎𝙄𝘾 𓆪 🎶\n\n✦ Good music, good mood ✦\n\n⚡ 𝙋𝙧𝙚𝙥𝙖𝙧𝙞𝙣𝙜 𝙮𝙤𝙪𝙧 𝙢𝙚𝙣𝙪...",

	"✨ 𓆩 𝙆𝘼𝙉𝙃𝘼 𝙈𝙐𝙎𝙄𝘾 𓆪 ✨\n\n💜 Your music journey starts here\n\n✅ 𝙍𝙚𝙖𝙙𝙮 — 𝙇𝙚𝙩'𝙨 𝙫𝙞𝙗𝙚!",
}

for _, frame := range frames {
	time.Sleep(600 * time.Millisecond)

	if _, err := loader.EditText(c, frame, nil); err != nil {
		c.Logger.Error(
			"Command animation edit failed",
			"command", command,
			"error", err,
		)
		break
	}
}

// Remove the temporary loading message after animation.
if err := loader.Delete(c, true); err != nil {
	c.Logger.Error(
		"Command animation cleanup failed",
		"command", command,
		"error", err,
	)
}

}
