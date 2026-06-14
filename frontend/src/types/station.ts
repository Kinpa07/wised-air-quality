export type StationBand = "Good" | "Moderate" | "Elevated" | "Unhealthy";
export type ConnectionQuality = "Good" | "Poor";
export type Pollutant = "pm2_5" | "pm10";

export interface Station {
  id: string;
  lat: number;
  lng: number;
  pm2_5: number | null;
  pm10: number | null;
  measured_at: string | null;
  band: StationBand | null;
  district: string;
  stability: number;
  connection: ConnectionQuality;
}
