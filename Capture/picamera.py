from picamera2 import Picamera2, Preview
import time
import datetime
import sys, getopt
import signal

import ConfigLoader
from ConfigLoader import *

#used for closing
picam2=None

#Ctrl C catch and close
def handler(signum, frame):
    print("closing")
    picam2.close()
    exit(1)

def main(argv):


    global picam2
    try:
        opts, args = getopt.getopt(argv, "h", ["--help"])
    except getopt.GetoptError:
        sys.exit(2)

    ConfigLoader.load_config(argv[0])
    if ConfigLoader.get_value("CAPTURE_PICAMERA_FPS") ==0:
        pass
    else:

        sleep = float(ConfigLoader.get_value("CAPTURE_PICAMERA_FPS"))

        #initiate Picamera and set configure size photo
        picam2 = Picamera2()
        preview_config = picam2.preview_configuration(main={"size": (800, 600),"ExposureTime": 1000})
        picam2.configure(preview_config)
        picam2.start()

        time.sleep(2)

        #catch Ctrl C
        signal.signal(signal.SIGINT, handler)

        #loop and save captured image every x seconds
        while True:
            start_time = time.time()

            now = datetime.datetime.now()
            metadata = picam2.capture_file('./pictures/'+str(now.hour)+str(now.minute)+str(now.second)+'.jpg')
            #metadata = picam2.capture_file('capture.jpg')
            time.sleep(sleep)
            #print(metadata)
            fps=1.0 / (time.time() - start_time)
            if fps>ConfigLoader.get_value("CAPTURE_PICAMERA_FPS"):
                time.sleep()

    else:


if __name__ == "__main__":
   main(sys.argv[1:])
   