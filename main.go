// main.go - FINAL VERSION

package main

import (
	"context"
	"database/sql"
	
	"log"
	"os"
	"os/signal"
	"path/filepath"
	
	"strings"
	
	"telegrarmchatbot/db"
	
	"telegrarmchatbot/internal/service"
	"telegrarmchatbot/internal/state"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	_ "github.com/lib/pq"
)

var (
	database       *sql.DB
	bookingService *service.BookingService
)

func main() {
	// Read token
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get working directory: %v", err)
	}

	filePath := filepath.Join(wd, "Token.txt")
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	token := strings.TrimSpace(string(contentBytes))

	// Connect to database
	database, err = db.Connect()
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize tables
	if err := db.InitTables(database); err != nil {
		log.Fatalf("unable to create tables: %v", err)
	}

	// Seed rooms
	if err := db.SeedRooms(database); err != nil {
		log.Fatalf("unable to seed rooms: %v", err)
	}

	// Initialize booking service
	bookingService = service.NewBookingService(database)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	// Register handlers
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/book", bot.MatchTypeExact, bookHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, cancelHandler)

	log.Println("Bot started successfully!")
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.Message == nil || update.Message.From == nil {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Unknown command. Type /help for available commands.",
	})
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Create or get user
	telegramID := update.Message.From.ID
	username := update.Message.From.Username
	fullName := update.Message.From.FirstName + " " + update.Message.From.LastName

	_, err := db.CreateOrGetUser(database, telegramID, username, fullName)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}

	welcomeText := `ສະບາຍດີ! Welcome to Room Booking Bot 🏢

Available commands:
/book - Book a meeting room
/cancel - Cancel your booking
/help - Show help message

Let's get started!`

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   welcomeText,
	})
}

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	helpText := `*Room Booking Bot Help* 🏢

*Available Commands:*
/book - Book a meeting room
/cancel - Cancel your booking
/help - Show this help message

*How to book:*
1. Type /book
2. View available time slots
3. Select a room and time
4. Enter meeting details
5. Confirm booking

*Rooms Available:*
- Room A
- Room B
- Room C

*Operating Hours:*
09:00 - 17:00 (1-hour slots)`

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      helpText,
		ParseMode: models.ParseModeMarkdown,
	})
}

func bookHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Generate today's timetable
	schedules, err := bookingService.GenerateTodayTimetable()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Sorry, unable to retrieve schedule. Please try again later.",
		})
		log.Printf("Error getting timetable: %v", err)
		return
	}

	// Format and send schedule
	message := bookingService.FormatTimetableMessage(schedules)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: models.ParseModeMarkdown,
	})

	// Start booking session
	userID := update.Message.From.ID
	state.Manager.StartBooking(userID)

	// Show room selection buttons
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🏢 Room A", CallbackData: "room | Room A"},
				{Text: "🏢 Room B", CallbackData: "room | Room B"},
				{Text: "🏢 Room C", CallbackData: "room | Room C"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Please select a room:",
		ReplyMarkup: keyboard,
	})
}



func cancelHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Get user
	telegramID := update.Message.From.ID
	user, err := db.GetUserByTelegramID(database, telegramID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error retrieving your information.",
		})
		log.Printf("Error getting user: %v", err)
		return
	}
	// Get user's bookings
	bookings, err := db.GetUserBookings(database, user.UserID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error retrieving your bookings.",
		})
		log.Printf("Error getting bookings: %v", err)
		return
	}

	// Format and send message
	message := bookingService.FormatUserBookings(bookings)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: models.ParseModeMarkdown,
	})

} // TODO: Add inline keyboard for cancellation
