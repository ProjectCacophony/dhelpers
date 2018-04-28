package state

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

func initGuildBans(session *discordgo.Session, guildID string) (err error) {
	// check if bot is allowed to see bans
	apermissions, err := UserPermissions(session.State.User.ID, guildID)
	if err != nil {
		return err
	}
	if apermissions&discordgo.PermissionBanMembers != discordgo.PermissionBanMembers {
		//fmt.Println("resetting bans for", guildID, "because no permissions")
		// reset ban list if not allowed
		deleteStateObject(guildBannedUserIDsSetKey(guildID))
		removeFromStateSet(guildBannedUserIDInitializedGuildIDsSetKey(), guildID)
		return nil
	}

	// have we already cached the guild bans for this guild?
	initializedGuildIDs, err := readStateSet(guildBannedUserIDInitializedGuildIDsSetKey())
	if err != nil {
		return err
	}

	var guildInitialized bool
	for _, initializedGuildID := range initializedGuildIDs {
		if initializedGuildID == guildID {
			guildInitialized = true
		}
	}

	if guildInitialized {
		//fmt.Println("ignoring initializing bans for", guildID, "because already initialized")
		return
	}

	// reset guild bans
	//fmt.Println("resetting bans for", guildID, "because caching new ones")
	deleteStateObject(guildBannedUserIDsSetKey(guildID))

	// cache new guild bans
	bans, err := session.GuildBans(guildID)
	if err != nil {
		return err
	}

	newSet := make([]string, 0)
	for _, ban := range bans {
		newSet = append(newSet, ban.User.ID)
	}
	err = addToStateSet(guildBannedUserIDsSetKey(guildID), newSet...)
	if err != nil {
		return err
	}
	//fmt.Println("setting bans for", guildID, ":", strings.Join(newSet, ", "))

	err = addToStateSet(guildBannedUserIDInitializedGuildIDsSetKey(), guildID)
	return err
}

func onReady(session *discordgo.Session, ready *discordgo.Ready) (err error) {
	//fmt.Println("running onReady")
	stateLock.Lock()
	defer stateLock.Unlock()

	// cache bot user
	err = updateStateObject(userKey(ready.User.ID), ready.User)
	if err != nil {
		return err
	}

	// cache guilds
	for _, guild := range ready.Guilds {
		err = updateStateObject(guildKey(guild.ID), guild)
		if err != nil {
			return err
		}
		err = addToStateSet(allGuildIDsSetKey(), guild.ID)
		if err != nil {
			return err
		}
		err = addToStateSet(guildBotIDsSetKey(guild.ID), ready.User.ID)
		if err != nil {
			return err
		}

		// cache guild channels
		for _, channel := range guild.Channels {
			err = updateStateObject(channelKey(channel.ID), channel)
			if err != nil {
				return err
			}
			err = addToStateSet(allChannelIDsSetKey(), channel.ID)
			if err != nil {
				return err
			}
		}

		// cache guild members and users
		for _, member := range guild.Members {
			err = updateStateObject(memberKey(member.GuildID, member.User.ID), member)
			if err != nil {
				return err
			}
			err = updateStateObject(userKey(member.User.ID), member.User)
			if err != nil {
				return err
			}
			err = addToStateSet(allUserIDsSetKey(), member.User.ID)
			if err != nil {
				return err
			}
			err = addToStateSet(guildUserIDsSetKey(member.GuildID), member.User.ID)
			if err != nil {
				return err
			}
		}

		// init guild bans (async)
		go func(gS *discordgo.Session, gGuildID string) {
			err = initGuildBans(gS, gGuildID)
			if err != nil {
				cache.GetLogger().WithField("module", "state").Errorln("error initializing bans for", gGuildID+":", err.Error())
			}
		}(session, guild.ID)
	}

	// cache private channels
	for _, privateChannel := range ready.PrivateChannels {
		err = updateStateObject(channelKey(privateChannel.ID), privateChannel)
		if err != nil {
			return err
		}
	}

	return nil
}

