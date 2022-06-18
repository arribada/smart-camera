#!/usr/bin/env python

# Device specific patches for Jetson Nano (needs to be before importing cv2)

import cv2
import os
import sys, getopt
import numpy as np
from edge_impulse_linux.image import ImageImpulseRunner
import datetime
import time
runner = None

def help():
    print('python classify-image.py <path_to_model.eim> <path_to_image.jpg>')

def main(argv):
    try:
        opts, args = getopt.getopt(argv, "h", ["--help"])
    except getopt.GetoptError:
        help()
        sys.exit(2)

    for opt, arg in opts:
        if opt in ('-h', '--help'):
            help()
            sys.exit()

    if len(args) != 2:
        help()
        sys.exit(2)

    model = args[0]

    dir_path = os.path.dirname(os.path.realpath(__file__))
    modelfile = os.path.join(dir_path, model)

    print('MODEL: ' + modelfile)

    with ImageImpulseRunner(modelfile) as runner:
        try:
            model_info = runner.init()
            print('Loaded runner for "' + model_info['project']['owner'] + ' / ' + model_info['project']['name'] + '"')
            labels = model_info['model_parameters']['labels']

            while True:
                dir_name=args[1]
                # Get list of all files only in the given directory
                list_of_files = filter( lambda x: os.path.isfile(os.path.join(dir_name, x)),
                            os.listdir(dir_name) )
                # Sort list of files based on last modification time in ascending order
                list_of_files = sorted( list_of_files,
                                        key = lambda x: os.path.getmtime(os.path.join(dir_name, x))
                                        )
                # Iterate over sorted list of files and print file path 
                # along with last modification time of file 
                for file_name in list_of_files:
                    file_path = os.path.join(dir_name, file_name)
                    img = cv2.imread(file_path)
                    if img is None:
                        print('Failed to load image', file_path)
                        exit(1)

                    # imread returns images in BGR format, so we need to convert to RGB
                    img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)

                    # get_features_from_image also takes a crop direction arguments in case you don't have square images
                    features, cropped = runner.get_features_from_image(img)

                    # the image will be resized and cropped, save a copy of the picture here
                    # so you can see what's being passed into the classifier
                    cv2.imwrite('debug.jpg', cv2.cvtColor(cropped, cv2.COLOR_RGB2BGR))

                    res = runner.classify(features)

                    if "classification" in res["result"].keys():
                        print('Result (%d ms.) ' % (res['timing']['dsp'] + res['timing']['classification']), end='')
                        for label in labels:
                            score = res['result']['classification'][label]
                            print('%s: %.2f\t' % (label, score), end='')
                        print('', flush=True)

                    elif "bounding_boxes" in res["result"].keys():
                        print('Found %d bounding boxes (%d ms.)' % (len(res["result"]["bounding_boxes"]), res['timing']['dsp'] + res['timing']['classification']))
                        for bb in res["result"]["bounding_boxes"]:
                            print('\t%s (%.2f): x=%d y=%d w=%d h=%d' % (bb['label'], bb['value'], bb['x'], bb['y'], bb['width'], bb['height']))
                            #cropped = cv2.rectangle(cropped, (bb['x'], bb['y']), (bb['x'] + bb['width'], bb['y'] + bb['height']), (255, 0, 0), 1)
                            if bb['value']>0.7:
                                now = datetime.datetime.now()

                                #save image on folder
                                cv2.imwrite('./detected/'+str(now.hour)+str(now.minute)+str(now.second)+'.jpg',  cv2.cvtColor(cropped, cv2.COLOR_RGB2BGR)) 
                    os.remove(file_path)      
                time.sleep(5) 

        finally:
            if (runner):
                runner.stop()

if __name__ == "__main__":
   main(sys.argv[1:])
