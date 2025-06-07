package main

import (
	"discord-bot/handlers"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega o .env
	err := godotenv.Load()
	if err != nil {
		panic("Erro ao carregar .env")
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Erro ao criar sessão:", err)
		return
	}

	dg.AddHandler(handlers.OnBotReady)
	dg.AddHandler(handlers.OnVoiceUpdate)
	dg.AddHandler(handlers.OnMessageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent

	err = dg.Open()
	if err != nil {
		fmt.Println("Erro ao abrir conexão:", err)
		return
	}

	fmt.Println("Bot online. Pressione CTRL+C para sair.")
	select {}
}
