<template>
  <v-list-item
    @click="state.selectNav(navID)"
    active-class="active-item"
    :class="{ 'active-item': navID == state.selectedNav }"
  >
    <v-nav-icon :showTooltip="showTooltip" :tooltip="name">
      <slot></slot>
    </v-nav-icon>
    <v-list-item-content>
      <v-list-item-title class="font-weight-light">{{ name }}</v-list-item-title>
    </v-list-item-content>
  </v-list-item>
</template>

<script>
import CommonStore from '../stores/commonStore';
import NavIcon from './NavIcon.vue';

import { observer } from 'mobx-vue';
export default observer({
  props: {
    state: {
      type: CommonStore,
      required: true,
    },
    navID: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    showTooltip: {
      type: Boolean,
      default: false,
    },
  },
  components: {
    'v-nav-icon': NavIcon,
  },
});
</script>

<style>
.active-item > a {
  background-color: #1976d2 !important;
  border-color: #1976d2 !important;
}
</style>
