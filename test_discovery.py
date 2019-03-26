#!/usr/bin/env python3
"""discovery-node endpoints test
"""

import requests
import provoj
import time

URL0 = "http://127.0.0.1:3000"
URL1 = "http://127.0.0.1:4000"

t = provoj.NewTest("discovery-node")

r = requests.get(URL0 + "/")
t.rStatus("get info from URL0", r)
r = requests.get(URL1 + "/")
t.rStatus("get info from URL1", r)

id0 = {"idAddr": "0x47a2b2353f1a55e4c975b742a7323c027160b4e3"}

r = requests.post(URL0 + "/id", json=id0)
t.rStatus("post id to URL0 " + id0["idAddr"], r)

time.sleep(1)


r = requests.get(URL1 + "/id/" + id0["idAddr"])
t.rStatus("get id to URL1 " + id0["idAddr"], r)

t.printScores()
