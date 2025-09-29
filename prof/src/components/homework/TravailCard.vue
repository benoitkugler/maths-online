<template>
  <v-card density="compact">
    <v-card-text>
      <v-row justify="space-between">
        <v-col align-self="center" cols="auto">
          <v-menu v-if="sheet.Origin.Visibility == Visibility.Admin">
            <template v-slot:activator="{ isActive, props }">
              <v-card
                v-on="{ isActive }"
                v-bind="props"
                variant="elevated"
                :title="sheet.Sheet.Title"
                :subtitle="subtitle"
                :color="color"
              >
              </v-card>
            </template>
            <PreviewSheet :sheet="sheet"></PreviewSheet>
          </v-menu>
          <v-card
            v-else
            variant="elevated"
            @click="emit('editSheet', props.sheet)"
            :title="sheet.Sheet.Title"
            :subtitle="subtitle"
            :color="color"
          >
          </v-card>
        </v-col>
        <v-col cols="auto" align-self="center">
          <v-menu>
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                density="comfortable"
                icon
                size="small"
                class="mr-1"
                v-bind="props"
                v-on="{ isActive }"
              >
                <v-icon color="secondary"> mdi-heart </v-icon>
              </v-btn>
            </template>
            <v-card v-if="sheet.Sheet.Anonymous.Valid">
              <v-card-text> Feuille anonyme </v-card-text>
              <v-card-actions>
                <v-spacer> </v-spacer>
                <v-btn @click="emit('setFavorite', sheet.Sheet)"
                  >Enregistrer dans les favoris</v-btn
                >
              </v-card-actions>
            </v-card>
            <v-card v-else :color="colorForOrigin(sheet.Origin)">
              <v-card-text v-if="sheet.Origin.Visibility == Visibility.Admin">
                Feuille favorite de la base officielle
              </v-card-text>
              <v-card-text v-else> Feuille favorite personelle </v-card-text>
            </v-card>
          </v-menu>

          <v-tooltip
            text="Modifier les exceptions..."
            location="top"
            v-if="inner.Noted"
          >
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                density="comfortable"
                icon
                size="small"
                class="mr-1"
                v-bind="props"
                v-on="{ isActive }"
                @click="emit('showDispenses')"
              >
                <v-icon> mdi-account-supervisor </v-icon>
              </v-btn>
            </template>
          </v-tooltip>
          <v-menu offset-y close-on-content-click>
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                v-on="{ isActive }"
                v-bind="props"
                density="comfortable"
                icon
                size="small"
                title="Copier vers une autre classe..."
                class="mr-1"
              >
                <v-icon icon="mdi-content-copy" size="small"></v-icon>
              </v-btn>
            </template>
            <v-list>
              <v-list-subheader>Copier vers...</v-list-subheader>
              <v-list-item
                v-for="(classroom, index) in classrooms"
                :key="index"
                link
                @click="emit('copy', classroom.id)"
              >
                {{ classroom.name }}
              </v-list-item>
            </v-list>
          </v-menu>

          <v-btn
            density="comfortable"
            @click="emit('delete')"
            icon
            title="Supprimer le travail"
            size="small"
          >
            <v-icon color="red" icon="mdi-delete" size="small"></v-icon>
          </v-btn>
        </v-col>
      </v-row>

      <v-row>
        <v-col align-self="center">
          <span class="text-grey text-subtitle-1 ml-1">
            Afficher à partir du
          </span>
        </v-col>
        <v-col cols="auto" align-self="center">
          <DateTimeChip
            title="Modifier le début du travail"
            v-model="inner.ShowAfter"
            @update:model-value="emit('update', inner)"
          ></DateTimeChip>
        </v-col>
      </v-row>
      <v-row justify="space-between" class="mt-2" no-gutters>
        <v-col cols="auto" align-self="center">
          <v-switch
            label="Limiter l'accès"
            v-model="inner.Noted"
            @update:model-value="emit('update', inner)"
            hide-details
            color="primary"
          >
          </v-switch>
        </v-col>

        <v-col cols="auto" align-self="center">
          <DateTimeChip
            prefix="clôture le"
            title="Modifier la clôture du travail"
            v-if="inner.Noted"
            v-model="inner.Deadline"
            @update:model-value="emit('update', inner)"
            :min-date="inner.ShowAfter"
          ></DateTimeChip>
          <v-chip v-else> Feuille en accès libre, sans clôture</v-chip>
        </v-col>
      </v-row>
      <v-row v-if="inner.Noted">
        <v-col>
          <v-select
            variant="outlined"
            density="compact"
            hide-details
            :items="selectItems(QuestionRepeatLabels)"
            label="Nombre d'essais par question"
            v-model="inner.QuestionRepeat"
            @update:model-value="emit('update', inner)"
          ></v-select>
        </v-col>
        <v-col align-self="center">
          <v-menu :close-on-content-click="false" v-model="showEditTimeLimit">
            <template #activator="{ props: menuProps }">
              <v-chip
                variant="outlined"
                label
                v-bind="menuProps"
                @click="editedTimeLimit = inner.QuestionTimeLimit"
                >Limite de temps par question :
                {{
                  inner.QuestionTimeLimit > 0
                    ? `${inner.QuestionTimeLimit} secs.`
                    : "Aucune"
                }}</v-chip
              >
            </template>
            <v-card>
              <v-card-text>
                <v-combobox
                  variant="outlined"
                  density="compact"
                  label="Limite"
                  v-model.number="editedTimeLimit"
                  suffix="secondes"
                  hint="0 pour ne pas limiter."
                  persistent-hint
                  :items="[5, 10, 15, 20, 30, 60, 0]"
                ></v-combobox>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  @click="
                    showEditTimeLimit = false;
                    inner.QuestionTimeLimit = editedTimeLimit;
                    emit('update', inner);
                  "
                  :disabled="isNaN(Number(editedTimeLimit))"
                  >Enregistrer</v-btn
                >
              </v-card-actions>
            </v-card>
          </v-menu>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  Visibility,
  type Classroom,
  type Sheet,
  type SheetExt,
  type Travail,
  IdClassroom,
  QuestionRepeatLabels,
  Int,
} from "@/controller/api_gen";
import { computed, ref, watch } from "vue";
import PreviewSheet from "./PreviewSheet.vue";
import DateTimeChip from "../DateTimeChip.vue";
import { colorForOrigin, copy, selectItems } from "@/controller/utils";

interface Props {
  travail: Travail;
  sheet: SheetExt;
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "update", tr: Travail): void;
  (e: "copy", idTarget: IdClassroom): void;
  (e: "setFavorite", sheet: Sheet): void;
  (e: "editSheet", sheet: SheetExt): void;
  (e: "showDispenses"): void;
}>();

const inner = ref(copy(props.travail));

watch(
  () => props.travail,
  () => (inner.value = copy(props.travail))
);

const nbTasks = computed(() => props.sheet.Tasks?.length || 0);

const subtitle = computed(() => {
  if (nbTasks.value == 0) {
    return "Aucune tâche";
  } else if (nbTasks.value == 1) {
    return "1 tâche";
  } else {
    return `${nbTasks.value} tâches`;
  }
});

const color = computed(() => {
  const baseColor = "grey-lighten-3";
  if (!inner.value.Noted) return baseColor;
  const now = new Date(Date.now());
  const start = new Date(inner.value.ShowAfter);
  const end = new Date(inner.value.Deadline);
  return start <= now && now <= end ? "blue-lighten-2" : baseColor;
});

const editedTimeLimit = ref(0 as Int);
const showEditTimeLimit = ref(false);
</script>
