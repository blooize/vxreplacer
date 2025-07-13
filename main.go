package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func main() {
	file, err := os.OpenFile("logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	// err = godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	discord_bot_token := os.Getenv("DISCORD_TOKEN")

	err = StartBot(discord_bot_token)
	if err != nil {
		log.Fatalf("[DISCORD] Error starting Discord bot: %v", err)
	}
	log.Printf("[DISCORD] Discord bot started successfully\n")
}

func StartBot(token string) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("[DISCORD] Error creating Discord session: %v", err)
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	err = dg.Open()
	if err != nil {
		log.Printf("[DISCORD] Error opening connection to Discord: %v", err)
		return fmt.Errorf("error opening connection to Discord: %w", err)
	}
	defer dg.Close()

	dg.AddHandler(twitterMessageHandler)
	dg.AddHandler(bskyMessageHandler)
	dg.AddHandler(instaMessageHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Printf("[DISCORD] Discord bot is now running. Press CTRL+C to exit.\n")
	<-stop

	return nil
}

func twitterMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	twitterRegex := regexp.MustCompile(`(https?:\/\/)?(?:www\.)?(\w*\.com)(\S+)`)

	matches := twitterRegex.FindAllStringSubmatch(m.Content, -1)

	if len(matches) > 0 {
		log.Printf("[DISCORD] Found links in message from %s: %v", m.Author.Username, matches)

		var vxlinks []string
		for _, match := range matches {
			log.Printf("[DISCORD] Twitter link found: %s", match[3])
			if match[2] == "twitter.com" || match[2] == "x.com" {
				vxlink := fmt.Sprintf("https://vxtwitter.com%s", match[3])
				vxlinks = append(vxlinks, vxlink)
			} else {
				return
			}
		}

		if len(vxlinks) > 0 {
			response := fmt.Sprintf("\n%s", strings.Join(vxlinks, "\n"))
			_, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
			if err != nil {
				log.Printf("[DISCORD] Error sending reply: %v", err)
			}
		}
		m.Message.Flags = 1 << 2
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: m.ChannelID,
			ID:      m.Message.ID,
			Flags:   m.Message.Flags,
		})
		if err != nil {
			log.Printf("[DISCORD] Error editing message flags: %v", err)
		}
	}
}

func bskyMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	bskyRegex := regexp.MustCompile(`(https?:\/\/)?(?:www\.)?(bsky\.app|bsky\.social)(\S+)`)

	matches := bskyRegex.FindAllStringSubmatch(m.Content, -1)

	if len(matches) > 0 {
		log.Printf("[DISCORD] Found BlueSky links in message from %s: %v", m.Author.Username, matches)

		var bskx []string
		for _, match := range matches {
			log.Printf("[DISCORD] BlueSky link found: %s", match[3])
			if match[2] == "bsky.app" || match[2] == "bsky.social" {
				bskxlink := fmt.Sprintf("https://bskx.app%s", match[3])
				bskx = append(bskx, bskxlink)
			} else {
				return
			}
		}

		if len(bskx) > 0 {
			response := fmt.Sprintf("\n%s", strings.Join(bskx, "\n"))
			_, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
			if err != nil {
				log.Printf("[DISCORD] Error sending reply: %v", err)
			}
		}
		m.Message.Flags = 1 << 2
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: m.ChannelID,
			ID:      m.Message.ID,
			Flags:   m.Message.Flags,
		})
		if err != nil {
			log.Printf("[DISCORD] Error editing message flags: %v", err)
		}
	}
}

func instaMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	instaRegex := regexp.MustCompile(`(https?:\/\/)?(?:www\.)?(instagram\.com)(\S+)`)
	ddregex := regexp.MustCompile(`(https?:\/\/)?(?:www\.)?(ddinstagram\.com)(\S+)`)

	matches := instaRegex.FindAllStringSubmatch(m.Content, -1)
	ddmatches := ddregex.FindAllString(m.Content, -1)

	if len(matches) > 0 && len(ddmatches) == 0 {
		log.Printf("[DISCORD] Found Instagram links in message from %s: %v", m.Author.Username, matches)

		var instaLinks []string
		for _, match := range matches {
			log.Printf("[DISCORD] Instagram link found: %s", match[3])
			if match[2] == "instagram.com" && match[2] != "ddinstagram.com" {
				instaLink := fmt.Sprintf("https://ddinstagram.com%s", match[3])
				instaLinks = append(instaLinks, instaLink)
			} else {
				return
			}
		}

		if len(instaLinks) > 0 {
			response := fmt.Sprintf("\n%s", strings.Join(instaLinks, "\n"))
			_, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
			if err != nil {
				log.Printf("[DISCORD] Error sending reply: %v", err)
			}
		}
		m.Message.Flags = 1 << 2
		_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: m.ChannelID,
			ID:      m.Message.ID,
			Flags:   m.Message.Flags,
		})
		if err != nil {
			log.Printf("[DISCORD] Error editing message flags: %v", err)
		}
	}
}
