<template>
  <a :href="`${docsUrl}`" target="_blank" class="text-blue-500 hover:text-blue-600 cursor-pointer mt-1"
    ><Icon name="question" class="!w-4 !h-4"
  /></a>
</template>

<script lang="ts" setup>
import { computed, defineProps, toRef } from 'vue';

import Icon from '~/components/atomic/Icon.vue';

const props = defineProps({
  url: {
    type: String,
    required: true,
  },
});

const docsBaseUrl = window.WOODPECKER_DOCS;
const url = toRef(props, 'url');
const docsUrl = computed(() => (url.value.startsWith('http') ? url.value : `${docsBaseUrl}${url.value}`));
</script>
