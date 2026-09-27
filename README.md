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

5. ~~Add your important birthdates to `birthdays.json`~~ BirthdayBot now runs on SQLite. Convert your existing data with the (non liable, as is) `convert.py` script, or use the `addBirthday` slash command to add to your database
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
9.  `systemctl enable birthday-bot.service`
10. `systemctl start birthday-bot.service`
11. `systemctl status birthday-bot.service`
12. BirthdayBot is set to run at 8am, and can be tested with any of the commands below 

## Available Commands
1. "/addbirthday first last {yyyy/}mm/dd {deceased}" adds a birthday entry 
2. "/getperson first last" retrieves a person record
3. "/getdate mm/dd" retrieves all people born on that date
4. "/getmonth mm" retrieves all people born in the month
5. "/getyear yyyy" retrieves all people born in that year 
6. "/update first last {yyyy/}mm/dd {deceased}" updates record of the person with name "first last"
7. "/schedule h" updates the scheduler to the 24h time specified 
8. "/remove first last" deletes the record of a person with name "first last"
9. "/help" will display the list of available commands and their arguments