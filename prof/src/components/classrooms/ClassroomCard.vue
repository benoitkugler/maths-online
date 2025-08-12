<template>
  <v-card :title="props.classroom.Classroom.name" class="ma-2">
    <template #append>
      <v-btn icon>
        <v-badge
          :content="props.classroom.SharedWith?.length"
          floating
          :color="props.classroom.SharedWith?.length ? 'pink' : 'transparent'"
        >
          <v-icon>mdi-account-group</v-icon>
        </v-badge>

        <v-menu activator="parent" :close-on-content-click="false">
          <v-card
            title="Partager la classe"
            subtitle="Autres enseignants liés à cette classe"
            min-width="600px"
          >
            <v-card-text class="text-center">
              <i v-if="!props.classroom.SharedWith?.length">
                La classe n'est pas partagée.
              </i>
              <v-chip
                v-for="(teacherMail, index) in props.classroom.SharedWith"
                :key="index"
                >{{ teacherMail }}</v-chip
              >
              <v-divider thickness="2" class="my-4"></v-divider>
              <v-row>
                <v-col align-self="center" cols="7">
                  <v-text-field
                    density="compact"
                    variant="outlined"
                    label="Adresse mail du compte"
                    v-model="mailToInvite"
                    hide-details
                  >
                  </v-text-field>
                </v-col>
                <v-col align-self="center" cols="5">
                  <v-btn
                    color="green"
                    :disabled="
                      !mailToInvite.length ||
                      props.classroom.SharedWith?.includes(mailToInvite)
                    "
                    @click="emit('invite', mailToInvite)"
                    >Ajouter un collègue</v-btn
                  >
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-menu>
      </v-btn>
    </template>
    <v-card-text>
      <v-row justify="center" class="mt-1">
        <v-col cols="auto">
          <v-chip
            link
            variant="elevated"
            color="primary"
            @click="emit('showStudents')"
          >
            {{ eleveText }}
          </v-chip>
        </v-col>
      </v-row>
    </v-card-text>

    <v-card-actions>
      <v-btn icon color="red" title="Supprimer" @click.stop="emit('delete')">
        <v-icon icon="mdi-delete"></v-icon>
      </v-btn>
      <v-spacer></v-spacer>
      <v-btn @click.stop="emit('update')">
        <template v-slot:append>
          <v-icon>mdi-cog</v-icon>
        </template>
        Modifier
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { ClassroomExt } from "@/controller/api_gen";
import { computed, ref } from "vue";

interface Props {
  classroom: ClassroomExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "showStudents"): void;
  (e: "delete"): void;
  (e: "update"): void;
  (e: "invite", mail: string): void;
}>();

const eleveText = computed(() => {
  switch (props.classroom.NbStudents) {
    case 0:
      return "Ajouter des élèves...";
    case 1:
      return "Un élève";
    default:
      return `${props.classroom.NbStudents} élèves`;
  }
});

const mailToInvite = ref("");
</script>
