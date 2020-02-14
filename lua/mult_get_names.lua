
--
-- 批量获取当前在线用户的用户名
-- 首先查询 clientids 得到所有的连接cid
-- 其次， user:client:{cid} 获取对应的顾客userid
-- 最后，通过user_id 遍历查询所有的用户信息
--

local clientidskey = KEYS[1]

names = {}
clientids = redis.call("smembers", clientidskey)
for i = 0, #clientids 
do
    clientid = clientids[i]
    key = "user:client:" .. clientid
    userid = redis.call("get", key)
    userkey = "user:" .. userid
    user = redis.call("get", userkey)
    detail = cmsgpack.unpack(user)

    table.insert(names, detail['username'])
end

return names