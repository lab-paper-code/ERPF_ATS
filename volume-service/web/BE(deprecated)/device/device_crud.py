from domain.device.device_schema import AuthInfo, DeviceRegister, DeviceUpdate, Device

import requests
import json
from utils import handle_response
with open('./conf.json','r') as conf:
    config = json.load(conf)


def register_device(authInfo: AuthInfo, deviceRegister: DeviceRegister) -> dict:
    admin_id = authInfo.id
    admin_password = authInfo.password
    Serv_url = config['devices'][0]                                             # palisade server url
    
    data = {
        'ip': deviceRegister.ip,
        'password': deviceRegister.password,
        'description': deviceRegister.description, 
    }
    
    response = requests.post(Serv_url, json=data, auth = (admin_id, admin_password))

    if response.status_code != 200:
        handle_response(response)
    
    response2 = requests.get(Serv_url, auth = (admin_id, admin_password))                   # GET request to get DeviceID
    
    if response2.status_code != 200:
        handle_response(response2)
    device_id = json.loads(response2.text)['devices'][-1]['id']

    return {"device_id": device_id}                                             # return device_info['id']
    
def get_device(authInfo: AuthInfo, device_id: str) -> Device:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['devices'][1] + device_id

    response=requests.get(Serv_url, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text) 

def list_devices(authInfo: AuthInfo) -> list[Device]:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['devices'][0]
    
    response=requests.get(Serv_url, auth=(dev_id, dev_password))
    
    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text)['devices']

def update_device(authInfo: AuthInfo, deviceUpdate: DeviceUpdate, device_id: str):
    dev_id = authInfo.id
    dev_password = authInfo.password

    Serv_url = config['devices'][1] + device_id

    data = {
        'ip': deviceUpdate.mod_ip,
        'password': deviceUpdate.mod_password
    }
    
    response=requests.patch(Serv_url, json=data, auth = (dev_id, dev_password))
    
    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text)