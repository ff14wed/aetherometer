<template>
  <v-tooltip right :disabled="!showTooltip" open-delay="0" close-delay="0" content-class="black">
    <template v-slot:activator="{ on }">
      <v-list-item-avatar size="36" tile v-on="on">
        <slot></slot>
      </v-list-item-avatar>
    </template>
    <span>{{ tooltip }}</span>
  </v-tooltip>
</template>

<script>
export default {
  props: {
    showTooltip: {
      type: Boolean,
      default: false,
    },
    tooltip: {
      type: String,
      default: '',
    },
  },
};
</script>
