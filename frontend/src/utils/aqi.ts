import { aqiColors } from "../styles/tokens";

export type AqiBand = "good" | "moderate" | "elevated" | "unhealthy";

export function bandFor(value: number): AqiBand {
  if (value < 12) return "good";
  if (value < 35) return "moderate";
  if (value < 55) return "elevated";
  return "unhealthy";
}

export function colorFor(value: number): string {
  return aqiColors[bandFor(value)];
}
