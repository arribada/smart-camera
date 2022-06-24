import cv2
import numpy as np
import datetime
vcap = cv2.VideoCapture("rtsp://admin:1qazxsw2@95.87.206.5", cv2.CAP_FFMPEG)

while(1):
    ret, frame = vcap.read()
    if ret == False:
        print("Frame is empty")
        break
    else:
        now = datetime.datetime.now()
        cv2.imwrite('./pictures/'+str(now.hour)+str(now.minute)+str(now.second)+'.jpg', frame)
        cv2.waitKey(1)
    vcap.release()