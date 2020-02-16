
--
-- 获取所有的组
-- 返回组信息
--
--

local key = KEYS[1]

local groups = {}
local group_ids = redis.call("smembers", key)
for index, group_id in pairs(group_ids) 
do
    group_key = "group:" .. group_id
    msg = redis.call(group_key)
    group = msgpack.unpack(msg)

    table.insert(groups, group)
end

return groups