func guildAdd(session *discordgo.Session, guild *discordgo.Guild) (err error) {
	//fmt.Println("running guildAdd", guild.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// cache guild channels
	for _, channel := range guild.Channels {
		err = updateStateObject(channelKey(channel.ID), channel)
		if err != nil {
			return err
		}
		err = addToStateSet(allChannelIDsSetKey(), channel.ID)
		if err != nil {
			return err
		}
	}

	// cache guild members and users
	for _, member := range guild.Members {
		err = updateStateObject(memberKey(member.GuildID, member.User.ID), member)
		if err != nil {
			return err
		}
		err = updateStateObject(userKey(member.User.ID), member.User)
		if err != nil {
			return err
		}
		err = addToStateSet(allUserIDsSetKey(), member.User.ID)
		if err != nil {
			return err
		}
		err = addToStateSet(guildUserIDsSetKey(member.GuildID), member.User.ID)
		if err != nil {
			return err
		}
	}

	// carry over previous guild fields if set
	previousGuild, err := Guild(guild.ID)
	if err == nil {
		if guild.Roles == nil {
			guild.Roles = previousGuild.Roles
		}
		if guild.Emojis == nil {
			guild.Emojis = previousGuild.Emojis
		}
		if guild.Members == nil {
			guild.Members = previousGuild.Members
		}
		if guild.Presences == nil {
			guild.Presences = previousGuild.Presences
		}
		if guild.Channels == nil {
			guild.Channels = previousGuild.Channels
		}
		if guild.VoiceStates == nil {
			guild.VoiceStates = previousGuild.VoiceStates
		}
	}

	// cache guild
	err = updateStateObject(guildKey(guild.ID), guild)
	if err != nil {
		return err
	}
	err = addToStateSet(allGuildIDsSetKey(), guild.ID)
	if err != nil {
		return err
	}
	err = addToStateSet(guildBotIDsSetKey(guild.ID), session.State.User.ID)
	if err != nil {
		return err
	}

	// init guild bans (async)
	go func(gS *discordgo.Session, gGuildID string) {
		err = initGuildBans(gS, gGuildID)
		if err != nil {
			cache.GetLogger().WithField("module", "state").Errorln("error initializing bans for", gGuildID+":", err.Error())
		}
	}(session, guild.ID)

	return nil
}

func guildRemove(session *discordgo.Session, guild *discordgo.Guild) (err error) {
	//fmt.Println("running guildRemove", guild.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// remove guild
	err = deleteStateObject(guildKey(guild.ID))
	if err != nil {
		return err
	}
	err = removeFromStateSet(allGuildIDsSetKey(), guild.ID)
	if err != nil {
		return err
	}
	err = removeFromStateSet(guildBotIDsSetKey(guild.ID), session.State.User.ID)
	if err != nil {
		return err
	}

	// remove channels
	for _, channel := range guild.Channels {
		err = deleteStateObject(channelKey(channel.ID))
		if err != nil {
			return err
		}
		err = removeFromStateSet(allChannelIDsSetKey(), channel.ID)
		if err != nil {
			return err
		}
	}
	return err
}

func memberAdd(session *discordgo.Session, member *discordgo.Member) (err error) {
	//fmt.Println("running memberAdd", member.GuildID, member.User.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read member guild
	previousGuild, err := Guild(member.GuildID)
	if err != nil {
		return err
	}

	// read previous member if exists
	previousMember, err := Member(member.GuildID, member.User.ID)
	if err == nil {
		// carry over previous member fields if set
		if member.JoinedAt == "" {
			member.JoinedAt = previousMember.JoinedAt
		}
	} else {
		// update member guild
		previousGuild.Members = append(previousGuild.Members, member)
		previousGuild.MemberCount++
		// cache guild
		err = updateStateObject(guildKey(previousGuild.ID), previousGuild)
		if err != nil {
			return err
		}
	}

	// cache member and user
	err = updateStateObject(memberKey(member.GuildID, member.User.ID), member)
	if err != nil {
		return err
	}
	err = updateStateObject(userKey(member.User.ID), member.User)
	if err != nil {
		return err
	}
	err = addToStateSet(allUserIDsSetKey(), member.User.ID)
	if err != nil {
		return err
	}
	err = addToStateSet(guildUserIDsSetKey(member.GuildID), member.User.ID)
	if err != nil {
		return err
	}

	if member.User.ID == session.State.User.ID {
		// init guild bans (async) (could be giving or revoking the bot ban permission)
		go func(gS *discordgo.Session, gGuildID string) {
			err = initGuildBans(gS, gGuildID)
			if err != nil {
				cache.GetLogger().WithField("module", "state").Errorln("error initializing bans for", gGuildID+":", err.Error())
			}
		}(session, member.GuildID)
	}

	return nil
}

