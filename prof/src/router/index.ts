import ClassroomsListVue from "@/views/ClassroomsList.vue";
import ExerciceEditorVue from "@/views/ExerciceEditor.vue";
import HomeworkActivityVue from "@/views/HomeworkActivity.vue";
import QuestionEditor from "@/views/QuestionEditor.vue";
import TrivialPoursuit from "@/views/TrivialPoursuit.vue";
import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";

const router = createRouter({
  history: createWebHistory("/prof/"),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
      meta: { Label: "Accueil" },
    },
    {
      path: "/classrooms",
      name: "classrooms",
      component: ClassroomsListVue,
      meta: { Label: "Classes et élèves" },
    },
    {
      path: "/editor-question",
      name: "editor-question",
      component: QuestionEditor,
      meta: { Label: "Editeur de question" },
    },
    {
      path: "/editor-exercice",
      name: "editor-exercice",
      component: ExerciceEditorVue,
      meta: { Label: "Editeur d'exercices" },
    },
    {
      path: "/trivial",
      name: "trivial",
      component: TrivialPoursuit,
      meta: { Label: "Triv'Maths" },
    },
    {
      path: "/homework",
      name: "homework",
      component: HomeworkActivityVue,
      meta: { Label: "Travail à la maison" },
    },
    // {
    //   path: '/about',
    //   name: 'Activités',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/AboutView.vue')
    // },
    {
      path: "/:catchAll(.*)",
      redirect: { name: "home" },
    },
  ],
});

export default router;
