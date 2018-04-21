package dhelpers

import (
	"io"
	"strings"

	"gitlab.com/project-d-collab/dhelpers/cache"

	"github.com/bwmarrin/discordgo"
)

// SendMessage sends a message to a specific channel, takes care of splitting and sanitising the content
func SendMessage(channelID, content string) (messages []*discordgo.Message, err error) {
	var message *discordgo.Message
	content = cleanDiscordContent(T(content))
	if len(content) > 2000 {
		for _, page := range autoPagify(content) {
			message, err = cache.GetDiscord().ChannelMessageSend(channelID, page)
			if err != nil {
				return messages, err
			}
			messages = append(messages, message)
		}
	} else {
		message, err = cache.GetDiscord().ChannelMessageSend(channelID, content)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// SendMessagef sends a message to a specific channel, takes care of splitting and sanitising the content, and replacing the fields
func SendMessagef(channelID, content string, fields ...interface{}) (messages []*discordgo.Message, err error) {
	var message *discordgo.Message
	content = cleanDiscordContent(Tf(content, fields...))
	if len(content) > 2000 {
		for _, page := range autoPagify(content) {
			message, err = cache.GetDiscord().ChannelMessageSend(channelID, page)
			if err != nil {
				return messages, err
			}
			messages = append(messages, message)
		}
	} else {
		message, err = cache.GetDiscord().ChannelMessageSend(channelID, content)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// SendMessagefc sends a message to a specific channel, takes care of splitting and sanitising the content, and replacing the fields, and applying pluralization
func SendMessagefc(channelID, content string, count int, fields ...interface{}) (messages []*discordgo.Message, err error) {
	var message *discordgo.Message
	content = cleanDiscordContent(Tfc(content, count, fields...))
	if len(content) > 2000 {
		for _, page := range autoPagify(content) {
			message, err = cache.GetDiscord().ChannelMessageSend(channelID, page)
			if err != nil {
				return messages, err
			}
			messages = append(messages, message)
		}
	} else {
		message, err = cache.GetDiscord().ChannelMessageSend(channelID, content)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// SendMessageBoxed sends a message to a specific channel, will put a box around it, takes care of splitting and sanitising the content
func SendMessageBoxed(channelID, content string) (messages []*discordgo.Message, err error) {
	var newMessages []*discordgo.Message
	content = cleanDiscordContent(T(content))
	for _, page := range autoPagify(content) {
		newMessages, err = SendMessage(channelID, "```"+page+"```")
		if err != nil {
			return messages, err
		}
		messages = append(messages, newMessages...)
	}
	return messages, nil
}

// SendEmbed sends an embed to a specific channel, takes care of splitting and sanitising the content
func SendEmbed(channelID string, embed *discordgo.MessageEmbed) (messages []*discordgo.Message, err error) {
	var message *discordgo.Message
	message, err = cache.GetDiscord().ChannelMessageSendEmbed(channelID, truncateEmbed(embed))
	if err != nil {
		return messages, err
	}
	messages = append(messages, message)
	return messages, nil
}

// SendFile sends a file to a specific channel, takes care of splitting and sanitising the content
func SendFile(channelID string, filename string, reader io.Reader, message string) (messages []*discordgo.Message, err error) {
	return SendComplex(channelID, &discordgo.MessageSend{File: &discordgo.File{Name: filename, Reader: reader}, Content: message})
}

// SendComplex sends a discordgo.MessageSend object to a specific channel, takes care of splitting and sanitising the content
func SendComplex(channelID string, data *discordgo.MessageSend) (messages []*discordgo.Message, err error) {
	var message *discordgo.Message
	if data.Embed != nil {
		data.Embed = truncateEmbed(data.Embed)
	}
	data.Content = cleanDiscordContent(data.Content)
	pages := autoPagify(data.Content)
	if len(pages) > 0 {
		for i, page := range pages {
			if i+1 < len(pages) {
				message, err = cache.GetDiscord().ChannelMessageSend(channelID, page)
			} else {
				data.Content = page
				message, err = cache.GetDiscord().ChannelMessageSendComplex(channelID, data)
			}
			if err != nil {
				return messages, err
			}
			messages = append(messages, message)
		}
	} else {
		message, err = cache.GetDiscord().ChannelMessageSendComplex(channelID, data)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// EditMessage edits a specific message, takes care of sanitising the content
func EditMessage(channelID, messageID, content string) (message *discordgo.Message, err error) {
	content = cleanDiscordContent(T(content))
	message, err = cache.GetDiscord().ChannelMessageEdit(channelID, messageID, content)
	if err != nil {
		return nil, err
	}
	return message, err
}

// EditMessagef edits a specific message, takes care of sanitising the content, and replacing the fields
func EditMessagef(channelID, messageID, content string, fields ...interface{}) (message *discordgo.Message, err error) {
	content = cleanDiscordContent(Tf(content, fields...))
	message, err = cache.GetDiscord().ChannelMessageEdit(channelID, messageID, content)
	if err != nil {
		return nil, err
	}
	return message, err
}

// EditMessagefc edits a specific message, takes care of sanitising the content, and replacing the fields, and applying pluralization
func EditMessagefc(channelID, messageID, content string, count int, fields ...interface{}) (message *discordgo.Message, err error) {
	content = cleanDiscordContent(Tfc(content, count, fields...))
	message, err = cache.GetDiscord().ChannelMessageEdit(channelID, messageID, content)
	if err != nil {
		return nil, err
	}
	return message, err
}

// EditEmbed edits a specific embed, takes care of sanitising the content
func EditEmbed(channelID, messageID string, embed *discordgo.MessageEmbed) (message *discordgo.Message, err error) {
	message, err = cache.GetDiscord().ChannelMessageEditEmbed(channelID, messageID, truncateEmbed(embed))
	if err != nil {
		return nil, err
	}
	return message, err
}

// EditComplex edits a specific message using a discordgo.MessageEdit object, takes care of sanitising the content
func EditComplex(data *discordgo.MessageEdit) (message *discordgo.Message, err error) {
	if data.Embed != nil {
		data.Embed = truncateEmbed(data.Embed)
	}
	if data.Content != nil {
		content := cleanDiscordContent(*data.Content)
		data.Content = &content
	}
	message, err = cache.GetDiscord().ChannelMessageEditComplex(data)
	if err != nil {
		return nil, err
	}
	return message, err
}

func pagify(text string, delimiter string) []string {
	result := make([]string, 0)
	textParts := strings.Split(text, delimiter)
	currentOutputPart := ""
	for _, textPart := range textParts {
		if len(currentOutputPart)+len(textPart)+len(delimiter) <= 1992 {
			if len(currentOutputPart) > 0 || len(result) > 0 {
				currentOutputPart += delimiter + textPart
			} else {
				currentOutputPart += textPart
			}
		} else {
			result = append(result, currentOutputPart)
			currentOutputPart = ""
			if len(textPart) <= 1992 {
				currentOutputPart = textPart
			}
		}
	}
	if currentOutputPart != "" {
		result = append(result, currentOutputPart)
	}
	return result
}

func autoPagify(text string) (pages []string) {
	for _, page := range pagify(text, "\n") {
		if len(page) <= 1992 {
			pages = append(pages, page)
		} else {
			for _, page := range pagify(page, ",") {
				if len(page) <= 1992 {
					pages = append(pages, page)
				} else {
					for _, page := range pagify(page, "-") {
						if len(page) <= 1992 {
							pages = append(pages, page)
						} else {
							for _, page := range pagify(page, " ") {
								if len(page) <= 1992 {
									pages = append(pages, page)
								} else {
									panic("unable to pagify text")
								}
							}
						}
					}
				}
			}
		}
	}
	return pages
}

func cleanDiscordContent(content string) (output string) {
	output = strings.Replace(content, "@everyone", "@"+ZeroWidthSpace+"everyone", -1)
	output = strings.Replace(output, "@here", "@"+ZeroWidthSpace+"here", -1)
	return output
}

// Applies Embed Limits to the given Embed
// Source: https://discordapp.com/developers/docs/resources/channel#embed-limits
func truncateEmbed(embed *discordgo.MessageEmbed) (result *discordgo.MessageEmbed) {
	if embed == nil || (&discordgo.MessageEmbed{}) == embed {
		return nil
	}
	if embed.Title != "" && len(embed.Title) > 256 {
		embed.Title = embed.Title[0:255] + "…"
	}
	if len(embed.Description) > 2048 {
		embed.Description = embed.Description[0:2047] + "…"
	}
	if embed.Footer != nil && len(embed.Footer.Text) > 2048 {
		embed.Footer.Text = embed.Footer.Text[0:2047] + "…"
	}
	if embed.Author != nil && len(embed.Author.Name) > 256 {
		embed.Author.Name = embed.Author.Name[0:255] + "…"
	}
	newFields := make([]*discordgo.MessageEmbedField, 0)
	for _, field := range embed.Fields {
		if field.Value == "" {
			continue
		}
		if len(field.Name) > 256 {
			field.Name = field.Name[0:255] + "…"
		}
		// TODO: better cutoff (at commas and stuff)
		if len(field.Value) > 1024 {
			field.Value = field.Value[0:1023] + "…"
		}
		newFields = append(newFields, field)
		if len(newFields) >= 25 {
			break
		}
	}
	embed.Fields = newFields

	if calculateFullEmbedLength(embed) > 6000 {
		if embed.Footer != nil {
			embed.Footer.Text = ""
		}
		if calculateFullEmbedLength(embed) > 6000 {
			if embed.Author != nil {
				embed.Author.Name = ""
			}
			if calculateFullEmbedLength(embed) > 6000 {
				embed.Fields = []*discordgo.MessageEmbedField{{}}
			}
		}
	}

	result = embed
	return result
}

func calculateFullEmbedLength(embed *discordgo.MessageEmbed) (count int) {
	count += len(embed.Title)
	count += len(embed.Description)
	if embed.Footer != nil {
		count += len(embed.Footer.Text)
	}
	if embed.Author != nil {
		count += len(embed.Author.Name)
	}
	for _, field := range embed.Fields {
		count += len(field.Name)
		count += len(field.Value)
	}
	return count
}
