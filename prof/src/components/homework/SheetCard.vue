<template>
  <v-list-item
    :class="'py-2 mb-1 px-2 bg-' + colorForOrigin(props.sheet.Origin)"
    style="break-inside: avoid-column"
    color="grey"
    rounded
  >
    <v-menu>
      <template v-slot:activator="{ isActive, props }">
        <v-card v-on="{ isActive }" v-bind="props" class="pa-2 mr-2">
          <v-list-item-title>
            {{ sheet.Sheet.Title }}
          </v-list-item-title>
          <v-list-item-subtitle>
            {{ subtitle }}
          </v-list-item-subtitle>
        </v-card>
      </template>
      <PreviewSheet :sheet="sheet"></PreviewSheet>
    </v-menu>
    <template v-slot:append="{}">
      <v-list-item-action>
        <v-menu offset-y close-on-content-click>
          <template v-slot:activator="{ isActive, props: props2 }">
            <v-btn
              v-on="{ isActive }"
              v-bind="props2"
              variant="outlined"
              color="primary-darken-1"
              title="Ajouter à une classe"
            >
              Assigner
            </v-btn>
          </template>
          <v-list>
            <v-list-subheader v-if="!props.classrooms.length">
              Aucune classe.
            </v-list-subheader>
            <v-list-subheader v-else
              >Assigner cette feuille à...</v-list-subheader
            >
            <v-list-item
              v-for="(classroom, index) in props.classrooms"
              :key="index"
              link
              @click="emit('assign', classroom.id)"
            >
              {{ classroom.name }}
            </v-list-item>
          </v-list>
        </v-menu>
        <!--  -->
        <v-btn
          v-if="sheet.Origin.Visibility == Visibility.Admin"
          @click="emit('duplicate')"
          title="Dupliquer et importer"
          icon
          class="ml-4"
          size="x-small"
        >
          <v-icon color="secondary"> mdi-content-copy </v-icon>
        </v-btn>

        <v-menu v-else offset-y close-on-content-click>
          <template v-slot:activator="{ isActive, props }">
            <v-btn
              v-on="{ isActive }"
              v-bind="props"
              icon
              class="ml-4"
              size="x-small"
            >
              <v-icon>mdi-dots-vertical</v-icon>
            </v-btn>
          </template>
          <v-list density="compact">
            <v-list-item>
              <v-btn
                variant="flat"
                class="mr-2"
                size="small"
                @click="emit('edit')"
              >
                <template v-slot:prepend>
                  <v-icon icon="mdi-pencil" class="mr-4"></v-icon>
                </template>
                Editer
              </v-btn>
            </v-list-item>
            <v-list-item>
              <v-btn
                variant="flat"
                class="mr-2"
                size="small"
                @click="emit('duplicate')"
              >
                <template v-slot:prepend>
                  <v-icon color="secondary" class="mr-4">
                    mdi-content-copy
                  </v-icon>
                </template>
                Dupliquer
              </v-btn>
            </v-list-item>
            <v-list-item>
              <v-btn
                variant="flat"
                class="mr-2"
                size="small"
                @click="emit('delete')"
              >
                <template v-slot:prepend>
                  <v-icon color="red" class="mr-4">mdi-delete</v-icon>
                </template>
                Supprimer
              </v-btn>
            </v-list-item>

            <v-list-item>
              <OriginButton
                :origin="sheet.Origin"
                @update-public="(b) => emit('updatePublic', b)"
                @create-review="emit('createReview')"
              ></OriginButton>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-list-item-action>
    </template>
  </v-list-item>
</template>

<script setup lang="ts">
import {
  Visibility,
  type Classroom,
  type SheetExt,
} from "@/controller/api_gen";
import { computed } from "vue";
import OriginButton from "../OriginButton.vue";
import PreviewSheet from "./PreviewSheet.vue";
import { colorForOrigin } from "@/controller/utils";

interface Props {
  sheet: SheetExt;
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "assign", idClassroom: number): void;
  (e: "edit"): void;
  (e: "duplicate"): void;
  (e: "delete"): void;
  (e: "updatePublic", pub: boolean): void;
  (e: "createReview"): void;
}>();

const subtitle = computed(() => {
  if (props.sheet.NbTravaux == 0) {
    return "Inactive";
  } else if (props.sheet.NbTravaux == 1) {
    return "Active dans 1 travail";
  } else {
    return `Active dans ${props.sheet.NbTravaux} travaux`;
  }
});
</script>
