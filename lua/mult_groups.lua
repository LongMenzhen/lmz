
--
-- 获取所有的组
-- 返回组信息
-- 

local key = KEYS[1]

local groups = {}
local group_ids = redis.call("smembers", key)
for index, group_id in pairs(group_ids) 
do
    local group_key = "group:" .. group_id
    local msg = redis.call("get", group_key)

    table.insert(groups, msg)
end

return groups