package models

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gorcon/rcon"
)

func CreateMessageTicket(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "HelpMenu" && m.Author.ID == "1075427758898950174" {
		messages, err := s.ChannelMessages(m.ChannelID, 1, "", "", "")
		if err != nil {
			log.Println(err)
		}

		if len(messages) != 0 {
			messageIDs := make([]string, len(messages))
			for i, msg := range messages {
				messageIDs[i] = msg.ID
			}

			err = s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
			if err != nil {
				log.Println(err)
			}
		}

		embed := []*discordgo.MessageEmbed{
			{
				Title:       "Поддержка",
				Description: "Что бы получить проходку нажмите на **Создать Заявку**\nТак же можно создать тикет поддержки  **Создать тикет**",
				Color:       0x00b0f4, // Используем шестнадцатеричное значение цвета
			},
		}

		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Content: "", Embeds: embed, Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "✅",
						},
						Label:    "Создать заявку",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						CustomID: "createApplication",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "📝",
						},
						Label:    "Создать тикет",
						Style:    discordgo.SecondaryButton,
						Disabled: false,
						CustomID: "createTiсkets",
					},
				},
			},
		}})
	}
}

func ButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userName []string

	if strings.HasPrefix(i.MessageComponentData().CustomID, "addMemberListServer") {
		for r := 0; r < len(i.Member.Roles); r++ {
			if i.Member.Roles[r] == "1168929755647717485" {
				re := regexp.MustCompile(`addMemberListServer_([a-zA-Z0-9_.]+)`)
				userName = re.FindStringSubmatch(i.MessageComponentData().CustomID)

				conn, err := rcon.Dial("116.202.214.243:24435", "B6PqNhNxtkGuQglDQ3igHkFDR8NQ_g")
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Восстановление после паники:", r)
					}
				}()

				conn.Execute(fmt.Sprintf("easywl add %v", userName[1]))

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Вы добавили в вайтлист: " + userName[1],
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}

			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      i.Message.ID,
				Channel: i.ChannelID,
				Components: &[]discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "❌",
						},
						Label:    "Закрыть тикет1",
						Style:    discordgo.DangerButton,
						Disabled: true,
						CustomID: "closeTicket",
					},
				},
			})
		}
	} else {
		log.Println(324243)
	}

	switch i.MessageComponentData().CustomID {
	case "closeTicket":
		s.ChannelDelete(i.ChannelID)
	case "accept":
		acceptOrNot := false
		for r := 0; r < len(i.Member.Roles); r++ {
			if i.Member.Roles[r] == "1168929755647717485" {
				embed := []*discordgo.MessageEmbed{
					{
						Title:       "Принятое обращение",
						Description: fmt.Sprintf("Ваш администратор: <@%v>", i.Member.User.ID),
						Color:       0x00b0f4,
					},
				}
				s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{Content: "", Embeds: embed})
				s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:      i.Message.ID,
					Channel: i.ChannelID,
					Components: &[]discordgo.MessageComponent{
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "❌",
							},
							Label:    "Закрыть тикет1",
							Style:    discordgo.DangerButton,
							Disabled: false,
							CustomID: "closeTicket",
						},
					},
				})
				if !acceptOrNot {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Вы приняли, пока всё в бете",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
				}
				acceptOrNot = true
				return
			}
		}

		if !acceptOrNot {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Вы не администратор",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	}
}

