import { definePreset } from "@primeuix/themes";
import Aura from "@primeuix/themes/aura";
import { aqiColors } from "./tokens";

// Minimal theme override: point PrimeVue's `primary` ramp at the mockup's
// orange accent so themed components (buttons, focus rings, active states)
// match the palette without per-component CSS. Everything else stays Aura.
export const WisedPreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: "#fdf3ec",
      100: "#fadcc7",
      200: "#f5c39f",
      300: "#efa873",
      400: "#ea934f",
      500: aqiColors.elevated, // base — the shared orange token
      600: "#cf6f2e",
      700: "#b25c25",
      800: "#8f491d",
      900: "#6b3715",
      950: "#3f200c",
    },
  },
});
