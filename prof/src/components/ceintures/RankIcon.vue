<template>
  <v-avatar>
    <div :class="'bg-' + rankColors[props.rank]">
      <v-img
        :src="beltImg"
        width="60"
        height="50"
        :style="{
          filter: props.rank == Rank.Blanche ? 'invert(30%)' : undefined,
        }"
      >
      </v-img>
    </div>
  </v-avatar>
</template>

<script setup lang="ts">
import { Rank } from "@/controller/api_gen";
import beltImg from "@/assets/belt-mask.png";
import { rankColors } from "@/controller/utils";

interface Props {
  rank: Rank;
}

const props = defineProps<Props>();

//   const emit = defineEmits<{
//     (e: "click"): void;
//   }>();
</script>
