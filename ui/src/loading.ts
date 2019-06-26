import 'typeface-roboto';

import Vue from 'vue';
import './plugins/vuetify';
import Loading from './components/Loading.vue';

Vue.config.productionTip = false;

new Vue({
  render: (h) => h(Loading),
}).$mount('#app');
