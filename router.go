package dhelpers

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

type EventType string

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
type rawRoutingEntry struct {
	Active       bool
	Type         []EventType
	Endpoint     string
	Requirements []rawRoutingRequirementEntry // will only get matched with EventTypeMessageCreate, EventTypeMessageUpdate, or EventTypeMessageDelete, will match everything if slice is empty
	Always       bool                         // if true: will run even if there have been previous (higher priority) matches
	Priority     int                          // higher runs before lower
}
type rawRoutingRequirementEntry struct {
	Beginning          string // can be empty, will match all
	Regex              string // can be empty, will match all
	DoNotPrependPrefix bool   // if false, prepends guild prefix to regex
	CaseSensitive      bool   // prepends (?i) to regex on go, language dependent#
	Alias              string
}

// Routing Compiled Config
type RoutingRule struct {
	Type               EventType
	Endpoint           string
	Beginning          string
	Regex              *regexp.Regexp
	DoNotPrependPrefix bool
	CaseSensitive      bool
	Always             bool
	Alias              string
}

// returns a sorted slice (by priority) with all rules
func GetRoutings() (routingRules []RoutingRule, err error) {
	// read and unmarshal config from file
	// TODO: load from S3 instead
	var rawRouting []rawRoutingEntry
	routingFileData, err := ioutil.ReadFile("routing.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(routingFileData, &rawRouting)
	if err != nil {
		return nil, err
	}

	// group rules by priorities
	rawEntriesByPriority := make(map[int][]rawRoutingEntry, 0)

	for _, rawRoutingEntry := range rawRouting {
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
			if rawRule.Type == nil || len(rawRule.Type) < 0 {
				continue
			}
			// skip empty endpoints
			if rawRule.Endpoint == "" {
				continue
			}
			// generate route for each type
			for _, ruleType := range rawRule.Type {
				newEntry := RoutingRule{
					Type:     ruleType,
					Endpoint: rawRule.Endpoint,
					Always:   rawRule.Always,

					Beginning:          "",
					Regex:              nil,
					DoNotPrependPrefix: false,
					CaseSensitive:      false,
				}
				if (ruleType == MessageCreateEventType ||
					ruleType == MessageUpdateEventType ||
					ruleType == MessageDeleteEventType) &&
					rawRule.Requirements != nil && len(rawRule.Requirements) > 0 {
					for _, requirement := range rawRule.Requirements {
						newEntryCopy := newEntry
						newEntryCopy.Beginning = requirement.Beginning
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
					routingRules = append(routingRules, newEntry)
				}
			}
		}
	}

	return routingRules, nil
}

// checks if a message content matches the requirements of the routing rule
func RoutingMatchMessage(routingEntry RoutingRule, content, prefix string) (match bool) {
	// match beginning if beginning is set
	if routingEntry.Beginning != "" {
		if routingEntry.CaseSensitive {
			if routingEntry.DoNotPrependPrefix {
				if !strings.HasPrefix(content, routingEntry.Beginning) {
					return false
				}
			} else {
				if !strings.HasPrefix(content, prefix+routingEntry.Beginning) {
					return false
				}
			}
		} else {
			if routingEntry.DoNotPrependPrefix {
				if !strings.HasPrefix(strings.ToLower(content), strings.ToLower(routingEntry.Beginning)) {
					return false
				}
			} else {
				if !strings.HasPrefix(strings.ToLower(content), prefix+strings.ToLower(routingEntry.Beginning)) {
					return false
				}
			}
		}
	}
	// match regex if regex is set
	if routingEntry.Regex != nil {
		if routingEntry.DoNotPrependPrefix {
			if !routingEntry.Regex.MatchString(content) {
				return false
			}
		} else {
			if !strings.HasPrefix(content, prefix) {
				return false
			}
			matchContent := strings.TrimLeft(content, prefix)
			if !routingEntry.Regex.MatchString(matchContent) {
				return false
			}
		}
	}

	return true
}
