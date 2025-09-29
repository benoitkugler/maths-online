<template>
  <v-card
    class="my-5 mx-auto"
    width="80%"
    max-width="600px"
    title="Ré-initialisation du mot de passe"
  >
    <v-card-text class="ma-2">
      <v-form>
        <v-row>
          <v-col>
            <v-text-field
              label="Nouveau mot de passe"
              v-model="password"
              density="compact"
              variant="outlined"
            ></v-text-field>
          </v-col>
        </v-row>
      </v-form>
      <v-card-actions>
        <v-spacer> </v-spacer>
        <v-btn :disabled="password.length < 4" @click="resetPassword"
          >Changer mon mot de passe</v-btn
        >
      </v-card-actions>
    </v-card-text>

    <v-dialog v-model="showSuccess" max-width="600px">
      <v-card title="Mot de passe modifié avec succès">
        <v-card-text>
          Vous allez être redirigé vers la page de connection.
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script setup lang="ts">
import { controller } from "@/controller/controller";
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

onMounted(handleQuery);

const seal = ref("");

async function handleQuery() {
  await router.isReady();
  seal.value = (router.currentRoute.value.query["seal"] || "") as string;
}

const password = ref("");

const showSuccess = ref(false);
async function resetPassword() {
  const res = controller.TeacherResetPassword({
    seal: seal.value,
    password: password.value,
  });
  if (res === undefined) return;
  showSuccess.value = true;
  setTimeout(() => {
    showSuccess.value = false;
    window.location.replace("/prof");
  }, 3000);
}
</script>

<style></style>
