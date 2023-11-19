from domain.volume.volume_schema import AuthInfo, VolumeCreate, VolumeUpdate, Volume

import requests
import json
from utils import handle_response
with open('./conf.json','r') as conf:
    config = json.load(conf)

def create_volume(authInfo: AuthInfo, volumeCreate: VolumeCreate) -> dict:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['volumes'][0]                                             # palisade server url
    
    data = {
        'device_id': volumeCreate.device_id,
        'volume_size': volumeCreate.volume_size,
    }
    
    response = requests.post(Serv_url, json=data, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)
    
    response2 = requests.get(Serv_url, auth = (dev_id, dev_password))                   # GET request to get volumeID
    
    if response2.status_code != 200:
        handle_response(response2)
    volume_id = json.loads(response2.text)['volumes'][-1]['id']

    return {"volume_id": volume_id}                                             # return volume_info['id']
    
def get_volume(authInfo: AuthInfo, volume_id: str) -> Volume:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['volumes'][1] + volume_id

    response=requests.get(Serv_url, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text) 

def list_volumes(authInfo: AuthInfo) -> list[Volume]:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['volumes'][0]
    
    response=requests.get(Serv_url, auth=(dev_id, dev_password))
    
    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text)['volumes']

def update_volumes(authInfo: AuthInfo, volumeUpdate: VolumeUpdate, volume_id: str):
    dev_id = authInfo.id
    dev_password = authInfo.password
    volumeUpdate.volume_id = volume_id

    Serv_url = config['volumes'][1] + volume_id

    data = {
        'volume_size': volumeUpdate.mod_volume_size
    }
    
    response=requests.patch(Serv_url, json=data, auth = (dev_id, dev_password))
    
    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text)

def mount_volume(authInfo, volume_id: str):
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['mounts'][0] + volume_id                                           # palisade server url
    
    response = requests.post(Serv_url, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)

    return json.loads(response.text)

def unmount_volume(authInfo, volume_id: str):
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['mounts'][0] + volume_id                                            # palisade server url
    
    response = requests.delete(Serv_url, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)

    return json.loads(response.text)