func memberRemove(member *discordgo.Member) (err error) {
	//fmt.Println("running memberRemove", member.GuildID, member.User.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read member guild
	previousGuild, err := Guild(member.GuildID)
	if err != nil {
		return err
	}

	// remove member and user
	err = deleteStateObject(memberKey(member.GuildID, member.User.ID))
	if err != nil {
		return err
	}

	// viable?
	allGuildIDs, err := AllGuildIDs()
	if err == nil {
		var isMember bool
		for _, guildID := range allGuildIDs {
			if guildID == member.GuildID {
				continue
			}
			isMember, err = IsMember(guildID, member.User.ID)
			if err != nil {
				return err
			}
			if isMember {
				break
			}
		}
		if !isMember {
			err = deleteStateObject(userKey(member.User.ID))
			if err != nil {
				return err
			}
			err = removeFromStateSet(allUserIDsSetKey(), member.User.ID)
			if err != nil {
				return err
			}
		}
	}

	// update previous guild
	for i, previousMember := range previousGuild.Members {
		if previousMember.User.ID == member.User.ID {
			previousGuild.Members = append(previousGuild.Members[:i], previousGuild.Members[i+1:]...)
			previousGuild.MemberCount--
			break
		}
	}
	err = removeFromStateSet(guildUserIDsSetKey(member.GuildID), member.User.ID)
	if err != nil {
		return err
	}

	// cache guild
	err = updateStateObject(guildKey(previousGuild.ID), previousGuild)
	return err
}

func roleAdd(session *discordgo.Session, guildID string, role *discordgo.Role) (err error) {
	//fmt.Println("running roleAdd", guildID, role.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read role guild
	previousGuild, err := Guild(guildID)
	if err != nil {
		return err
	}

	// update previous guild
	var updated bool
	for i, previousRole := range previousGuild.Roles {
		if previousRole.ID == role.ID {
			previousGuild.Roles[i] = role
			updated = true
			break
		}
	}
	if !updated {
		previousGuild.Roles = append(previousGuild.Roles, role)
	}

	// cache guild
	err = updateStateObject(guildKey(previousGuild.ID), previousGuild)
	if err != nil {
		return err
	}

	if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator ||
		role.Permissions&discordgo.PermissionBanMembers == discordgo.PermissionBanMembers {
		// init guild bans (async) (could be giving or revoking the bot ban permission)
		go func(gS *discordgo.Session, gGuildID string) {
			err = initGuildBans(gS, gGuildID)
			if err != nil {
				cache.GetLogger().WithField("module", "state").Errorln("error initializing bans for", gGuildID+":", err.Error())
			}
		}(session, guildID)
	}

	return nil
}

func roleRemove(guildID, roleID string) (err error) {
	//fmt.Println("running roleRemove", guildID, roleID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read role guild
	previousGuild, err := Guild(guildID)
	if err != nil {
		return err
	}

	// remove role
	for i, previousRole := range previousGuild.Roles {
		if previousRole.ID == roleID {
			previousGuild.Roles = append(previousGuild.Roles[:i], previousGuild.Roles[i+1:]...)
			return nil
		}
	}

	// cache guild
	err = updateStateObject(guildKey(previousGuild.ID), previousGuild)
	return err
}

func emojiAdd(guildID string, emoji *discordgo.Emoji) (err error) {
	//fmt.Println("running emojiAdd", guildID, emoji.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read emoji guild
	previousGuild, err := Guild(guildID)
	if err != nil {
		return err
	}

	// update previous guild
	var updated bool
	for i, previousEmoji := range previousGuild.Emojis {
		if previousEmoji.ID == emoji.ID {
			previousGuild.Emojis[i] = emoji
			updated = true
			break
		}
	}
	if !updated {
		previousGuild.Emojis = append(previousGuild.Emojis, emoji)
	}

	// cache guild
	err = updateStateObject(guildKey(previousGuild.ID), previousGuild)
	return err
}

func emojisAdd(guildID string, emojis []*discordgo.Emoji) (err error) {
	//fmt.Println("running emojisAdd", guildID, len(emojis))
	for _, emoji := range emojis {
		err = emojiAdd(guildID, emoji)
		if err != nil {
			return err
		}
	}
	return nil
}

