<template>
  <v-menu>
    <template v-slot:activator="{ isActive, props }">
      <v-chip
        v-on="{ isActive }"
        v-bind="props"
        class="mx-1"
        size="small"
        title="Selectionner la matière"
        @click.stop
        variant="elevated"
        :color="MatiereColor"
      >
        {{ sigle }}
      </v-chip>
    </template>
    <v-card density="compact">
      <v-card-text>
        <v-radio-group
          :model-value="props.matiere"
          @update:model-value="o => emit('update:matiere', o)"
          hide-details
          :color="MatiereColor"
        >
          <v-radio
            :value="matiere"
            v-for="(matiere, index) in Object.keys(MatiereTagLabels)"
            :key="index"
            :label="matiere"
          >
          </v-radio>
        </v-radio-group>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { MatiereTag } from "@/controller/api_gen";
import { MatiereTagLabels } from "@/controller/api_gen";
import { MatiereColor } from "@/controller/editor";
import { computed } from "vue";

interface Props {
  matiere: MatiereTag;
}
const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update:matiere", matiere: MatiereTag): void;
}>();

const sigle = computed(() => {
  switch (props.matiere) {
    case MatiereTag.Allemand:
      return "ALL";
    case MatiereTag.Anglais:
      return "ANG";
    case MatiereTag.Autre:
      return "AUT";
    case MatiereTag.Espagnol:
      return "ESP";
    case MatiereTag.Francais:
      return "FRA";
    case MatiereTag.HistoireGeo:
      return "H-G";
    case MatiereTag.Italien:
      return "ITA";
    case MatiereTag.Mathematiques:
      return "MAT";
    case MatiereTag.PhysiqueChimie:
      return "P-C";
    case MatiereTag.SES:
      return "SES";
    case MatiereTag.SVT:
      return "SVT";
    default:
      return "";
  }
});
</script>
