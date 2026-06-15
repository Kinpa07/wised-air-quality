// The four AQI band colours, matching the legend (Good 0–12, Moderate 12–35,
// Elevated 35–55, Unhealthy 55+). Warm palette read off the mockup's badges,
// markers, and chart line.
export const aqiColors = {
  good: "#3FA776", // green
  moderate: "#E5B53C", // amber
  elevated: "#E5823A", // orange (also the chart trend line)
  unhealthy: "#D6483B", // red
} as const;

// Surface / background / text colours.
export const surface = {
  bg: "#F1F0EE", // app background — warm light gray
  card: "#FFFFFF", // KPI tiles, panels, watchlist
  border: "#E5E3DF", // hairline borders between sections
  text: "#1A1A1A", // primary text / KPI numbers
  muted: "#8C8A86", // labels, secondary text
} as const;

// Fallback marker fill for stations with no current reading (no AQI band).
export const markerNoData = surface.muted;

// Chart-specific fills (PM10 bars are a neutral tan; PM2.5 line reuses
// aqiColors.elevated).
export const chart = {
  bar: "#E7DECF", // PM10 bars — muted tan
  line: aqiColors.elevated, // PM2.5 trend line — orange
} as const;

// Spacing scale (px) read off the mockup's gaps and paddings.
export const space = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
} as const;

// Corner radii.
export const radius = {
  sm: 4,
  md: 8,
  lg: 12,
} as const;

// Fixed pixel sizes for elements that need a consistent width/height.
export const sizes = {
  badge: 96, // uniform AQI pill width; fits the longest label ("Unhealthy")
  map: 460, // Leaflet needs an explicit pixel height
  chart: 340, // trend chart body height (Chart.js needs an explicit height)
  search: 280,
  maxWidth: 1600, // cap + centre the dashboard so it doesn't stretch on 4K
  titleRow: 40, // so a card header with a toggle matches a text-only one
} as const;
