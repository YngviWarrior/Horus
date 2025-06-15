package main

import (
	"discord-bot/handlers"
	"discord-bot/jobs"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	timeAlive := time.Now()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Aviso: não foi possível carregar .env (normal no deploy em Railway)")
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

	go jobs.ResetData(dg, timeAlive)

	err = dg.Open()
	if err != nil {
		fmt.Println("Erro ao abrir conexão:", err)
		return
	}
	defer dg.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Encerrando bot...")
}
