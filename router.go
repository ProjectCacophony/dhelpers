package dhelpers

import (
	"regexp"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
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

	Events        []EventType
	Module        string
	Destination   string
	Requirement   []rawRoutingRequirementEntry // will only get matched with EventTypeMessageCreate, EventTypeMessageUpdate, or EventTypeMessageDelete, will match everything if slice is empty
	Priority      int                          // higher runs before lower
	ErrorHandlers []string
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
	ErrorHandlers      []ErrorHandlerType
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
					DestinationMain: replaceDestinations(parts[0]),
					DestinationSub:  replaceDestinations(parts[1]),
					Always:          rawRule.Always,
					AllowMyself:     rawRule.AllowMyself,
					AllowBots:       rawRule.AllowBots,
					AllowDM:         rawRule.AllowDM,

					Beginning:          "",
					Regex:              nil,
					DoNotPrependPrefix: false,
					CaseSensitive:      false,
				}
				for _, errorHandler := range rawRule.ErrorHandlers {
					switch errorHandler {
					case string(SentryErrorHandler):
						newEntry.ErrorHandlers = append(newEntry.ErrorHandlers, SentryErrorHandler)
					case string(DiscordErrorHandler):
						newEntry.ErrorHandlers = append(newEntry.ErrorHandlers, DiscordErrorHandler)
					}
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
	if routingEntry.Beginning != "" || routingEntry.Regex != nil {
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
func ContainerDestinations(session *discordgo.Session, routingConfig []RoutingRule, container EventContainer) (destinations []DestinationData) {
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

		switch routingEntry.DestinationMain {
		case "kafka":
			destinations = append(destinations, DestinationData{
				Type:          KafkaDestinationType,
				Name:          routingEntry.DestinationSub,
				ErrorHandlers: routingEntry.ErrorHandlers,
				Alias:         routingEntry.Alias,
			})
		}
	}

	return
}

// replaceDestinations reples various placeholders
// currently supported {ENV} => current environment
func replaceDestinations(input string) string {
	return strings.Replace(input, "{ENV}", string(GetEnvironment()), -1)
}
