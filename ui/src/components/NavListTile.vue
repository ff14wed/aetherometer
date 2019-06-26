<template>
  <v-list-tile
    @click="state.selectNav(navID)"
    active-class="active-tile"
    :class="{ 'active-tile': navID == state.selectedNav }"
  >
    <v-nav-icon :showTooltip="showTooltip" :tooltip="name">
      <slot></slot>
    </v-nav-icon>
    <v-list-tile-content>
      <v-list-tile-title class="font-weight-light">{{ name }}</v-list-tile-title>
    </v-list-tile-content>
  </v-list-tile>
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
  .v-list__tile {
    margin-left: 10px;
    margin-right: 5px;
    padding: 0 6px;
  }

  .active-tile > a {
    background-color: #1976d2 !important;
    border-color: #1976d2 !important;
  }
</style>