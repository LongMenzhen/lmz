
--
-- 批量获取当前在线用户的用户名
-- 首先查询 clientids 得到所有的连接cid
-- 其次， user:client:{cid} 获取对应的顾客userid
-- 最后，通过user_id 遍历查询所有的用户信息
--

local clientidskey = KEYS[1]

local names = {}
local clientids = redis.call("smembers", clientidskey)
for index, clientid in pairs(clientids)
do
    local key = "user:client:" .. clientid
    local userid = redis.call("get", key)
    local userkey = "user:" .. userid
    local user = redis.call("get", userkey)
    local detail = cmsgpack.unpack(user)

    table.insert(names, detail['Username'])
end

return names