func channelAdd(channel *discordgo.Channel) (err error) {
	//fmt.Println("running channelAdd", channel.GuildID, channel.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read channel
	previousChannel, err := Channel(channel.ID)
	if err == nil {
		// carry over previous fields if set
		if channel.Messages == nil {
			channel.Messages = previousChannel.Messages
		}
		if channel.PermissionOverwrites == nil {
			channel.PermissionOverwrites = previousChannel.PermissionOverwrites
		}

		// cache channel
		err = updateStateObject(channelKey(channel.ID), channel)
		return err
	}

	if channel.Type != discordgo.ChannelTypeDM && channel.Type != discordgo.ChannelTypeGroupDM {
		// read channel guild
		var previousGuild *discordgo.Guild
		previousGuild, err = Guild(channel.GuildID)
		if err != nil {
			return err
		}

		// update guild
		previousGuild.Channels = append(previousGuild.Channels, channel)

		// cache guild
		err = updateStateObject(guildKey(previousGuild.ID), previousGuild)
		if err != nil {
			return err
		}
	}

	// cache channel
	err = updateStateObject(channelKey(channel.ID), channel)
	if err != nil {
		return err
	}
	err = addToStateSet(allChannelIDsSetKey(), channel.ID)
	return err
}

func channelRemove(channel *discordgo.Channel) (err error) {
	//fmt.Println("running channelRemove", channel.GuildID, channel.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read channel
	previousChannel, err := Channel(channel.ID)
	if err != nil {
		return err
	}

	if channel.Type != discordgo.ChannelTypeDM && channel.Type != discordgo.ChannelTypeGroupDM {
		// read channel guild
		var previousGuild *discordgo.Guild
		previousGuild, err = Guild(previousChannel.GuildID)
		if err != nil {
			return err
		}

		// update guild
		for i, previousGuildChannel := range previousGuild.Channels {
			if previousGuildChannel.ID == channel.ID {
				previousGuild.Channels = append(previousGuild.Channels[:i], previousGuild.Channels[i+1:]...)
				break
			}
		}

		// cache guild
		err = updateStateObject(guildKey(previousGuild.ID), previousGuild)
		if err != nil {
			return err
		}
	}

	// cache channel
	err = deleteStateObject(channelKey(channel.ID))
	if err != nil {
		return err
	}
	err = removeFromStateSet(allChannelIDsSetKey(), channel.ID)
	return err
}

func presenceAdd(guildID string, presence *discordgo.Presence) (err error) {
	//fmt.Println("running presenceAdd", guildID, presence.User.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read presence guild
	previousGuild, err := Guild(guildID)
	if err != nil {
		return err
	}

	// update presence
	var updated bool
	for i, previousPresence := range previousGuild.Presences {
		if previousPresence.User.ID == presence.User.ID {
			//Update status
			previousGuild.Presences[i].Game = presence.Game
			previousGuild.Presences[i].Roles = presence.Roles
			if presence.Status != "" {
				previousGuild.Presences[i].Status = presence.Status
			}
			if presence.Nick != "" {
				previousGuild.Presences[i].Nick = presence.Nick
			}

			//Update the optionally sent user information
			//ID Is a mandatory field so you should not need to check if it is empty
			previousGuild.Presences[i].User.ID = presence.User.ID

			if presence.User.Avatar != "" {
				previousGuild.Presences[i].User.Avatar = presence.User.Avatar
			}
			if presence.User.Discriminator != "" {
				previousGuild.Presences[i].User.Discriminator = presence.User.Discriminator
			}
			if presence.User.Email != "" {
				previousGuild.Presences[i].User.Email = presence.User.Email
			}
			if presence.User.Token != "" {
				previousGuild.Presences[i].User.Token = presence.User.Token
			}
			if presence.User.Username != "" {
				previousGuild.Presences[i].User.Username = presence.User.Username
			}

			updated = true
		}
	}
	if !updated {
		previousGuild.Presences = append(previousGuild.Presences, presence)
	}

	// cache guild
	err = updateStateObject(guildKey(previousGuild.ID), previousGuild)
	return err
}

