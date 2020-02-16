
--
-- 根据groupid 返回下面所有的所有username
-- 查询组下的所有user_ids
-- 遍历所有的user_id 对应的用户信息，得到他们的username
--

local group_key = KEYS[1]

local names = {}
local userids = redis.call("smembers", group_key)
for index, userid in pairs(userids)
do
    local user_key = "user:" .. userid
    local msg = redis.call("get", user_key)
    local user = msgpack.unpack(msg)

    table.insert(names, user['username'])
end

return names