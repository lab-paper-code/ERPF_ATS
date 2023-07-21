import requests, sys
sys.path.append('/home/palisade1/Desktop/ksv/volume-service/test')
from utils import get_login_info, handle_response

def device_register_request(Serv_url, IP, PASSWORD, Description, ID, PASSWD):
    data = {
        'ip': IP,
        'password': PASSWORD,
        'description': Description
    }
    response = requests.post(Serv_url, json=data, auth=(ID,PASSWD)) # post request
    return response

if __name__ == "__main__":
    ID, PASSWORD=get_login_info() # get LoginInfo
    Serverurl = "http://155.230.36.25:31200/devices" # set serverURL, change this to var.

    print("디바이스 등록을 위한 정보를 받습니다.")
    ip=input("IP를 입력하세요: ")
    password=input("Password를 입력하세요: ")
    description=input("디바이스에 대한 설명을 입력하세요: ")
    response = device_register_request(Serverurl, ip, password, description, ID, PASSWORD)
    
    handle_response(response)

    