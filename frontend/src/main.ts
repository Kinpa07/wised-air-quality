import { createApp } from "vue";
import { createPinia } from "pinia";
import PrimeVue from "primevue/config";

import "primeicons/primeicons.css";
import "leaflet/dist/leaflet.css";
import "./style.css";

import App from "./App.vue";
import { WisedPreset } from "./styles/preset";

const app = createApp(App);

app.use(createPinia());
app.use(PrimeVue, {
  theme: {
    preset: WisedPreset,
    options: {
      darkModeSelector: false,
    },
  },
});

app.mount("#app");
