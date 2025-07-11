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

	dg.AddHandler(messageHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Printf("[DISCORD] Discord bot is now running. Press CTRL+C to exit.\n")
	<-stop

	return nil
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	twitterRegex := regexp.MustCompile(`(https?:\/\/)?(?:www\.)?(\w*\.com)(\S+)`)

	matches := twitterRegex.FindAllStringSubmatch(m.Content, -1)

	if len(matches) > 0 {
		log.Printf("[DISCORD] Found Twitter links in message from %s: %v", m.Author.Username, matches)

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
	}
}
