package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	webhookSecret   []byte
	discordSession  *discordgo.Session
	discordChannel  string
)

func init() {
	loadEnvFile("../.env")
	loadEnvFile("/project/.env")

	secret := getEnv("EVENT_SECRET", "")
	if secret == "" {
		log.Fatal("EVENT_SECRET is not set in .env or environment")
	}
	webhookSecret = []byte(secret)

	botToken := getEnv("DISCORD_BOT_TOKEN", "")
	if botToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set in .env or environment")
	}

	discordChannel = getEnv("DISCORD_CHANNEL_ID", "")
	if discordChannel == "" {
		log.Fatal("DISCORD_CHANNEL_ID is not set in .env or environment")
	}

	var err error
	discordSession, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	if err = discordSession.Open(); err != nil {
		log.Fatalf("Failed to connect Discord bot: %v", err)
	}

	log.Println("Discord bot connected successfully")
}

func main() {
	defer discordSession.Close()

	http.HandleFunc("/webhook", handleWebhook)

	port := getEnv("PORT", "8080")
	log.Printf("Webhook service running on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func sendToDiscord(embed *discordgo.MessageEmbed) {
	_, err := discordSession.ChannelMessageSendEmbed(discordChannel, embed)
	if err != nil {
		log.Printf("Failed to send Discord message: %v", err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("X-Hub-Signature-256")
	if !verifySignature(signature, body) {
		http.Error(w, "Invalid signature", http.StatusForbidden)
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	log.Printf("Received GitHub Event: %s", event)

	switch event {
	case "ping":
		log.Println("Initial connection (Ping) successful!")
		sendToDiscord(&discordgo.MessageEmbed{
			Title:       "🏓 GitHub Webhook Connected",
			Description: "Ping received. Webhook is live and connected!",
			Color:       0x57F287,
		})

	case "push":
		var payload PushPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Failed to parse push JSON: %v", err)
			break
		}
		branch := strings.TrimPrefix(payload.Ref, "refs/heads/")
		log.Printf("Repository: %s | Branch: %s | Pusher: %s", payload.Repository.FullName, branch, payload.Pusher.Name)

		commitLines := ""
		for i, c := range payload.Commits {
			if i >= 5 {
				commitLines += fmt.Sprintf("_...and %d more commit(s)_", len(payload.Commits)-5)
				break
			}
			short := c.ID
			if len(short) > 7 {
				short = short[:7]
			}
			commitLines += fmt.Sprintf("• [`%s`](%s) %s — **%s**\n", short, c.URL, c.Message, c.Author.Name)
		}
		if commitLines == "" {
			commitLines = "_No commit details available_"
		}

		sendToDiscord(&discordgo.MessageEmbed{
			Title: fmt.Sprintf("🚀 Push to %s", payload.Repository.FullName),
			URL:   payload.Repository.URL,
			Color: 0x5865F2,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Branch", Value: fmt.Sprintf("`%s`", branch), Inline: true},
				{Name: "Pusher", Value: payload.Pusher.Name, Inline: true},
				{Name: fmt.Sprintf("Commits (%d)", len(payload.Commits)), Value: commitLines},
			},
			Footer: &discordgo.MessageEmbedFooter{Text: "GitHub Push Event"},
		})

		go runMakeFullInstall(payload.Repository.FullName, branch)

	default:
		log.Printf("Event '%s' is not specifically handled.", event)
		sendToDiscord(&discordgo.MessageEmbed{
			Title:       fmt.Sprintf("📦 GitHub Event: %s", event),
			Description: "An unhandled GitHub event was received.",
			Color:       0xFEE75C,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event received successfully"))
}

func runMakeFullInstall(repo, branch string) {
	projectDir := "/project"

	log.Printf("Starting make full-install in %s for %s@%s", projectDir, repo, branch)
	sendToDiscord(&discordgo.MessageEmbed{
		Title:       "⚙️ Deploy Started",
		Description: fmt.Sprintf("Running `make full-install` for **%s** on branch `%s`...", repo, branch),
		Color:       0xFEE75C,
	})

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("make", "full-install")
	cmd.Dir = projectDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("make full-install failed: %v\nSTDERR: %s", err, stderr.String())
		sendToDiscord(&discordgo.MessageEmbed{
			Title:       "❌ Deploy Failed",
			Description: fmt.Sprintf("**%s** @ `%s`", repo, branch),
			Color:       0xED4245,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Error", Value: fmt.Sprintf("```%s```", truncate(stderr.String(), 1000))},
			},
			Footer: &discordgo.MessageEmbedFooter{Text: "make full-install"},
		})
		return
	}

	log.Printf("make full-install succeeded:\n%s", stdout.String())
	sendToDiscord(&discordgo.MessageEmbed{
		Title:       "✅ Deploy Successful",
		Description: fmt.Sprintf("**%s** @ `%s` deployed successfully!", repo, branch),
		Color:       0x57F287,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Output", Value: fmt.Sprintf("```%s```", truncate(stdout.String(), 1000))},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: "make full-install"},
	})
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[len(s)-max:]
}

func verifySignature(signatureHeader string, body []byte) bool {
	if signatureHeader == "" || len(signatureHeader) < 8 {
		return false
	}

	receivedHex := signatureHeader[7:]

	mac := hmac.New(sha256.New, webhookSecret)
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedHex := hex.EncodeToString(expectedMAC)

	return hmac.Equal([]byte(receivedHex), []byte(expectedHex))
}

func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

type CommitAuthor struct {
	Name string `json:"name"`
}

type Commit struct {
	ID      string       `json:"id"`
	Message string       `json:"message"`
	URL     string       `json:"url"`
	Author  CommitAuthor `json:"author"`
}

type PushPayload struct {
	Ref     string   `json:"ref"`
	Commits []Commit `json:"commits"`
	Pusher  struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Repository struct {
		FullName string `json:"full_name"`
		URL      string `json:"html_url"`
	} `json:"repository"`
}
