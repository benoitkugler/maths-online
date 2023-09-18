<template>
  <v-card title="Exceptions" :subtitle="props.sheet.Title">
    <v-dialog
      max-width="500px"
      :model-value="toEdit != null"
      @update:model-value="toEdit = null"
    >
      <v-card v-if="toEdit != null" title="Modifier les exceptions">
        <v-card-text class="my-2">
          <v-form>
            <v-row>
              <v-col>
                <v-select
                  label="Élève"
                  density="compact"
                  variant="outlined"
                  color="primary"
                  :items="studentItems"
                  :model-value="toEdit.IdStudent < 0 ? null : toEdit.IdStudent"
                  @update:model-value="v => (toEdit!.IdStudent = v)"
                  hide-details
                ></v-select>
              </v-col>
            </v-row>
            <v-row>
              <v-col align-self="center">
                <v-switch
                  label="Clôture personnalisée"
                  hide-details
                  v-model="toEdit.Deadline.Valid"
                  color="primary"
                >
                </v-switch>
              </v-col>
              <v-col align-self="center">
                <DateTimeChip
                  v-if="toEdit.Deadline.Valid"
                  title="Modifier la clôture"
                  :model-value="toEdit.Deadline.Time"
                  @update:model-value="
                    v => {
                      toEdit!.Deadline.Time = v;
                    }
                  "
                ></DateTimeChip>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-checkbox
                  label="Ne pas compter la note dans la moyenne"
                  color="primary"
                  v-model="toEdit.IgnoreForMark"
                  hide-details
                ></v-checkbox>
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="success"
            :disabled="toEdit.IdStudent < 0"
            @click="setDispense(toEdit)"
            >Enregistrer</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!--  -->

    <template v-slot:append>
      <v-btn
        @click="
          toEdit = {
            IdTravail: props.travail.Id,
            IdStudent: -1,
            Deadline: emptyTime(),
            IgnoreForMark: false
          }
        "
      >
        <template v-slot:prepend>
          <v-icon color="green">mdi-plus</v-icon>
        </template>
        Ajouter
      </v-btn>
    </template>
    <v-card-text>
      <v-list>
        <v-list-item>
          <v-row class="text-center">
            <v-col cols="5">Élève</v-col>
            <v-col cols="4">Clôture</v-col>
            <v-col cols="3">Note ignorée</v-col>
          </v-row>
        </v-list-item>
        <v-list-item
          v-if="!data.Exceptions?.length"
          class="text-center"
          subtitle="Aucune permission n'est enregistrée."
        >
        </v-list-item>
        <v-list-item
          v-for="(exp, index) in data.Exceptions"
          :key="index"
          @click="toEdit = copy(exp)"
        >
          <template v-slot:prepend>
            <v-btn
              icon
              size="x-small"
              class="mr-2"
              @click.stop="removeDispense(exp)"
            >
              <v-icon color="red">mdi-close</v-icon>
            </v-btn>
          </template>
          <v-row>
            <v-col cols="5">
              <v-list-item-title>
                {{ formatName((data.Students || {})[exp.IdStudent]) }}
              </v-list-item-title>
            </v-col>
            <v-col cols="4" class="text-center">
              <span v-if="exp.Deadline.Valid">{{
                formatTime(exp.Deadline.Time, true)
              }}</span>
              <small v-else>-</small>
            </v-col>
            <v-col cols="3" class="text-center">
              <v-icon v-if="exp.IgnoreForMark">mdi-check</v-icon>
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Travail,
  Exceptions,
  Student,
  TravailException,
  Time,
  Sheet,
  NullTime
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { ref } from "vue";
import DateTimeChip from "../DateTimeChip.vue";
import { computed } from "vue";
import { copy, formatTime } from "@/controller/utils";

interface Props {
  travail: Travail;
  sheet: Sheet;
}

const props = defineProps<Props>();

// const emit = defineEmits<{}>();

const data = ref<Exceptions>({ Exceptions: [], Students: {} });

const toEdit = ref<TravailException | null>(null);

onMounted(fetchDispenses);

const studentItems = computed(() => {
  const out = Object.values(data.value.Students || {});
  out.sort((a, b) => a.Name.localeCompare(b.Name));
  return out.map(st => ({ value: st.Id, title: formatName(st) }));
});

async function fetchDispenses() {
  const res = await controller.HomeworkGetDispenses({
    "id-travail": props.travail.Id
  });
  if (res === undefined) return;
  data.value = res;
}

function formatName(student: Student) {
  return `${student.Name} ${student.Surname}`;
}

async function setDispense(params: TravailException) {
  toEdit.value = null;
  const res = await controller.HomeworkSetDispense(params);
  if (res === undefined) return;
  fetchDispenses();
}

async function removeDispense(params: TravailException) {
  params.Deadline = emptyTime();
  params.IgnoreForMark = false;
  setDispense(params);
}

function emptyTime(): NullTime {
  return {
    Valid: false,
    Time: new Date(Date.now()).toISOString() as Time
  };
}
</script>
