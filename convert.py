import sqlite3
import json

# Read old format data file
input = open('data/birthdays.json.example').read()
data = json.loads(input)

insert = 'INSERT INTO birthdays(first_name, last_name, month, day) VALUES(?,?,?,?)'

# Open db connection
con = sqlite3.connect('data/birthdays.db')
cursor = con.cursor()

# Cascade through dictionary key-value pairs
for month, date in data.items():
    for date, people in date.items():
        for person in people:
            # Check that record isn't blank
            if person.strip():
                first = person.split(' ')[0]
                last = person.split(' ')[1]
                # Make new record
                info = (first, last, month, date)
                print(info)
                cursor.execute(insert, info)
                con.commit()

con.close()