from socket import SocketKind
from time import sleep
import psutil
from psutil._common import sconn
import json
import requests
import os
import gc
from typing import List,Set

REMOTE_ENDPOINT =  os.environ["HOST"] if "HOST" in os.environ else "http://localhost:2709/monitor"
POD_NAME = os.environ["POD_NAME"] if "POD_NAME" in os.environ else "aName"

def to_dict(conns: List[sconn]):
    return json.loads(json.dumps([x._asdict() for x in conns]))
def to_json(conns: List[sconn]):
    return json.dumps([x._asdict() for x in conns]) 
def try_send_monitor(conns: List[sconn]):
    try:
       requests.post(REMOTE_ENDPOINT, json=to_dict(conns),headers={
           "POD_NAME" : POD_NAME
       })
    except Exception as e :
        print(e)
        pass 

def with_filter(conns: List[sconn]):
    return list(filter(lambda x : x.status == "LISTEN" or (x.type == SocketKind.SOCK_DGRAM and x.raddr == ()), conns))

def update_connections(prev_conns: List[sconn]):
    curr_conns: List[sconn] = with_filter(psutil.net_connections(kind='all'))
    new_conns_list: List[sconn] = list(set(curr_conns) - set(prev_conns))
    new_registered_conns = list(set(prev_conns + new_conns_list))
    del curr_conns
    del new_conns_list
    return new_registered_conns

conns: List[sconn] = with_filter(psutil.net_connections(kind='all'))
#try_send_monitor(conns)
#print(to_dict(conns))
while True:
    conns = update_connections(conns)
    #try_send_monitor(conns)
    print(to_json(conns),flush=True)
    del conns
    gc.collect()
    sleep(10)
