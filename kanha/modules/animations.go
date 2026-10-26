
package modules

import (
	"time"

	td "github.com/Kanha/Meow"
)

// runCommandAnimation is a best-effort visual effect.
// Animation/reaction failures never stop the original command.
func runCommandAnimation(c *td.Client, m *td.Message, command string) {
	if c == nil || m == nil {
		return
	}

	// Try a large unicorn reaction. Unsupported reactions are ignored.
	reaction, err := td.UnmarshalReactionType(
		[]byte(`{"@type":"reactionTypeEmoji","emoji":"🦄"}`),
	)
	if err == nil {
		_ = c.AddMessageReaction(
			m.ChatID(),
			m.Id,
			reaction,
			&td.AddMessageReactionOpts{IsBig: true},
		)
	}

	loader, err := m.ReplyText(
		c,
		"✨ KANHA MUSIC\n⏳ Waking up...",
		nil,
	)
	if err != nil {
		return
	}

	frames := []string{
		"✨ KANHA MUSIC\n⚡ Initializing...",
		"💜 KANHA MUSIC\n🎶 Preparing your menu...",
		"🦄 KANHA MUSIC\n✅ Almost ready...",
	}

	for _, frame := range frames {
		time.Sleep(250 * time.Millisecond)

		if _, err := loader.EditText(c, frame, nil); err != nil {
			break
		}
	}

	// Remove the temporary loader, leaving the actual command response.
	_ = loader.Delete(c, true)
}
