import Vue from 'vue'
import App from './App.vue'
import BootstrapVue from 'bootstrap-vue'
import VueRouter from 'vue-router'
/*
import { library } from '@fortawesome/fontawesome-svg-core'
import { faUserSecret } from '@fortawesome/free-solid-svg-icons'
*/
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { routes } from "./routes";

Vue.config.productionTip = false

const moment = require('moment')
require('moment/locale/tr')
Vue.use(require('vue-moment'),
    {
      moment
    }
);
Vue.use(BootstrapVue);
Vue.use(VueRouter);
Vue.component('font-awesome-icon', FontAwesomeIcon)

const router = new VueRouter({
  base:'/',
  routes,
  mode : 'history'
});

new Vue({
  el: '#app',
  router,
  render: h => h(App)
})/*.$mount('#app')*/
