<template>
    <b-container>
        <b-row class="justify-content-md-center">
            <b-col cols="12" sm="6">
                <h3 class="text-center">E-Fatura Bilgi</h3>
                <br/>
                <b-form inline @submit.prevent="submit">
                    <label class="mb-2 mr-sm-3 mb-sm-0">Autoline Dok.No:</label>
                    <b-input id="docid" class="mb-2 mr-sm-3 mb-sm-0" v-model="docid" type="number"></b-input>
                    <b-button @click.prevent="submit" variant="primary">Ara</b-button>
                </b-form>
                <div class="text-center" v-if="loading">
                    <font-awesome-icon icon="spinner" class="fa-spin fa-2x"/>
                </div>
            </b-col>
        </b-row>
        <br>
        <b-row class="justify-content-md-center">
            <b-col cols="12" sm="6">
                <div v-if="CompanyData.data">
                    <table class="table table-bordered table-sm">
                        <tbody>
                        <tr>
                            <td><strong>Müşteri Adı/Ünvanı</strong></td>
                            <td>{{CompanyData.data.CustCompanyName}}</td>
                        </tr>
                        <tr>
                            <td><strong>Fatura No</strong></td>
                            <td>{{CompanyData.data.InvoiceID}}</td>
                        </tr>
                        <tr>
                            <td><strong>Fatura Tutarı</strong></td>
                            <td>{{CompanyData.data.InvoiceTotal}}</td>
                        </tr>
                        <tr>
                            <td><strong>Fatura Tarihi</strong></td>
                            <td>{{CompanyData.data.PostingDate | moment("DD-MM-YYYY")}}</td>
                            <!--2019-08-05T00:00:00+03:00-->
                        </tr>
                        <tr>
                            <td><strong>Posta Kutusu</strong></td>
                            <td>{{CompanyData.data.ReceiverEmail}}</td>
                        </tr>
                        <tr>
                            <td><strong>ERP Referans No</strong></td>
                            <td>{{CompanyData.data.UniqueDocKey}}</td>
                        </tr>
                        <tr>
                            <td><strong>ETTN</strong></td>
                            <td>{{CompanyData.data.UniqueDocRef}}</td>
                        </tr>
                        <tr>
                            <td><strong>WIP NO</strong></td>
                            <td>{{CompanyData.data.WIPNumber}}</td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            </b-col>
        </b-row>
    </b-container>
</template>

<script>
    import axios from "axios"
    import {library} from '@fortawesome/fontawesome-svg-core'
    import {faSpinner, faAlignLeft} from '@fortawesome/free-solid-svg-icons'

    library.add(faSpinner, faAlignLeft)

    export default {
        data() {
            return {
                docid: '',
                Error: '',
                CompanyData: {
                    CustCompanyName: '',
                    InvoiceID: '',
                    InvoiceTotal: 0,
                    PostingDate: '',
                    ReceiverEmail: '',
                    UniqueDocKey: '',
                    UniqueDocRef: '',
                    WIPNumber: 0,
                },
                loading: false,
            }
        },
        methods: {
            submit() {
                this.Error = '';
                this.CompanyData = {};
                this.loading = true;
                axios
                    .get('/api/einvoiceinfo/' + this.docid)
                    .then(response => {
                        if (response.data.data != null) {
                            this.CompanyData = response.data
                        } else {
                            this.Error = 'Kayıt bulunamadı.';
                        }
                    })
                    .catch(error => {
                        this.Error = error
                    })
                    .finally(() => this.loading = false)
                this.docid = null;
            }
        }
    }
</script>

<style scoped>
    input[type=number]::-webkit-inner-spin-button,
    input[type=number]::-webkit-outer-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }
</style>
