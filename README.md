# La Tavernière

Simple discord bot for moderation (and more).
<br />
Made in Go - my first project in this programming language.
<br /><br />
It was originally made for a friend, then edited for public release.
<br /><br />
Run with ```make```.

## Actions

- [x]  Announcement for newly posted youtube videos - NOT PUBLICLY OPEN TO EVERY GUILD FOR NOW
    - [ ]  Youtube Live Announcements
    - [ ]  Add reel in youtube announcements. See if /shorts/{idvideo} return something if shorts or error if not ou l'inverse
- [x]  Log message for each bot action
- [x]  Levels for conversation and messages

### Commands
>By default, the owner of a guild is an admin.

- [x]  Help command to describe all existing commands
- [x]  Config command to set channels and other thing for the bot to work properly
    - [ ]  Set youtube channels to track
- [x]  Message comand to send a message/an embed, with thumbnail, attachment... through the bot - admin only
- [x]  Blacklist command to ban a person and add their id and a reason (with date) to a chan - admin only
- [x]  Kick command to kick a personn - admin only
- [x]  Command to add or delete a handler to add a role to a user when a reaction is made on a specific message - admin only
    - [ ] Reaction are optionnal (all reaction taken in account if not specified)
- [x]  No youtube live command for 'today' or until a specified date - admin only - NOT PUBLICLY OPEN TO EVERY GUILD FOR NOW
- [x]  Level command to see someone's level or reset it - reset someone else level is admin only

## TODO
- [ ]  Vote command with limited time
- [ ]  Personalized join and leave message
- [ ]  Handler-reaction-role: choice to make the message
- [ ]  Stop bot properly
- [ ]  Set default channel for log
- [ ]  Warn user of an error by response or log message
- [ ]  Profile command (username, name on server, server name, account creation date, server join date, roles, about me, level, nb message, badges, pp, baner)

## Libraries
-  [discordgo](https://github.com/bwmarrin/discordgo), for a simple utilisation of the api discord
-  [gomoji](https://github.com/forPelevin/gomoji), for actions with emojis in string
-  [bun](https://github.com/uptrace/bun), to simplify interactions with postgres/sql database
-  [gg](https://github.com/fogleman/gg), for rendering 2D graphics
-  [freetype](https://github.com/golang/freetype), to load a font for graphics render (use of truetype)

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

## Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/GAsNA/la-taverniere/issues).

## License

Copyright © 2023 [GAsNa](https://github.com/GAsNa).<br />
This project is [GPL](https://github.com/GAsNA/la-taverniere/blob/main/LICENSE) licensed.

## Badges
![goBadge](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![discordBadge](https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white)
![youtubeBadge](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)
![dockerBadge](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)
