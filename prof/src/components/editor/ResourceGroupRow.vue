<template>
  <v-list-item
    link
    @click="emit('clicked')"
    rounded
    :class="'px-1 bg-' + colorForOrigin(props.group.Origin)"
  >
    <v-row no-gutters justify="space-between">
      <v-col cols="auto" align-self="center">
        <v-btn
          v-if="props.group.Origin.Visibility == Visibility.Admin"
          class="mx-2 my-1"
          size="x-small"
          icon
          @click.stop="emit('duplicate')"
          :title="
            props.isQuestion
              ? 'Dupliquer et importer cette question'
              : 'Dupliquer et importer cet exercice'
          "
        >
          <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
        </v-btn>
        <v-menu v-else>
          <template v-slot:activator="{ isActive, props }">
            <v-btn
              v-on="{ isActive }"
              v-bind="props"
              class="mx-2 my-1"
              size="x-small"
              icon="mdi-dots-vertical"
            ></v-btn>
          </template>

          <v-list>
            <v-list-item>
              <v-btn
                class="my-1"
                size="small"
                @click.stop="emit('duplicate')"
                flat
              >
                <template v-slot:prepend>
                  <v-icon
                    class="mr-4"
                    icon="mdi-content-copy"
                    color="secondary"
                  ></v-icon>
                </template>
                Dupliquer
              </v-btn>
            </v-list-item>

            <v-list-item>
              <OriginButton
                :origin="props.group.Origin"
                @update-public="b => emit('updatePublic', b)"
                @create-review="emit('createReview')"
              ></OriginButton>
            </v-list-item>

            <v-list-item>
              <v-btn class="my-1" size="small" flat @click="emit('delete')">
                <template v-slot:prepend>
                  <v-icon class="mr-4" icon="mdi-delete" color="red"></v-icon>
                </template>
                Supprimer
              </v-btn>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>

      <v-col class="my-3 mx-1" style="text-align: left" align-self="center">
        {{ props.group.Title }}

        <VariantList :variants="variants"></VariantList>
      </v-col>
      <v-col cols="auto" align-self="center">
        <TagIndex :tags="props.group.Tags || []"></TagIndex>
      </v-col>
    </v-row>
  </v-list-item>
</template>

<script setup lang="ts">
import type { TagsDB } from "@/controller/api_gen";
import { Visibility } from "@/controller/api_gen";
import {
  colorForOrigin,
  type ResourceGroup,
  type VariantG
} from "@/controller/editor";
import { computed } from "vue";
import OriginButton from "../OriginButton.vue";
import VariantList from "./VariantList.vue";
import TagIndex from "../TagIndex.vue";

interface Props {
  group: ResourceGroup;
  allTags: TagsDB;
  isQuestion: boolean;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "clicked"): void;
  (e: "duplicate"): void;
  (e: "delete"): void;
  (e: "updatePublic", isPublic: boolean): void;
  (e: "createReview"): void;
}>();

const variants = computed<VariantG[]>(() => {
  return props.group.Variants || [];
});
</script>

<style scoped></style>
