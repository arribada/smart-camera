import json
from os import environ


class Singleton(type):
    _instances = {}
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            cls._instances[cls] = super(Singleton, cls).__call__(*args, **kwargs)
        return cls._instances[cls]

class Config(metaclass=Singleton):
    def set_data(cls, data):
        Config.data=data

def load_config(path):
    f = open(path)
    data = json.load(f)
    Config().set_data(data)


def get_value(value_name):
    if environ.get(value_name) is not None:
        return environ.get(value_name)
    else:
        values=value_name.split("_")
        data_value=Config.data
        for value in values:
            if value in data_value.keys():
                data_value=data_value[value]
            else:
                data_value=None
                break
        return data_value



#load_config("C:\\Users\\Windows\\PycharmProjects\\loadConfig\\Config.json")
#print(get_value("CAPTURE_PICAMERA_FPS"))
#print(get_value("DATAFOLDER"))
#print(get_value("CAPTURE_STREAM_FPS"))