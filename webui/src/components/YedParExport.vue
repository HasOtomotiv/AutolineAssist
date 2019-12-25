<template>
    <b-container>
        <b-row class="justify-content-md-center">
            <b-col cols="12" sm="6">
                <h3 class="text-center">Yedek Parça Transfer Dosyası</h3>
                <br/>
                <b-form>
                    <b-form-group label-cols-sm="4" label-cols-lg="3" label="WIP / Sipariş" label-for="wippo">
                        <b-form-select class="w-10 mb-2 mr-sm-3 mb-sm-0" id="wippo" v-model="wippo"
                                       :options="wipoptions"></b-form-select>
                    </b-form-group>
                    <b-form-group label-cols-sm="4" label-cols-lg="3" label="No" label-for="yedparid">
                        <b-input id="yedparid" class="mb-2 mr-sm-3 mb-sm-0" v-model="yedparid"
                                 type="number" name="yedparid"></b-input>
                    </b-form-group>
                    <b-form-group label-cols-sm="4" label-cols-lg="3" label="Lokasyon Kodu" label-for="branchid">
                        <b-input autocomplete="on" id="branchid" class="mb-2 mr-sm-3 mb-sm-0" v-model="branchid"
                                 type="number" name="brancid"></b-input>
                    </b-form-group>
                    <b-form-group label-cols-sm="4" label-cols-lg="3" label="Dosya Formatı" label-for="doctype">
                        <b-form-select id="doctype" v-model="doctype" :options="docoptions"></b-form-select>
                    </b-form-group>
                    <b-button @click.prevent="submit" block variant="primary bloc">İndir</b-button>
                </b-form>
            </b-col>
        </b-row>
    </b-container>
</template>

<script>
    import axios from "axios"

    export default {
        data() {
            return {
                yedparid: '',
                doctype: '',
                branchid: '',
                wippo: 'P',
                Error: '',
                wipoptions: [
                    {text: 'Sipariş', value: 'P'},
                    {text: 'WIP', value: 'W'},
                ],
                docoptions: [
                    {value: 'so', text: 'SO'},
                    {value: 'csv', text: 'CSV'},
                    {value: 'xfr', text: 'XFR'},
                ]
            }
        },
        methods: {
            submit() {
                let url = '';
                switch (this.wippo) {
                    case "P":
                        url = '/api/getporecords/' + this.doctype + '/' + this.branchid + '/' + this.yedparid;
                        break;
                    case "W":
                        url = '/api/getwiprecords/' + this.doctype + '/' + this.branchid + '/' + this.yedparid;
                }
                axios({
                    url: url,
                    method: 'GET',
                    responseType: 'blob',
                }).then((response) => {
                    var fileURL = window.URL.createObjectURL(new Blob([response.data]));
                    var fileLink = document.createElement('a');

                    fileLink.href = fileURL;
                    fileLink.setAttribute('download', this.wippo + this.yedparid + '.' + this.doctype);
                    document.body.appendChild(fileLink);

                    fileLink.click();
                });
            }
        }
    }
</script>

<style scoped>
    /* Chrome, Safari, Edge, Opera */
    input::-webkit-outer-spin-button,
    input::-webkit-inner-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }

    /* Firefox */
    input[type=number] {
        -moz-appearance: textfield;
    }
</style>
