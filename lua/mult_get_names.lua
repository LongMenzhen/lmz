
--
-- 批量获取当前在线用户的用户名
-- 首先查询 user:online:userids 得到所有的userids
--  再通过user_id 遍历查询所有的用户信息
--

local user_online_keys = KEYS[1]

local names = {}
local userids = redis.call("smembers", user_online_keys)
for index, userid in pairs(userids)
do
    local key = "user:client:" .. userid
    local userkey = "user:" .. userid
    local user = redis.call("get", userkey)
    local detail = cmsgpack.unpack(user)

    table.insert(names, detail['Username'])
end

return names