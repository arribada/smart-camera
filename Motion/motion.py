import cv2
import datetime
import sys, getopt
import os
import time
import config_loader
from random import randint

class Rois():
    def __init__(self):
        self.listRois=list()

    def choose_ROIs(self,frame):
        self.frame=frame
        self.save=frame
        self.listRois=[(config_loader.get_value("MOTION_ROI")[0],
                        config_loader.get_value("MOTION_ROI")[1],
                         config_loader.get_value("MOTION_ROI")[2],
                          config_loader.get_value("MOTION_ROI")[3])]
     
        
    def check_if_overlapping(self, x, y, w, h):
        for roi in self.listRois:
            return not (roi[0][0]+roi[1][0] < x or roi[0][0] > x+w or roi[0][1] < y+h or roi[0][1]+roi[1][1] > y)

    def overlap(self,frame1, x, y, w, h):
        for roi in self.listRois:
            #sunt inversate, e x1 si y1 jos stanga
            newx1=x
            newx2=x+w
            newy2=y
            newy1=y+h
            x1=roi[0]
            x2=roi[2]
            y2=roi[1]
            y1=roi[3]

            if config_loader.get_value("DEBUG") == 1:
                cv2.rectangle(frame1, (x1, y1), (x2, y2), (0, 255, 0), 2)
            if (x1 < newx2 and x2 > newx1 and
                    y1 > newy2 and y2 < newy1):
                return True
        return False



def get_image(dir_name,file_name):
    attempts = 0
    while attempts < 10:
        try:
            image=cv2.imread(os.path.join(dir_name, file_name), cv2.IMREAD_COLOR)
            return image
        except Exception as ex:
            print(ex)
        attempts += 1
        time.sleep(0.2)
    return None


def main(argv):
    config_loader.load_config(argv[0])

    while True:
        try:
            opts, args = getopt.getopt(argv, "h", ["--help"])
        except getopt.GetoptError:
            sys.exit(2)
        try:
            dir_name = config_loader.get_value("DATAFOLDER")+'/captured/'
            # Get list of all files only in the given directory
            list_of_files = filter( lambda x: os.path.isfile(os.path.join(dir_name, x)),
                                    os.listdir(dir_name) )
            # Sort list of files based on last modification time in ascending order
            list_of_files = sorted( list_of_files,
                                    key = lambda x: os.path.getmtime(os.path.join(dir_name, x)))
            # Iterate over sorted list of files and print file path
            # along with last modification time of file
    
            file_path = os.path.join(dir_name, list_of_files[0])
            frame1 = get_image(dir_name,file_path)
            RoisClass=Rois()
            RoisClass.choose_ROIs(frame1)
            frame2 = get_image(dir_name,file_path)
    
            file_name_to_del=''
    
            while True:
                time.sleep(1)
                list_of_files = filter(lambda x: os.path.isfile(os.path.join(dir_name, x)),
                                       os.listdir(dir_name))
                # Sort list of files based on last modification time in ascending order
                list_of_files = sorted(list_of_files,
                                       key=lambda x: os.path.getmtime(os.path.join(dir_name, x)))
                # Iterate over sorted list of files and print file path
                # along with last modification time of file
    
                for file_name in list_of_files:
                
                        #find diff between frames
                        diff = cv2.absdiff(frame1,frame2)
    
                        gray = cv2.cvtColor(diff, cv2.COLOR_BGR2GRAY)
                        blur=cv2.GaussianBlur(gray,(5,5),0)
    
                        #find threshold
                        _,thresh = cv2.threshold(blur,20, 255, cv2.THRESH_BINARY)
                        #cv2.imshow("feed3", thresh)
                        dilated = cv2.dilate(thresh,None, iterations=3)
    
                        #cv2.imshow("feed2", dilated)
                        contours,_ = cv2.findContours(dilated,cv2.RETR_TREE,cv2.CHAIN_APPROX_SIMPLE)
    
    
                        for contour in contours:
                            (x,y,w,h) = cv2.boundingRect(contour)
    
                            if cv2.contourArea(contour) < 700:
                                continue
                            if RoisClass.overlap(frame1,x,y,w,h):
                                now = datetime.datetime.now()
                                cv2.imwrite(config_loader.get_value("DATAFOLDER")+'/motion/'+str(now.hour)+str(now.minute)+str(now.second)+str(randint(0, 100))+'.jpg',  cv2.cvtColor(frame1, 0))

                                if config_loader.get_value("DEBUG") ==1:
                                    cv2.rectangle(frame1,(x,y),(x+w,y+h), (0,255,245), 2)
    
                        #draw
                        frame3=frame1

                        if config_loader.get_value("DEBUG") == 1:
                            cv2.drawContours(frame3,contours,-1,(0,255,0), 2)
                            cv2.imwrite(config_loader.get_value("DATAFOLDER")+'/debug/motion.jpg', frame3)
    
                        frame1 = frame2
    
                        frame2 = get_image(dir_name,file_name)
    
                        #remove old photo
                        if file_name_to_del!='':
                            try:
                                os.remove(file_name_to_del)
                            except Exception as ex:
                                print(ex)
                        file_name_to_del=os.path.join(dir_name, file_name)
        except Exception as ex:
            pass           


    #cv2.destroyAllWindows()
if __name__ == "__main__":
   main(sys.argv[1:])