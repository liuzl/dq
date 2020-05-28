import requests
import json
from tqdm import tqdm
data = {"title": "这里是标题", "content": "这是内容"}
for i in tqdm(range(10000)):
    data['id'] = i
    r = requests.post("http://127.0.0.1:9080/enqueue/", data={"data": json.dumps(data, ensure_ascii=False)})
    #print(r.text)
