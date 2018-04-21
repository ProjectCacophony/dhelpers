package dhelpers

import (
	"regexp"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

// EventType is the Type used for the EventContainer
type EventType string

// defines discordgo event types for the EventContainer
const (
	ChannelCreateEventType            EventType = "CHANNEL_CREATE"
	ChannelDeleteEventType                      = "CHANNEL_DELETE"
	ChannelPinsUpdateEventType                  = "CHANNEL_PINS_UPDATE"
	ChannelUpdateEventType                      = "CHANNEL_UPDATE"
	GuildBanAddEventType                        = "GUILD_BAN_ADD"
	GuildBanRemoveEventType                     = "GUILD_BAN_REMOVE"
	GuildCreateEventType                        = "GUILD_CREATE"
	GuildDeleteEventType                        = "GUILD_DELETE"
	GuildEmojisUpdateEventType                  = "GUILD_EMOJIS_UPDATE"
	GuildMemberAddEventType                     = "GUILD_MEMBER_ADD"
	GuildMemberRemoveEventType                  = "GUILD_MEMBER_REMOVE"
	GuildMemberUpdateEventType                  = "GUILD_MEMBER_UPDATE"
	GuildMembersChunkEventType                  = "GUILD_MEMBERS_CHUNK"
	GuildRoleCreateEventType                    = "GUILD_ROLE_CREATE"
	GuildRoleDeleteEventType                    = "GUILD_ROLE_DELETE"
	GuildRoleUpdateEventType                    = "GUILD_ROLE_UPDATE"
	GuildUpdateEventType                        = "GUILD_UPDATE"
	MessageCreateEventType                      = "MESSAGE_CREATE"
	MessageDeleteEventType                      = "MESSAGE_DELETE"
	MessageReactionAddEventType                 = "MESSAGE_REACTION_ADD"
	MessageReactionRemoveEventType              = "MESSAGE_REACTION_REMOVE"
	MessageReactionRemoveAllEventType           = "MESSAGE_REACTION_REMOVE_ALL"
	MessageUpdateEventType                      = "MESSAGE_UPDATE"
	PresenceUpdateEventType                     = "PRESENCE_UPDATE"
	//GuildIntegrationsUpdateEventType            = "GUILD_INTEGRATIONS_UPDATE"
	//PresencesReplaceEventType         = "PRESENCES_REPLACE"
	//ReadyEventType                    = "READY"
	//RelationshipAddEventType          = "RELATIONSHIP_ADD"
	//RelationshipRemoveEventType       = "RELATIONSHIP_REMOVE"
	//ResumedEventType                  = "RESUMED"
	//TypingStartEventType              = "TYPING_START"
	//UserGuildSettingsUpdateEventType  = "USER_GUILD_SETTINGS_UPDATE"
	//UserNoteUpdateEventType           = "USER_NOTE_UPDATE"
	//UserSettingsUpdateEventType       = "USER_SETTINGS_UPDATE"
	//UserUpdateEventType               = "USER_UPDATE"
	//VoiceServerUpdateEventType        = "VOICE_SERVER_UPDATE"
	//VoiceStateUpdateEventType         = "VOICE_STATE_UPDATE"
)

// Routing JSON Config
type rawRoutingEntryContainer struct {
	Module []rawRoutingEntry
}

type rawRoutingEntry struct {
	Active      bool
	Always      bool // if true: will run even if there have been previous (higher priority) matches
	AllowBots   bool // if set to true, will trigger for messages by bots
	AllowMyself bool // if set to true, will trigger for messages by this bot itself
	AllowDM     bool

	Events      []EventType
	Module      string
	Destination string
	Requirement []rawRoutingRequirementEntry // will only get matched with EventTypeMessageCreate, EventTypeMessageUpdate, or EventTypeMessageDelete, will match everything if slice is empty
	Priority    int                          // higher runs before lower
}

type rawRoutingRequirementEntry struct {
	Beginning          []string // can be empty, will match all
	Regex              string   // can be empty, will match all
	DoNotPrependPrefix bool     // if false, prepends guild prefix to regex
	CaseSensitive      bool     // prepends (?i) to regex on go, language dependent#
	Alias              string
}

// RoutingRule is a a compiled routing rule used for matching
type RoutingRule struct {
	Event              EventType
	Module             string
	DestinationMain    string
	DestinationSub     string
	Beginning          string
	Alias              string
	Regex              *regexp.Regexp
	DoNotPrependPrefix bool
	CaseSensitive      bool
	Always             bool
	AllowBots          bool
	AllowMyself        bool
	AllowDM            bool
}

// GetRoutings returns a sorted slice (by priority) with all rules
func GetRoutings() (routingRules []RoutingRule, err error) {
	// read and unmarshal config from file
	// TODO: load from S3 instead
	var rawRoutingContainer rawRoutingEntryContainer
	_, err = toml.DecodeFile("routing.toml", &rawRoutingContainer)
	if err != nil {
		return nil, err
	}

	// group rules by priorities
	rawEntriesByPriority := make(map[int][]rawRoutingEntry)

	for _, rawRoutingEntry := range rawRoutingContainer.Module {
		rawEntriesByPriority[rawRoutingEntry.Priority] = append(
			rawEntriesByPriority[rawRoutingEntry.Priority], rawRoutingEntry,
		)
	}

	// sort entries
	var keys []int
	for k := range rawEntriesByPriority {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
		keys[i], keys[j] = keys[j], keys[i]
	}

	// generated compiled routings
	for _, k := range keys {
		for _, rawRule := range rawEntriesByPriority[k] {
			// skip disabled rules
			if !rawRule.Active {
				continue
			}
			// skip empty event slices
			if rawRule.Events == nil || len(rawRule.Events) < 0 {
				continue
			}
			// skip empty or invalid destinations
			if rawRule.Destination == "" || !strings.Contains(rawRule.Destination, "/") {
				continue
			}
			// generate route for each type
			for _, ruleType := range rawRule.Events {
				parts := strings.SplitN(rawRule.Destination, "/", 2)
				newEntry := RoutingRule{
					Event:           ruleType,
					Module:          rawRule.Module,
					DestinationMain: parts[0],
					DestinationSub:  parts[1],
					Always:          rawRule.Always,
					AllowMyself:     rawRule.AllowMyself,
					AllowBots:       rawRule.AllowBots,
					AllowDM:         rawRule.AllowDM,

					Beginning:          "",
					Regex:              nil,
					DoNotPrependPrefix: false,
					CaseSensitive:      false,
				}
				if (ruleType == MessageCreateEventType ||
					ruleType == MessageUpdateEventType ||
					ruleType == MessageDeleteEventType) &&
					rawRule.Requirement != nil && len(rawRule.Requirement) > 0 {
					for _, requirement := range rawRule.Requirement {
						if requirement.Beginning != nil && len(requirement.Beginning) > 0 {
							for _, beginning := range requirement.Beginning {
								newEntryCopy := newEntry
								newEntryCopy.Beginning = beginning
								if requirement.Regex != "" {
									if requirement.CaseSensitive {
										newEntryCopy.Regex = regexp.MustCompile(requirement.Regex)
									} else {
										newEntryCopy.Regex = regexp.MustCompile("(?i)" + requirement.Regex)
									}
								}
								newEntryCopy.DoNotPrependPrefix = requirement.DoNotPrependPrefix
								newEntryCopy.CaseSensitive = requirement.CaseSensitive
								newEntryCopy.Alias = requirement.Alias
								routingRules = append(routingRules, newEntryCopy)
							}
						} else {
							newEntryCopy := newEntry
							if requirement.Regex != "" {
								if requirement.CaseSensitive {
									newEntryCopy.Regex = regexp.MustCompile(requirement.Regex)
								} else {
									newEntryCopy.Regex = regexp.MustCompile("(?i)" + requirement.Regex)
								}
							}
							newEntryCopy.DoNotPrependPrefix = requirement.DoNotPrependPrefix
							newEntryCopy.CaseSensitive = requirement.CaseSensitive
							newEntryCopy.Alias = requirement.Alias
							routingRules = append(routingRules, newEntryCopy)
						}
					}
				} else {
					routingRules = append(routingRules, newEntry)
				}
			}
		}
	}

	return routingRules, nil
}

// RoutingMatchMessage checks if a message content matches the requirements of the routing rule
func RoutingMatchMessage(routingEntry RoutingRule, author, bot *discordgo.User, channel *discordgo.Channel, content string, args []string, prefix string) (match bool) {
	// ignore bots?
	if !routingEntry.AllowBots {
		if author.Bot {
			return false
		}
	}
	// ignore itself?
	if !routingEntry.AllowMyself {
		if author.ID == bot.ID {
			return false
		}
	}
	// DMs?
	if !routingEntry.AllowDM {
		if channel.Type == discordgo.ChannelTypeDM {
			return false
		}
	}
	// check prefix if should check
	if !routingEntry.DoNotPrependPrefix {
		if prefix == "" {
			return false
		}
	}
	// match beginning if beginning is set
	if routingEntry.Beginning != "" {
		if routingEntry.CaseSensitive {
			if args[0] != routingEntry.Beginning {
				return false
			}
		} else {
			if strings.ToLower(args[0]) != strings.ToLower(routingEntry.Beginning) {
				return false
			}
		}
	}
	// match regex if regex is set
	if routingEntry.Regex != nil {
		matchContent := content
		if !routingEntry.DoNotPrependPrefix {
			matchContent = strings.TrimSpace(strings.TrimLeft(content, prefix))
		}
		if !routingEntry.Regex.MatchString(matchContent) {
			return false
		}
	}

	return true
}

// GetMessageArguments trims the prefix and returns all arguments, including the command, and the prefix used
func GetMessageArguments(content string, prefixes []string) (args []string, prefix string) {
	for _, possiblePrefix := range prefixes {
		if strings.HasPrefix(content, possiblePrefix) {
			content = strings.TrimLeft(content, possiblePrefix)
			prefix = possiblePrefix
			break
		}
	}

	args, err := ToArgv(content)
	if err == nil {
		return args, prefix
	}

	return []string{content}, prefix
}

// ContainerDestinations figures out the correct destinations for an event container
func ContainerDestinations(session *discordgo.Session, routingConfig []RoutingRule, container EventContainer) (lambdaDestinations, processorDestinations, aliases []string) {
	var handled int

	for _, routingEntry := range routingConfig {
		if handled > 0 && !routingEntry.Always {
			continue
		}

		if container.Type != routingEntry.Event {
			continue
		}

		// check requirements
		if container.Type == MessageCreateEventType {
			channel, err := session.State.Channel(container.MessageCreate.ChannelID)
			if err != nil {
				continue
			}

			if !RoutingMatchMessage(
				routingEntry,
				container.MessageCreate.Author,
				session.State.User,
				channel,
				container.MessageCreate.Content,
				container.Args,
				container.Prefix,
			) {
				continue
			}
		}
		if container.Type == MessageUpdateEventType {
			channel, err := session.State.Channel(container.MessageUpdate.ChannelID)
			if err != nil {
				continue
			}

			if !RoutingMatchMessage(
				routingEntry,
				container.MessageUpdate.Author,
				session.State.User,
				channel,
				container.MessageUpdate.Content,
				container.Args,
				container.Prefix,
			) {
				continue
			}
		}
		if container.Type == MessageDeleteEventType {
			channel, err := session.State.Channel(container.MessageDelete.ChannelID)
			if err != nil {
				continue
			}

			if !RoutingMatchMessage(
				routingEntry,
				container.MessageDelete.Author,
				session.State.User,
				channel,
				container.MessageDelete.Content,
				container.Args,
				container.Prefix,
			) {
				continue
			}
		}

		handled++
		aliases = append(aliases, routingEntry.Alias)
		if routingEntry.DestinationMain == "lambda" {
			lambdaDestinations = append(lambdaDestinations, routingEntry.DestinationSub)
		}
		if routingEntry.DestinationMain == "sqs" {
			processorDestinations = append(processorDestinations, routingEntry.DestinationSub)
		}
	}

	return
}