func banAdd(session *discordgo.Session, guildID string, user *discordgo.User) (err error) {
	// check if bot is allowed to see bans
	apermissions, err := UserPermissions(session.State.User.ID, guildID)
	if err != nil {
		return err
	}
	if apermissions&discordgo.PermissionBanMembers != discordgo.PermissionBanMembers {
		return nil
	}

	// add ban
	err = addToStateSet(guildBannedUserIDsSetKey(guildID), user.ID)
	return err
}

func banRemove(guildID string, user *discordgo.User) (err error) {
	err = removeFromStateSet(guildBannedUserIDsSetKey(guildID), user.ID)
	return err
}

// SharedStateEventHandler receives events from a discordgo Websocket and updates the shared state with them
func SharedStateEventHandler(session *discordgo.Session, i interface{}) error {
	ready, ok := i.(*discordgo.Ready)
	if ok {
		return onReady(session, ready)
	}

	switch t := i.(type) {
	case *discordgo.GuildCreate:
		return guildAdd(session, t.Guild)
	case *discordgo.GuildUpdate:
		return guildAdd(session, t.Guild)
	case *discordgo.GuildDelete:
		return guildRemove(session, t.Guild)
	case *discordgo.GuildMemberAdd:
		return memberAdd(session, t.Member)
	case *discordgo.GuildMemberUpdate:
		return memberAdd(session, t.Member)
	case *discordgo.GuildMemberRemove:
		return memberRemove(t.Member)
	case *discordgo.GuildMembersChunk:
		for i := range t.Members {
			t.Members[i].GuildID = t.GuildID
			err := memberAdd(session, t.Members[i])
			if err != nil {
				return err
			}
		}
		return nil
	case *discordgo.GuildRoleCreate:
		return roleAdd(session, t.GuildID, t.Role)
	case *discordgo.GuildRoleUpdate:
		return roleAdd(session, t.GuildID, t.Role)
	case *discordgo.GuildRoleDelete:
		return roleRemove(t.GuildID, t.RoleID)
	case *discordgo.GuildEmojisUpdate:
		return emojisAdd(t.GuildID, t.Emojis)
	case *discordgo.ChannelCreate:
		return channelAdd(t.Channel)
	case *discordgo.ChannelUpdate:
		return channelAdd(t.Channel)
	case *discordgo.ChannelDelete:
		return channelRemove(t.Channel)
	case *discordgo.GuildBanAdd:
		return banAdd(session, t.GuildID, t.User)
	case *discordgo.GuildBanRemove:
		return banRemove(t.GuildID, t.User)
	case *discordgo.PresenceUpdate:
		err := presenceAdd(t.GuildID, &t.Presence)
		if err != nil {
			return err
		}

		previousMember, err := Member(t.GuildID, t.User.ID)
		if err != nil {
			// Member not found; this is a user coming online
			previousMember = &discordgo.Member{
				GuildID: t.GuildID,
				Nick:    t.Nick,
				User:    t.User,
				Roles:   t.Roles,
			}

		} else {
			if t.Nick != "" {
				previousMember.Nick = t.Nick
			}

			if t.User.Username != "" {
				previousMember.User.Username = t.User.Username
			}

			if t.User.Discriminator != "" {
				previousMember.User.Discriminator = t.User.Discriminator
			}

			if t.User.Avatar != "" {
				previousMember.User.Avatar = t.User.Avatar
			}

			// PresenceUpdates always contain a list of roles, so there's no need to check for an empty list here
			previousMember.Roles = t.Roles

		}

		return memberAdd(session, previousMember)
		/*
		   case *discordgo.MessageCreate:
		       if s.MaxMessageCount != 0 {
		           err = s.MessageAdd(t.Message)
		       }
		   case *discordgo.MessageUpdate:
		       if s.MaxMessageCount != 0 {
		           err = s.MessageAdd(t.Message)
		       }
		   case *discordgo.MessageDelete:
		       if s.MaxMessageCount != 0 {
		           err = s.MessageRemove(t.Message)
		       }
		   case *discordgo.MessageDeleteBulk:
		       if s.MaxMessageCount != 0 {
		           for _, mID := range t.Messages {
		               s.messageRemoveByID(t.ChannelID, mID)
		           }
		       }
		   case *discordgo.VoiceStateUpdate:
		       if s.TrackVoice {
		           err = s.voiceStateUpdate(t)
		       }
		*/

	}

	return nil
}
