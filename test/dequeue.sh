# 查看队列头的数据
curl "http://127.0.0.1:9080/dequeue/?peek=true"

# 取走队列头数据
# curl "http://127.0.0.1:9080/dequeue/"
