<template>
    <div id="app">

        <router-link tag="a"  to="/"><h1 class="text-center">Autoline Asist</h1> </router-link>
        <router-view></router-view>
        <footer><a href="https://github.com/HasOtomotiv/AutolineAssist" target="_blank">Kaynak Kodu</a> | <a href="mailto:acikkaynak@has.com.tr">Görüş ve önerileriniz için</a> | <b-badge>v1.0</b-badge> </footer>
    </div>
</template>

<script>

</script>

<style>
    #app {
        font-family: 'Avenir', Helvetica, Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
        color: #2c3e50;
        margin-top: 60px;
    }
    footer {
        text-align: center;
        font-size: small;
        color : white;
        position: absolute;
        width: 100%;
        clear: both;
        bottom: 0;
        padding: 10px;
        margin: 0;
        /* ATTENTION! The following elements below
        can be set to whatever your heart desires */
        height: 40px; /* REMEMBER height = padding-bottom */
        background: #333;
    }
    footer a:link, footer a:visited, footer a:hover {
        color: white;
    }
</style>
