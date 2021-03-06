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
t.rStatus("check root endpoint from Node0", r)
r = requests.get(URL1 + "/")
t.rStatus("check root endpoint from Node1", r)

id0 = {"idAddr": "0x47a2b2353f1a55e4c975b742a7323c027160b4e3"}

r = requests.post(URL0 + "/id", json=id0)
t.rStatus("post id0 (" + id0["idAddr"] + ") to Node0 ", r)

time.sleep(1)


# get the data from the id0, as the discovery-node don't have the id data, will ask for it over Pss Swarm network
r = requests.get(URL1 + "/id/" + id0["idAddr"])
t.rStatus("get id0 (" + id0["idAddr"] + ") from Node1 ", r)
jsonR = r.json()
print(jsonR)

time.sleep(1)

# get again the data from the id0, this time, the discovery-node will return it from its db cache
r = requests.get(URL1 + "/id/" + id0["idAddr"])
t.rStatus("get id0 (" + id0["idAddr"] + ") from Node1 ", r)
jsonR2 = r.json()
print(jsonR2)

t.equal("checking that first and second get of the identity returns the same", jsonR, jsonR2)

t.printScores()

print("to make another single identity discover petition, use:\ncurl http://127.0.0.1:4000/id/0x47a2b2353f1a55e4c975b742a7323c027160b4e3 | jq")
