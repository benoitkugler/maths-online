<template>
  <span class="bg-grey-lighten-3 px-1 rounded editor-font">
    <slot></slot>
  </span>
</template>

<script setup lang="ts"></script>

<style scoped>
.editor-font {
  font-family: "Roboto Mono", monospace;
}
</style>
