import requests, sys
sys.path.append('/home/palisade1/Desktop/ksv/volume-service/test')
from utils import get_login_info, handle_response

def volume_register_request(Serv_url, DeviceID, VolumeSize, ID, PASSWD):
    data = {
        'device_id': DeviceID,
        'volume_size': VolumeSize
    }
    response = requests.post(Serv_url, json=data, auth=(ID,PASSWD)) # post request
    return response

if __name__ == "__main__":
    ID, PASSWORD=get_login_info() # get LoginInfo
    Serverurl = "http://155.230.36.25:31200/volumes"

    print("볼륨 생성을 위한 정보를 받습니다.")
    deviceid=input("디바이스 ID를 입력하세요: ")
    volumesize=input("요청할 볼륨 크기를 입력하세요: ") # need err handling?
    response = volume_register_request(Serverurl, deviceid, volumesize, ID, PASSWORD)

    handle_response(response)