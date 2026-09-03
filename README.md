# BirthdayBot

Discord bot that sends a reminder so you never forgot an important birthday again!

Set up tutorial by [Brian Morrison](https://youtu.be/XuFq7NW3ii4?si=g9wnxcFpjzr60P8O)
Database tutorial by [@nawazdhandala](https://oneuptime.com/blog/post/2026-02-02-sqlite-go/view)

## Set Up Instructions (Raspberry Pi)

1. Make sure Raspberry Pi is set up with Raspberry Pi OS
2. Follow instructions from [GoLang](https://go.dev/doc/install) to install Go. Pick an ARM version, such as [https://go.dev/dl/go1.26.5.linux-arm64.tar.gz](https://go.dev/dl/go1.26.5.linux-arm64.tar.gz)
3. Clone this repo
4. Add environment variables to `.env`
   1. Add your bot token, generated from [Discord Developer Portal](https://discord.com/developers/applications)
   2. Add your channel ID, found in the Discord app with developer settings turned on
5. Add your importants birthdates to `birthdays.json`
6. In the BirthdayBot folder, run `go build -o birthday-bot.exe`
7. `nano /etc/systemd/system/birthday-bot.service`
    ```
    [Unit]
    Description=Birthday Reminder Discord Bot
    After=network-online.target

    [Service]
    Type=simple
    WorkingDirectory=/home/{username}/BirthdayBot
    ExecStart=/home/{username}/BirthdayBot/birthday-bot.exe
    Restart=always
    RestartSec=5
    User={username}

    [Install]
    WantedBy=multi-user.target
    ```
8. `systemctl reload birthday-bot.service`
9. `systemctl enable birthday-bot.service`
10. `systemctl start birthday-bot.service`
11. `systemctl status birthday-bot.service`
12. BirthdayBot is set to run at 8:30am, or can be tested by messaging "Birthdays?" in any server channel it's invited to. Slash commands to come. 

## Available Commands
1. "Birthdays?" returns any birthdays for todays date
2. [Coming soon] "/addbirthday first last yyyy/mm/dd" adds a birthday entry to `birthdays.json`