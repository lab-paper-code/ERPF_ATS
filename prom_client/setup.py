import socket
import argparse
import pip
import os

# install package with pip
def install(package): 
    if hasattr(pip, 'main'):
        pip.main(['install', package])
    else:
        pip._internal.main(['install', package])

def create_directory(directory_path):
    # 디렉토리가 존재하지 않으면 생성
    if not os.path.exists(directory_path):
        os.makedirs(directory_path)
        print(f"디렉토리 생성: {directory_path}")
    else:
        print(f"디렉토리 이미 존재: {directory_path}")

install("prometheus-client")
# install("uvicorn") # seems unnecessary, installing fastapi also installs uvicorn
install("fastapi")
install("ultralytics")


# command line argument parsing
parser = argparse.ArgumentParser()
parser.add_argument("--image_path", type=str, action="store")          
parser.add_argument("--output_path", type=str, action="store")           
parser.add_argument("--port",type=int, default=8000)
args = parser.parse_args()

# get IP
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.connect(("8.8.8.8", 80))
IP = s.getsockname()[0]
s.close()

# set INPUT_PATH, OUTPUT_PATH, Port specified in command
YOLO_INPUT_PATH = args.image_path
YOLO_OUTPUT_PATH = args.output_path
Port = args.port

create_directory(YOLO_OUTPUT_PATH)
create_directory(f"{YOLO_OUTPUT_PATH}/predict")

# write config options to config.py
with open("config.py","w") as f:
    f.write(f'YOLO_INPUT_PATH = "{YOLO_INPUT_PATH}"\n')
    f.write(f'YOLO_OUTPUT_PATH = "{YOLO_OUTPUT_PATH}"\n')
    f.write(f'IP = "{IP}"\n')
    f.write(f'Port = {Port}\n')