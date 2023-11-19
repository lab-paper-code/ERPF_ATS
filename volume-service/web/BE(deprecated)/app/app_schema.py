from pydantic import BaseModel, validator
from datetime import datetime

class AuthInfo(BaseModel):       # used for request authentication
    id : str
    password : str

class AppRegister(BaseModel):    # used for app register
    name : str
    require_gpu : bool | None = None
    description : str | None = None
    docker_image : str
    arguments : str | None = None
    open_ports : list[int] | None = None

    @validator('name', 'docker_image')
    def not_empty(cls, v):
        if not v or not v.strip():
            raise ValueError('빈 값은 허용되지 않습니다.')
        return v
    
class App(BaseModel):            # return type of get/list apps
    id : str
    name : str
    require_gpu : bool | None = None
    description : str | None = None
    docker_image : str
    arguments : str | None = None
    open_ports : list | None = None
    created_at : datetime | None = None
    updated_at : datetime | None = None