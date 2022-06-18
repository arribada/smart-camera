import os
import time

while True:
    try:
        os.system("edge-impulse-uploader ./detected/*.jpg")
    except Exception as ex:
        print(ex)

    import glob
    removing_files = glob.glob('/root/edge/linux-sdk-python/working/detected/*.jpg')
    for i in removing_files:
        os.remove(i)
    print("done")
    time.sleep(100)