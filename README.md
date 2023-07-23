# La Tavernière

Simple discord bot for moderation (and more).
<br />
Made in Go - my first project in this programming language.
<br /><br />
Run with ```go run src/*```.

## Actions
For commands, the bot uses interactions (slash commands).

- [x]  Announcement for newly posted youtube videos and youtube live (for one youtube channel for now)
    - [ ]  Add reel in youtube announcements
    - [ ]  Stop bot properly
- [x]  No youtube live command for 'today' or until a specified date - admin only
- [x]  Blacklist command to ban a person and add their id and a reason (with date) to a chan - admin only
- [x]  Kick command to kick a personn - admin only
- [x]  Command to add a handler to add a role to a user when a reaction is made on a specific message - admin only
- [x]  Log message for each bot action

### Optionnels
- [x]  Commands for Ray, Feitan, Ukyim, Kentaro, GAsNa
- [x]  Salope command for Kentaro
- [ ]  Levels for conversation and messages

## Libraries
-  [discordgo](https://github.com/bwmarrin/discordgo), for a simple utilisation of the api discord in Go
-  [godotenv](https://github.com/joho/godotenv), to load env vars from a .env file in Go
-  [gomoji](https://github.com/forPelevin/gomoji), for actions with emojis in string in Go

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
