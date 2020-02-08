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
            contents: []
        }
    },
    methods: {
        sendMessage: function () {
            console.log(this.contents)
            this.contents.push(this.input)
            this.input = ''
        }
    }
}
</script>
<style scoped>
* {
    margin: 0;
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
