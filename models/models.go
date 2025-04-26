package models

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	ComponentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"createApplication": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "createApplication_" + i.Interaction.Member.User.ID,
					Title:    "Получение проходки",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "Nick_Minecraft",
									Label:       "Игровой ник",
									Style:       discordgo.TextInputShort,
									Placeholder: "Dipas_",
									Required:    true,
									MaxLength:   300,
									MinLength:   2,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "MyName",
									Label:       "Как вас зовут",
									Style:       discordgo.TextInputShort,
									Placeholder: "Ксюша",
									Required:    false,
									MaxLength:   300,
									MinLength:   2,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "MyDate",
									Label:       "Дата рождения и взораст",
									Style:       discordgo.TextInputShort,
									Placeholder: "5.10.2005/20",
									Required:    true,
									MaxLength:   300,
									MinLength:   2,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "MyCity",
									Label:       "Город",
									Style:       discordgo.TextInputShort,
									Placeholder: "Москва",
									Required:    true,
									MaxLength:   300,
									MinLength:   2,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "MyDescription",
									Label:       "Расскажите о себе",
									Style:       discordgo.TextInputParagraph,
									Placeholder: "Я занимаюсь спортом...",
									Required:    true,
									MaxLength:   300,
									MinLength:   2,
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Println("Model changeNikaName:", err)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Созданно.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
		"createTiсkets": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "createTicket_" + i.Interaction.Member.User.ID,
					Title:    "Разное",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "Nick_Minecraft",
									Label:       "Игровой ник",
									Style:       discordgo.TextInputShort,
									Placeholder: "Dipas_",
									Required:    false,
									MaxLength:   300,
									MinLength:   2,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:    "MyDescription",
									Label:       "Описание",
									Style:       discordgo.TextInputParagraph,
									Placeholder: "Помогите",
									Required:    true,
									MaxLength:   300,
									MinLength:   2,
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Println("Model changeNikaName:", err)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Созданно.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
	}
)
