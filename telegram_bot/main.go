package main

import (
	"log"
	"math/rand"
	"time"
	

	"github.com/hbagdi/go-unsplash/unsplash"
)

const (
	unsplashAccessKey = "YOUR_UNSPLASH_ACCESS_KEY"
)

func main() {
	unsplash := unsplash.New(unsplashUnauthenticatedConnection, unsplashAccessKey)

	// Create the Telegram bot
	bot, err := tgbotapi.NewBotAPI("YOUR_TELEGRAM_BOT_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	// Set up the update configuration
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	// Get updates from Telegram
	updates, err := bot.GetUpdatesChan(config)
	if err != nil {
		log.Fatal(err)
	}

	// Listen for incoming messages
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.IsCommand() || update.Message.Text == "image" {
			// Get a random photo from Unsplash
			photoURL, err := getRandomPhoto()
			if err != nil {
				log.Println(err)
				continue
			}

			// Send the photo to the user
			photoBytes, err := getImageBytes(photoURL)
			if err != nil {
				log.Println(err)
				continue
			}

			photo := tgbotapi.FileBytes{Name: "photo.jpg", Bytes: photoBytes}
			photoMsg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, photo)
			_, err = bot.Send(photoMsg)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

}
func getRandomPhoto() (string, error) {
	// Set the search options
	options := &unsplash.SearchOptions{
		Query:  "random",
		Orient: "landscape",
	}

	// Search for photos
	photos, err := unsplash.Search.Photos(options)
	if err != nil {
		return "", err
	}

	// Get a random photo from the search results
	rand.Seed(time.Now().UnixNano())
	photoIndex := rand.Intn(len(photos.Results))
	photo := photos.Results[photoIndex]

	// Return the photo URL
	return photo.Urls.Regular, nil
}
