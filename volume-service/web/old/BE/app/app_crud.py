from domain.app.app_schema import AuthInfo, AppRegister, App

import requests
import json
from utils import handle_response
with open('./conf.json','r') as conf:
    config = json.load(conf)


def register_app(authInfo: AuthInfo, appRegister: AppRegister) -> dict:
    admin_id = authInfo.id
    admin_password = authInfo.password
    Serv_url = config['apps'][0]                                             # palisade server url
    
    data = {
        'name': appRegister.name,
        'require_gpu': appRegister.require_gpu,
        'description': appRegister.description,
        'docker_image': appRegister.docker_image,
        'arguments': appRegister.arguments,
        'open_ports': appRegister.open_ports,
        
    }
    
    response = requests.post(Serv_url, json=data, auth = (admin_id, admin_password))

    if response.status_code != 200:
        handle_response(response)
    
    response2 = requests.get(Serv_url, auth = (admin_id, admin_password))                   # GET request to get DeviceID
    
    if response2.status_code != 200:
        handle_response(response2)
    app_id = json.loads(response2.text)['apps'][-1]['id']

    return {"app_id": app_id}                                             # return app_info['id']
    
def get_app(authInfo: AuthInfo, app_id: str) -> App:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['apps'][1] + app_id

    response=requests.get(Serv_url, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text) 

def list_apps(authInfo: AuthInfo) -> list[App]:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['apps'][0]
    
    response=requests.get(Serv_url, auth=(dev_id, dev_password))
    
    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text)['apps']