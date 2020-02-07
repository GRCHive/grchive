import requests
import os

url = "https://gitlab.com/api/v4/projects/grchive%2Fgrchive/variables"
token = os.environ["GITLAB_API_TOKEN"]

r = requests.get(url=url, headers={
    "PRIVATE-TOKEN" : token  
})
data = r.json()

for datum in data:
    key = datum['key']
    val = datum['value'].replace("\"", "\\\"")
    print("export {}=\"{}\"".format(key, val))
