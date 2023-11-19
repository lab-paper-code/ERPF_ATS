from fastapi import APIRouter, Header
from starlette import status

import domain.apprun.apprun_crud as apprun_crud

from domain.apprun.apprun_schema import AuthInfo, AppRunExecute, AppRun

router = APIRouter(
    prefix = '/appruns'
)

@router.post("/{app_id}", status_code=status.HTTP_200_OK)
def execute_apprun(app_id: str, apprun_execute: AppRunExecute, dev_id: str = Header(..., alias ='dev_id'), dev_password: str = Header(..., alias ='dev_password')):
    print(f"Incoming Headers: dev_id = {dev_id}, dev_password = {dev_password}")
    print(f"Incoming Body: apprun_execute = {apprun_execute}")
    auth_info = AuthInfo(id = dev_id, password = dev_password)
    return apprun_crud.execute_apprun(auth_info, apprun_execute, app_id)

@router.get("/{apprun_id}", response_model=AppRun)
def get_apprun(apprun_id: str, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return apprun_crud.get_apprun(auth_info, apprun_id)

@router.get("", response_model=list[AppRun])
def list_appruns(dev_id: str = Header(..., alias ='dev_id'), dev_password: str = Header(..., alias ='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return apprun_crud.list_appruns(auth_info)

@router.delete("/{apprun_id}", status_code=status.HTTP_200_OK)
def terminate_appruns(apprun_id: str, dev_id: str = Header(..., alias ='dev_id'), dev_password: str = Header(..., alias ='dev_password')):
    auth_info = AuthInfo(id = dev_id, password = dev_password)
    return apprun_crud.terminate_appruns(auth_info, apprun_id)