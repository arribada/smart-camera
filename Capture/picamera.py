from picamera2 import Picamera2, Preview
import time
import datetime
import sys, getopt
import signal

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
    sleep = float(args[0])

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
        now = datetime.datetime.now()
        metadata = picam2.capture_file('./pictures/'+str(now.hour)+str(now.minute)+str(now.second)+'.jpg')
        #metadata = picam2.capture_file('capture.jpg')
        time.sleep(sleep)
        #print(metadata)


if __name__ == "__main__":
   main(sys.argv[1:])
   