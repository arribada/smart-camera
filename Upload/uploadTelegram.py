import telegram.ext
from telegram import Bot
import os
with open('token.txt','r') as f:
    TOKEN = f.read()



channel_id="-790283961"
bot = Bot(TOKEN)
#bot.send_message(channel_id, "test")

#bot.send_photo(channel_id,photo=open("C:\\Users\\Practica 2019\\Desktop\\Capture.PNG","rb"))


import os
import time

while True:
    dir_name = "/root/edge/linux-sdk-python/working/detectedTelegram/"
    # Get list of all files only in the given directory
    list_of_files = filter( lambda x: os.path.isfile(os.path.join(dir_name, x)),
                            os.listdir(dir_name) )
    # Sort list of files based on last modification time in ascending order
    list_of_files = sorted( list_of_files,
                            key = lambda x: os.path.getmtime(os.path.join(dir_name, x)))
    # Iterate over sorted list of files and print file path
    # along with last modification time of file

    for file in list_of_files:
        bot.send_photo(channel_id,photo=open(os.path.join(dir_name, file),"rb"))
        time.sleep(3)
        os.remove(os.path.join(dir_name, file))
    time.sleep(10)