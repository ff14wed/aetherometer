import 'typeface-roboto';

import Vue from 'vue';
import Loading from './components/Loading.vue';
import vuetify from './plugins/vuetify';

Vue.config.productionTip = false;

new Vue({
  vuetify,
  render: (h) => h(Loading),
}).$mount('#app');

