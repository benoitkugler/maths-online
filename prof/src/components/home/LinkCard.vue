<template>
  <v-card
    :color="props.color"
    rounded
    link
    :to="props.href ? undefined : { name: props.routeName }"
    :href="props.href ? props.routeName : undefined"
    :target="props.href ? '_blank' : undefined"
  >
    <v-card-text class="pa-2">
      <v-icon>{{ props.icon }}</v-icon>
      {{ props.title }}
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
interface Props {
  title: string;
  routeName: string;
  href?: boolean;
  color: string;
  icon: string;
}

const props = defineProps<Props>();

// const emit = defineEmits<{
//   (e: "update:modelValue", modelValue: Date_): void;
// }>();
</script>
