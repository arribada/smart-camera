import config_loader
import os
import time
import sys
from telegram import Bot



def main(argv):
    config_loader.load_config(argv[0])
    TOKEN=config_loader.get_value("ALERT_TOKEN")
    channel_id=config_loader.get_value("ALERT_CHANNELID")
    bot = Bot(TOKEN)

    while True:
        try:
            dir_name = config_loader.get_value("DATAFOLDER")+"/detectedTelegram/"
            # Get list of all files only in the given directory
            list_of_files = filter( lambda x: os.path.isfile(os.path.join(dir_name, x)),
                                    os.listdir(dir_name) )
            # Sort list of files based on last modification time in ascending order
            list_of_files = sorted( list_of_files,
                                    key = lambda x: os.path.getmtime(os.path.join(dir_name, x)))
            # Iterate over sorted list of files and print file path
            # along with last modification time of file

            for file in list_of_files:
                try:
                    bot.send_photo(channel_id,photo=open(os.path.join(dir_name, file),"rb"))
                except Exception as ex:
                    print(ex)
                else:
                    time.sleep(0.5)
                    os.remove(os.path.join(dir_name, file))
            time.sleep(10)
        except Exception as ex:
            print(ex)

if __name__ == "__main__":
   main(sys.argv[1:])