from pydantic import BaseModel, validator
from datetime import datetime

class AuthInfo(BaseModel):          # info used for request authentication
    id : str
    password : str

class DeviceRegister(BaseModel):    # used for device register
    ip : str | None = None
    password : str                  # device password, used with deviceID to login
    description : str | None = None

    @validator('password')
    def not_empty(cls, v):
        if not v or not v.strip():
            raise ValueError('빈 값은 허용되지 않습니다.')
        return v

class DeviceUpdate(BaseModel):
    mod_ip : str | None = None
    mod_password: str | None = None

class Device(BaseModel):            # return type of get/list devices
    id : str
    ip : str
    password : str
    description : str | None = None
    created_at : datetime | None = None
    updated_at : datetime | None = None