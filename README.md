# La Tavernière

Simple discord bot for moderation (and more).
<br />
Made in Go - my first project in this programming language.
<br /><br />
It was originally made for a friend, then edited for public release.
<br /><br />
Run with ```make```.

## TODO
- [ ]  Multi-guild accessibility
- [ ]  Verification if text channel, when channel is required in command
        ChannelTypes:[]discordgo.ChannelType {
            discordgo.ChannelTypeGuildText,
        },
- [ ]  Description message command in readme
- [ ]  Verification message command when no title and no embed...
- [ ]  Verification all log error
- [ ]  Log message for message command: supp "in channel"

## Actions

- [x]  Announcement for newly posted youtube videos and youtube live (for one youtube channel for now)
    - [ ]  Add reel in youtube announcements. See if /shorts/{idvideo} return something if shorts or error if not ou l'inverse
    - [ ]  Stop bot properly
    - [ ]  Multi-guild accessibility
- [x]  Log message for each bot action
- [x]  Levels for conversation and messages

### Commands
>By default, the owner of a guild is an admin.

- [x]  No youtube live command for 'today' or until a specified date - admin only
- [x]  Blacklist command to ban a person and add their id and a reason (with date) to a chan - admin only
- [x]  Kick command to kick a personn - admin only
- [x]  Command to add or delete a handler to add a role to a user when a reaction is made on a specific message - admin only
    - [ ] Merge add and delete commands
    - [ ] Reaction are optionnal (all reaction taken in account if not specified)
- [x]  Commands for Ray, Feitan, Ukyim, Kentaro, GAsNa
- [x]  Salope command for Kentaro
- [ ]  Command to add role with nam e: /role us -> add role to the person who ask. If doesn't exist, add argument for role you want
- [ ]  Help command to describe all existing commands
- [ ]  Config command to set channels and other thing for the bot to work properly

## Libraries
-  [discordgo](https://github.com/bwmarrin/discordgo), for a simple utilisation of the api discord in Go
-  [gomoji](https://github.com/forPelevin/gomoji), for actions with emojis in string in Go
-  [bun](https://github.com/uptrace/bun), to simplify interactions with postgres/sql database

## Author
<table>
  <tr>
    <td align="center">
      <a href="https://github.com/GAsNA">
        <img src="https://avatars.githubusercontent.com/u/58465901?v=4" width="100px;" alt=""/>
      </a>
      <br />
      <sub>
        <a href="https://github.com/GAsNA">
          <b>@GAsNa</b>
        </a>
        <br />
      </sub>
    </td>
  </tr>
</table>

## Badges
![goBadge](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![discordBadge](https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white)
![youtubeBadge](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)
![dockerBadge](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)
