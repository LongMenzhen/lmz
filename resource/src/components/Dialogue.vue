<template>
    <div class="main">
        <div class="left">
            <!-- 内容框 -->
            <div class="content">
                <div v-for="content in contents" v-bind:key="content" >{{content}}</div>
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
        <div class="right"></div>
    </div>
</template>

<script>
export default {
    name: 'dialogue',
    data () {
        return {
            input: '',
            websock: null,
            contents: []
        }
    },
    created () {
        this.initWebsocket();
    },
    methods: {
        sendMessage: function () {
            // console.log(this.contents)
            // this.contents.push(this.input)
            let data = {"event": "say", "group_id": 1, "body": {"content": this.input}}
            this.websocketSend(JSON.stringify(data))
            this.input = ''
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
            let actions = {"event": "login", "body": {"name": "Tom"}};
            this.websock.send(JSON.stringify(actions));
        },
        websocketOnError() {
            this.initWebsocket();
        },
        websocketOnMessage(e) {
            // const redata = JSON.parse(e.data)
            console.log(123)
            console.log(e)
            this.contents.push(e.data)
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
.left {
    float: left;
    width: 800px;
}
.right {
    float: right;
    width: 200px;
    border: 1px solid;
    height: 798px;
}
.content {
    width: 800px;
    height: 600px;
    border: 1px solid;
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
    height: 154px;
    width: 596px;
    font-size: 16px;
    font-family: monospace;
}
</style>
