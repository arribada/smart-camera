from picamera2 import Picamera2, Preview
from random import randint

import time
import datetime
import sys, getopt
import signal
import config_loader
import cv2

#used for closing
picam2=None

def keep_fps( start_time,final_time, fps):
    fps_real = 1 / (final_time - start_time)
    if fps_real < fps:
        pass
    elif fps_real >= fps:
        time.sleep((fps_real - fps) / (fps * fps_real))

#Ctrl C catch and close
def handler(signum, frame):
    print("closing")
    picam2.close()
    exit(1)

def main(argv):
    global picam2
    while True:

        try:
            opts, args = getopt.getopt(argv, "h", ["--help"])
        except getopt.GetoptError:
            sys.exit(2)

        config_loader.load_config(args[0])
        if config_loader.get_value("CAPTURE_PICAMERA_FPS") == 0:
            try:
                fps_needed = float(config_loader.get_value("CAPTURE_STREAM_FPS"))

                vcap = cv2.VideoCapture(config_loader.get_value("CAPTURE_STREAM_URL"), cv2.CAP_FFMPEG)
                while True:
                    start_time=time.time()
                    ret, frame = vcap.read()
                    if ret == False:
                        print("Frame is empty")
                        break
                    else:
                        now = datetime.datetime.now()
                        cv2.imwrite(config_loader.get_value("DATAFOLDER")+'/captured/'+str(now.hour)+str(now.minute)+str(now.second)+str(randint(0, 100))+'.jpg', frame)
                    keep_fps(start_time,time.time(),fps_needed)

                vcap.release()
            except Exception as ex:
                print(ex)

        else:
            try:
                fps_needed = float(config_loader.get_value("CAPTURE_PICAMERA_FPS"))

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

                    metadata = picam2.capture_file(config_loader.get_value("DATAFOLDER")+'/captured/'+str(now.hour)+str(now.minute)+str(now.second)+str(randint(0, 100))+'.jpg')
                    if config_loader.get_value("DEBUG") == 1:
                        metadata = picam2.capture_file(config_loader.get_value("DATAFOLDER")+'/debug/capture.jpg')
                    if config_loader.get_value("DEBUG") == 1:
                        print(metadata)

                    keep_fps(start_time,time.time(),fps_needed)
            except Exception as ex:
                print(ex)
        time.sleep(5)


if __name__ == "__main__":
   main(sys.argv[1:])