/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import "@mdi/font/css/materialdesignicons.css";
import "vuetify/styles";

// Composables
import { createVuetify } from "vuetify";

// Translations provided by Vuetify
import { fr } from "vuetify/locale";

const myCustomLightTheme = {
  dark: false,
  colors: {
    primary: "#6fdec1",
    "primary-darken-1": "#52a38e",
    secondary: "#e8f241",
    "secondary-darken-1": "#b3bd0d",
    "secondary-lighten-1": "#f4fa91",
  },
};

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: "myCustomLightTheme",
    themes: {
      myCustomLightTheme,
    },
  },
  locale: {
    locale: "fr",
    messages: { fr },
  },
});
