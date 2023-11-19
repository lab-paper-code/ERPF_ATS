from domain.apprun.apprun_schema import AuthInfo, AppRunExecute, AppRun

import requests
import json
from utils import handle_response
with open('./conf.json','r') as conf:
    config = json.load(conf)


def execute_apprun(authInfo: AuthInfo, appRunExecute: AppRunExecute, app_id: str) -> dict:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['appruns'][1] + app_id                                            # palisade server url
    
    data = {
        'device_id': appRunExecute.device_id,
        'volume_id': appRunExecute.volume_id
    }
    
    response = requests.post(Serv_url, json=data, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)

    Serv_url2 = config['appruns'][0]
    response2 = requests.get(Serv_url2, auth = (dev_id, dev_password))                   # GET request to get DeviceID
    
    if response2.status_code != 200:
        handle_response(response2)
    appruns_id = json.loads(response2.text)

    return {"appruns_id": appruns_id}                                             # return device_info['id']


def get_apprun(authInfo: AuthInfo, apprun_id: str) -> AppRun:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['appruns'][1] + apprun_id

    response=requests.get(Serv_url, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text) 

def list_appruns(authInfo: AuthInfo) -> list[AppRun]:
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['appruns'][0]
    
    response=requests.get(Serv_url, auth=(dev_id, dev_password))
    
    if response.status_code != 200:
        handle_response(response)  

    return json.loads(response.text)['app_runs']

def terminate_appruns(authInfo: AuthInfo, apprun_id: str):
    dev_id = authInfo.id
    dev_password = authInfo.password
    Serv_url = config['appruns'][1] + apprun_id                                            # palisade server url
    
    response = requests.delete(Serv_url, auth = (dev_id, dev_password))

    if response.status_code != 200:
        handle_response(response)