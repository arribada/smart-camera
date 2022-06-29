import os
import time

import config_loader


def main(argv):
    config_loader.load_config(argv[0])
    while True:
        try:
            os.system("edge-impulse-uploader "+config_loader.get_value("DATAFOLDER") +"/detected/*.jpg")
        except Exception as ex:
            print(ex)

        import glob
        removing_files = glob.glob(config_loader.get_value("DATAFOLDER") +"/detected/*.jpg")
        for i in removing_files:
            os.remove(i)
        print("done")
        time.sleep(100)