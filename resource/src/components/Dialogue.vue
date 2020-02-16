<template>
    <div class="main">
        <div class="top">
            <label>账号： </label><input v-model="user_id" />
            <lable>密码： </lable><input v-model="password" />
            <button @click="login">登陆</button><br/>
            <input v-model="group_name" />
            <button @click="createGroup">创建消息组</button>
        </div>
        <div class="left">
            <!-- 内容框 -->
            <div class="content">
                <div v-for="content in contents">{{content}}</div>
            </div>

            <!-- 操作框 -->
            <div class="input">
                <textarea id="txta" v-model="input"></textarea>
            </div>
            <div @click="sendMessage" class="send">
                <span>发送</span>
            </div>
        </div>

        <!--右边-->
        <div class="right">
            <div class="name" v-for="name in names">{{name}}</div>
        </div>
    </div>
</template>

<script>
import "axios";

export default {
    name: 'dialogue',
    data () {
        return {
            user_id: '',
            password: '',
            group_name: '',
            names: [],
            input: '',
            websock: null,
            contents: []
        }
    },
    created () {
        // this.initWebsocket();
    },
    methods: {
        login: function() {
            if (0 >= this.user_id.length || 0 >= this.password.length) {
                alert('请输入登陆账号和密码')
                return
            }

            this.initWebsocket()
        },
        sendMessage: function() {
            if (0 >= this.input.length) {
                alert("输入内容不能为空")
                return
            }
            let data = {"event": "say", "body": {"content": this.input}}
            this.websocketSend(JSON.stringify(data))
            this.input = ''
        },
        createGroup: function() {
            let action = {"event": "create-group", "body": {"name": this.group_name}}
            websocketSend(JSON.stringify(action))
        },
        getGroups: function() {
            
        },
        initWebsocket() {
            const wsuri = 'ws://127.0.0.1:8000/ws';
            this.websock = new WebSocket(wsuri);
            this.websock.onopen = this.websocketOnOpen;
            this.websock.onmessage = this.websocketOnMessage;
            this.websock.onerror = this.websocketOnError;
            this.websock.onclose = this.websocketClose;
        },
        websocketOnOpen() {
            console.log('连接成功')
            let actions = {"event": "login", "body": {"user_id": this.user_id, "password": this.password}};
            this.websock.send(JSON.stringify(actions));
        },
        websocketOnError(ev) {
            console.log("出错了！")
            console.log(ev)
            // this.initWebsocket();
        },
        websocketOnMessage(event) {
            let data = JSON.parse(event.data)
            switch (data['event']) {
                case 'say':
                    let body = data['body']
                    let content = body['name'] + ' ' + body['created_at'] + ':' + body['content']
                    this.contents.push(content)
                    break
                case 'login':
                    let names = data['body']['names']
                    let name = data['body']['name']
                    if (0 >= this.names.length) {
                        this.names.push(name)
                    } else {
                        this.names = names
                    }

                    break;
            }
        },
        websocketSend(Data) {
            this.websock.send(Data)
        },
        websocketClose(e) {
            console.log('断开连接', e)
        }
    }
}
</script>
<style scoped>
* {
    margin: 0;
    padding: 0;
}
.main {
    width: 1024px;
    height: 800px;
    text-align: center;
    border: 1px solid;
    padding: 1px;
}
.top {
    height: 28px;
}
.left {
    float: left;
    width: 800px;
}
.right {
    float: right;
    width: 200px;
    border: 1px solid;
    height: 768px;
}

.right .name {
    border: 1px solid;
    text-align: right;
}

.content {
    width: 800px;
    height: 600px;
    border: 1px solid;
    text-align: left;
    line-height: 24px;
    /* margin-left: 10px; */
}
.input {
    margin-top: 10px;
    width: 600px;
    height: 160px;
    border: 1px solid;
    float: left;
}
.send {
    margin-top: 10px;
    /* margin-left: 30px; */
    width: 180px;
    height: 160px;
    border: 1px solid;
    float: right;
    display: table;
}
.send span {
    text-align: center;
    display: table-cell;
    vertical-align: middle;
}
#txta {
    margin: 0px;
    height: 158px;
    width: 598px;
    font-size: 16px;
    font-family: monospace;
}
</style>
