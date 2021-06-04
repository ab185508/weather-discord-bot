package discordbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotID string
var GoBot *discordgo.Session
var Log int = 0

func Start() {
	GoBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := GoBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID
	GoBot.AddHandler(messageHandler)

	err = GoBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	mux_start()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, config.BotPrefix) {
		mid := m.Reference().MessageID

		if m.Author.ID == BotID {
			return
		}

		if m.Content == "!ping" {
			_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
		}

		if m.Content == "!commands" {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Available Commands: \n !ping: A command to get a test reply from weather-bot. \n !temp (cityname): A command to get the temperature in a city of choice. \n")
		}

		if strings.Contains(m.Content, "!temp") {
			whole := strings.Fields(m.Content)

			if len(whole) > 1 {
				city := whole[1]
				url := "http://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + APIKey
				dat, error := api(url)
				if error == 0 {
					temp, f := msi_to_s(dat)
					_, _ = s.ChannelMessageSend(m.ChannelID, ("The temperature is " + temp + "F in " + city))

					if f > 80 {
						err := s.MessageReactionAdd(m.ChannelID, mid, "\U0001F525")
						if err != nil {
							fmt.Println(err.Error())
						}
					}
					if f >= 59 && f <= 79 {
						err := s.MessageReactionAdd(m.ChannelID, mid, "\U0001F610")
						if err != nil {
							fmt.Println(err.Error())
						}
					}
					if f < 59 && f >= 41 {
						err := s.MessageReactionAdd(m.ChannelID, mid, "\U0001F9E5")
						if err != nil {
							fmt.Println(err.Error())
						}
					}
					if f < 40 {
						err := s.MessageReactionAdd(m.ChannelID, mid, "\U000026C4")
						if err != nil {
							fmt.Println(err.Error())
						}
					}
					Log++
					lg := strconv.Itoa(Log)
					mains(city, temp, lg)
				} else {
					_, _ = s.ChannelMessageSend(m.ChannelID, ("Invalid city, please try again."))
					return
				}
			} else {
				return
			}
		}

	}

}

func api(url string) (map[string]interface{}, int) {
	response, err := http.Get(url)
	error := 0
	dat := make(map[string]interface{})
	if err != nil {
		fmt.Println("fix the code dude")
		fmt.Println(err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		if err := json.Unmarshal(data, &dat); err != nil {
			fmt.Println("oh no")
		}

		if len(dat) == 2 {
			error = 1
			return dat, error
		}
	}
	return dat, error
}

func msi_to_s(json map[string]interface{}) (string, float64) {
	whole := json["main"]
	conv := whole.(map[string]interface{})
	temp_i := conv["temp"]
	temp_s := fmt.Sprintf("%v", temp_i)
	temp_f, _ := strconv.ParseFloat(temp_s, 64)
	k_to_f := ((temp_f - 273.15) * 1.8) + 32
	temp := fmt.Sprintf("%f", k_to_f)
	return temp, k_to_f
}
