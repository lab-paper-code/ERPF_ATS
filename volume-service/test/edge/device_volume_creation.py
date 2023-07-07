''' 
목표: device_volume 생성
구현할 함수: 볼륨 등록(volume_register_request), 상태 출력(handle_response)
input: DeviceID, VolumeSize
output: 볼륨 생성, 상태 출력
'''
import requests

def volume_register_request(Serv_url, DeviceID, VolumeSize):
    data = {
        'DeviceID': DeviceID,
        'VolumeSize': VolumeSize
    }
    response = requests.post(Serv_url, data=data) # post request
    return response

Serverurl = '155.230.36.25:31200/volumes' # set serverURL, change this to var.
print("볼륨 생성을 위한 정보를 받습니다.") #순서: 볼륨 생성 요청 -> volume_handler가 볼륨 생성 처리->생성된 볼륨이 list로 등록?
DeviceID=input("DeviceID를 입력하세요: ")
VolumeSize=input("요청할 볼륨 크기를 입력하세요: ")
response = volume_register_request(Serverurl, DeviceID, VolumeSize)

def handle_response(response): # print status -> GPU 쓰는 거
    if response.status_code == 200:
        print("볼륨이 정상적으로 생성되었습니다.") 
    elif response.status_code == 400:
        print("Bad Request: 잘못된 요청입니다.")
    elif response.status_code == 401:
        print("Unauthorized: 인증되지 않은 요청입니다.")
    elif response.status_code == 403:
        print("Forbidden: 접근이 금지되었습니다.")
    elif response.status_code == 404:
        print("Not Found: 요청한 리소스를 찾을 수 없습니다.")
    else:
        print("Unknown Error: 알 수 없는 에러가 발생했습니다.")

handle_response(response)