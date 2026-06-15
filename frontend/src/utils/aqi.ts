import { aqiColors } from "../styles/tokens";

export type AqiBand = "good" | "moderate" | "elevated" | "unhealthy";

// Single source of truth for the AQI scale: thresholds + labels live here as
// data, so both bandFor (logic) and the legend (display) read the same table.
export const AQI_BANDS = [
  { band: "good", label: "Good", min: 0, max: 12 },
  { band: "moderate", label: "Moderate", min: 12, max: 35 },
  { band: "elevated", label: "Elevated", min: 35, max: 55 },
  { band: "unhealthy", label: "Unhealthy", min: 55, max: Infinity },
] as const;

// One lookup the rest derive from. The last band (max Infinity) always matches
function entryFor(value: number) {
  return AQI_BANDS.find((b) => value < b.max) ?? AQI_BANDS[AQI_BANDS.length - 1]!;
}

export function bandFor(value: number): AqiBand {
  return entryFor(value).band;
}

export function colorFor(value: number): string {
  return aqiColors[bandFor(value)];
}

export function labelFor(value: number): string {
  return entryFor(value).label;
}
