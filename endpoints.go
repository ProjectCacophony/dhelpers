package dhelpers

import "github.com/bwmarrin/discordgo"

// SetDiscordEndpoints sets endpoints with a custom base discord host
// https://github.com/bwmarrin/discordgo/blob/master/endpoints.go#L181
func SetDiscordEndpoints(endpointDiscord string) {
	discordgo.EndpointDiscord = endpointDiscord
	discordgo.EndpointAPI = discordgo.EndpointDiscord + "api/v" + discordgo.APIVersion + "/"
	discordgo.EndpointGuilds = discordgo.EndpointAPI + "guilds/"
	discordgo.EndpointChannels = discordgo.EndpointAPI + "channels/"
	discordgo.EndpointUsers = discordgo.EndpointAPI + "users/"
	discordgo.EndpointGateway = discordgo.EndpointAPI + "gateway"
	discordgo.EndpointGatewayBot = discordgo.EndpointGateway + "/bot"
	discordgo.EndpointWebhooks = discordgo.EndpointAPI + "webhooks/"
}
