import os, sys, time, io

import argparse
import numpy as np
from tqdm import tqdm

from keras.applications.inception_v3 import InceptionV3
from keras.applications.inception_v3 import preprocess_input
from tensorflow.keras.utils import img_to_array
from tensorflow.keras.utils import load_img
from keras.applications import imagenet_utils

from tensorflow.python.keras.models import load_model
from tensorflow.keras.applications import MobileNet
from tensorflow.keras.applications import InceptionV3
from silence_tensorflow import silence_tensorflow

from config import config

silence_tensorflow()    # Hide warnning messages of tensorflow
model = MobileNet(weights="imagenet")

def printConfiguration():
    print('Configuration: {')
    for key, value in config.items():
        print('\t' + key + ': ' + str(value))
    print('}')

def isImageFile(fileName):
    fileName = fileName.upper()
    validFormat = ['JPEG', 'JPG', 'PNG']
    for f in validFormat:
        if fileName.endswith('.' + f):
            return True
    return False

def getFilesByCnt(cnt):
    print('[ ] Get ' + str(cnt) + ' files')
    try:
        imageFileList = [element for element in os.listdir(IMAGE_PATH) if isImageFile(element)]
        targetFileList = []

        for i in range(0, cnt):
            targetFileList.append(IMAGE_PATH + '/' + imageFileList[i])

        print(' - ' + '\n - '.join(targetFileList))

        return targetFileList

    except Exception as e:
        print('[E] ' + str(e))
    
def getFiles(limitSize):
    try:
        fileList = [IMAGE_PATH + '/' + fileName for fileName in os.listdir(IMAGE_PATH)]
        totalFileSizes = 0

        for i in range(len(fileList)):
            currentFileName = fileList
            currentFileSize = os.path.getsize(fileList[i]) / (1024.0 * 1024.0 * 1024.0)

            if totalFileSizes + currentFileName > limitSize:
                print("TOTAL FILE SIZE: ", totalFileSizes)
                return fileList[:i]
            
            totalFileSizes += currentFileSize
        print("SEND FILE LIST LENGTH: {}".format(len(fileList)))
        return fileList

    except Exception as e:
        print(str(e))

def predictFu(path, model, outputFile):
    try:
        inputShape = (224, 224)
        image = load_img(path, target_size=inputShape)
        imageArray = img_to_array(image)
        imageArray = np.expand_dims(imageArray, axis=0)
        imageArray = preprocess_input(imageArray)

        predictions = imagenet_utils.decode_predictions(model.predict(imageArray))

        for (i, (imagenetID, label, prob)) in enumerate(predictions[0]):
            outputFile.write("{}: {:.2f}% \n".format(label, prob * 100))
            return imagenetID, label

    except Exception as e:
        print(str(e))

def imageClassification(cnt):
    outputFile = f"{RESULT_PATH}result_{IMG_CNT}.txt"
    imageFileList = getFilesByCnt(int(cnt))

    if os.path.isfile(outputFile):
        print('[ ] Log file ' + outputFile + ' already exists. Deleting the existing file.')
        os.unlink(outputFile)

    with open(outputFile, 'w') as f:
        start = time.time()
        for imageFile in tqdm(imageFileList):
            predictFu(imageFile, model, f)

    print("[ ] Elapsed time: %s" % (time.time() - start))
    print('[ ] A log file has been created: ' + outputFile)

if __name__ == "__main__" :
    print('[ ] Starting... ')

    print('[ ] Parsing arguments')
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", dest="img_cnt", action="store", default=10, type=int)
    config['img_cnt'] = parser.parse_args().img_cnt
    IMAGE_PATH = config['img_path']
    RESULT_PATH = config['result_path']
    RESULT_PATH = RESULT_PATH if RESULT_PATH[-1] == "/" else RESULT_PATH + "/"
    IMG_CNT = config['img_cnt']

    printConfiguration()

    print('[ ] Start the image classification')
    print('[ ] Buffer size:', io.DEFAULT_BUFFER_SIZE, 'Bytes')
    imageClassification(IMG_CNT)