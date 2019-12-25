<template>
    <b-container>
        <b-row class="justify-content-md-center">
            <b-col cols="6">
                <h3>Uygulamalar</h3>
                <br>
                <b-list-group>
                    <router-link tag="a" class="list-group-item list-group-item-action" to="/efaturainfo">e-Fatura
                        Bilgileri
                    </router-link>
                    <router-link tag="a" class="list-group-item list-group-item-action" to="/yedparexport">Yedek
                        Parça Transfer
                    </router-link>
                    <a href="#" class="list-group-item list-group-item-action disabled">Müşteri e-Fatura Bilgileri
                        <small>(Yakında)</small>
                    </a>
                </b-list-group>
            </b-col>
        </b-row>
    </b-container>
</template>

<script>
    export default {
        name: 'Home',
        props: {
            msg: String
        }
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    h3 {
        margin: 40px 0 0;
    }

    ul {
        list-style-type: none;
        padding: 0;
    }

    a {
        color: #42b983;
    }
</style>
