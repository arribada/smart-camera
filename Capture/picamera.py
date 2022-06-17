from picamera2 import Picamera2, Preview
import time

picam2 = Picamera2()

preview_config = picam2.preview_configuration(main={"size": (800, 600),"ExposureTime": 1000})


picam2.configure(preview_config)


#picam2.start_preview(Preview.QTGL)

picam2.start()
time.sleep(2)
import datetime

while True:
    now = datetime.datetime.now()
    metadata = picam2.capture_file('./pictures/'+str(now.hour)+str(now.minute)+str(now.second)+'.jpg')
    time.sleep(3)
    print(metadata)

picam2.close()