func InteractionModal(data discordgo.ModalSubmitInteractionData, s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch {
	case strings.HasPrefix(data.CustomID, "createApplication"):
		channel, _ := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
			Name:     fmt.Sprintf("проходка-%v", i.Member.User.GlobalName),
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: "1168944699378258023",
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    "1168929755647717485",
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel,
				},
				{
					ID:   i.GuildID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
				{
					ID:    i.Member.User.ID,
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: discordgo.PermissionViewChannel,
				},
			},
		})

		embed := []*discordgo.MessageEmbed{
			{
				Title:       "Поддержка",
				Description: "Благодарим вас за обращение в службу поддержки.\nПожалуйста, опишите вашу проблему и ждите ответа.",
				Color:       0x00b0f4,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "\nРежиме работы службы технической поддержки",
						Value:  "Будние дни с 10:00 до 22:00 по МСК. В остальное время сотрудники работают в свободном графике и по собственному желанию, — нам тоже необходим отдых. Рассчитываем на Ваше понимание!",
						Inline: false,
					},
					{
						Name:   "\nВажно",
						Value:  "Администратор никогда не попросит ваш пароль.\nАдминистрации Запрещено употреблять ненормативную лексику, если это произошло то подайте жалобу на администратора.\nРешение высший администрации не обсуждаются.",
						Inline: false,
					},
					{
						Name:   "\nИнформация",
						Value:  fmt.Sprintf("Открыл: <@%v>", i.Member.User.ID),
						Inline: false,
					},
				},
			},
		}

		s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{Content: fmt.Sprintf("<@&%v> <@&%v>", i.GuildID, "1168929755647717485"), Embeds: embed})

		embed = []*discordgo.MessageEmbed{
			{
				Color: 0x00b0f4,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Игровой ник",
						Value:  data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Inline: false,
					},
					{
						Name:   "Как вас зовут",
						Value:  data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Inline: false,
					},
					{
						Name:   "Дата рождения",
						Value:  data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Inline: false,
					},
					{
						Name:   "Город",
						Value:  data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Inline: false,
					},
					{
						Name:   "Расскажите о себе",
						Value:  data.Components[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Inline: false,
					},
				},
			},
		}
		s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{Content: "", Embeds: embed, Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "✅",
						},
						Label:    "Принять",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						CustomID: "accept",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "❌",
						},
						Label:    "Закрыть тикет",
						Style:    discordgo.DangerButton,
						Disabled: false,
						CustomID: "closeTicket",
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "✅",
						},
						Label:    "Добавить в список",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						CustomID: "addMemberListServer_" + data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
					},
				},
			},
		}})
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Созданно <#%v>", channel.ID),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	case strings.HasPrefix(data.CustomID, "createTicket"):

		channel, _ := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
			Name:     fmt.Sprintf("разное-%v", i.Member.User.GlobalName),
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: "1168944699378258023",
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    "1168929755647717485",
					Type:  discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel,
				},
				{
					ID:   i.GuildID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
				{
					ID:    i.Member.User.ID,
					Type:  discordgo.PermissionOverwriteTypeMember,
					Allow: discordgo.PermissionViewChannel,
				},
			},
		})

		embed := []*discordgo.MessageEmbed{
			{
				Title:       "Поддержка",
				Description: "Благодарим вас за обращение в службу поддержки.\nПожалуйста, опишите вашу проблему и ждите ответа.",
				Color:       0x00b0f4,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "\nРежиме работы службы технической поддержки",
						Value:  "Будние дни с 10:00 до 22:00 по МСК. В остальное время сотрудники работают в свободном графике и по собственному желанию, — нам тоже необходим отдых. Рассчитываем на Ваше понимание!",
						Inline: false,
					},
					{
						Name:   "\nВажно",
						Value:  "Администратор никогда не попросит ваш пароль.\nАдминистрации Запрещено употреблять ненормативную лексику, если это произошло то подайте жалобу на администратора.\nРешение высший администрации не обсуждаются.",
						Inline: false,
					},
					{
						Name:   "\nИнформация",
						Value:  fmt.Sprintf("Открыл: <@%v>", i.Member.User.ID),
						Inline: false,
					},
				},
			},
		}

		s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{Content: fmt.Sprintf("<@&%v> <@&%v>", i.GuildID, "1168929755647717485"), Embeds: embed})

		embed = []*discordgo.MessageEmbed{
			{
				Color: 0x00b0f4,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Игровой ник",
						Value:  data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Inline: false,
					},
					{
						Name:   "Описание",
						Value:  data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
						Inline: false,
					},
				},
			},
		}
		s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{Content: "", Embeds: embed, Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "✅",
						},
						Label:    "Принять",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						CustomID: "accept",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "❌",
						},
						Label:    "Закрыть тикет",
						Style:    discordgo.DangerButton,
						Disabled: false,
						CustomID: "closeTicket",
					},
				},
			},
		}})
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Созданно <#%v>", channel.ID),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
