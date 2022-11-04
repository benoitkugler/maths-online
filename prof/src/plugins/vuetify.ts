// Styles
import "@mdi/font/css/materialdesignicons.css";
// Vuetify
import { createVuetify } from "vuetify";
import { aliases, mdi } from "vuetify/lib/iconsets/mdi";
import "vuetify/styles";

const myCustomLightTheme: ThemeDefinition = {
  dark: false,
  colors: {
    primary: "#6fdec1",
    "primary-darken-1": "#52a38e",
    secondary: "#e8f241",
    "secondary-darken-1": "#e5f02e"
  }
};

export default createVuetify({
  theme: {
    defaultTheme: "myCustomLightTheme",
    themes: {
      myCustomLightTheme
    }
  },
  icons: {
    defaultSet: "mdi",
    aliases,
    sets: {
      mdi
    }
  }
});
