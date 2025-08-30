<template>
  <v-avatar>
    <div :class="'bg-' + rankColors[props.rank]">
      <v-img
        v-if="props.rank"
        :src="beltImg"
        :width="props.small ? 30 : 60"
        :height="props.small ? 25 : 50"
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
  small?: boolean;
}

const props = defineProps<Props>();

//   const emit = defineEmits<{
//     (e: "click"): void;
//   }>();
</script>
