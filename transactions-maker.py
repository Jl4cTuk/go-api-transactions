import requests
import json
import random

rng = range(1, 15, 1)
wlts = (
"3b489658",
"dd0d087a",
"c41510ac",
"5ffc3974",
"b3846c39",
"9ac28c41",
"bdf99b10",
"f614ab75",
"47b5cd35",
"8484198a")

url = "http://localhost:8080/api/send"

for i in range(10):
    payload = json.dumps({"from": random.choice(wlts),
                        "to": random.choice(wlts),
                        "amount": random.choice(rng)})

    headers = {
        "Accept": "*/*",
        "Accept-Encoding": "gzip, deflate, br",
        "User-Agent": "EchoapiRuntime/1.1.0",
        "Connection": "keep-alive",
        "Content-Type": "application/json"
    }

    response = requests.request("POST", url, data=payload, headers=headers)

    print(response.text)
