import 'typeface-roboto';

import Vue from 'vue';
import './plugins/vuetify';
import App from './components/App.vue';

Vue.config.productionTip = false;

new Vue({
  render: (h) => h(App),
}).$mount('#app');
