from fastapi import APIRouter, Header
from starlette import status

import domain.app.app_crud as app_crud

from domain.app.app_schema import AuthInfo, AppRegister, App

router = APIRouter(
    prefix = '/apps'
)

POST : /user/create

@router.post("")
def register_app(app_register: AppRegister, admin_id: str = Header(..., alias='admin_id'), admin_password: str = Header(..., alias='admin_password')):
    auth_info = AuthInfo(id = admin_id, password = admin_password)
    return app_crud.register_app(auth_info, app_register)

@router.get("/{app_id}", response_model=App)
def get_app(app_id: str, dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return app_crud.get_app(auth_info, app_id)

@router.get("", response_model=list[App])
def list_apps(dev_id: str = Header(..., alias='dev_id'), dev_password: str = Header(..., alias='dev_password')):
    auth_info = AuthInfo(id=dev_id, password=dev_password)
    return app_crud.list_apps(auth_info)