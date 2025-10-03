import { Plus, Minus, Dot } from "@lucide/svelte"

export const iconByState = {
  added: Plus,
  removed: Minus,
  normal: Dot,
} as const

export function formatTimestamp(time: Date) {
  return new Date(time).toLocaleTimeString()
}

// Values in Mbps
const SPEED_MAP: Record<string, number> = {
  low: 1.5,
  full: 12,
  high: 480,
  super: 5000,
}

export function formatSpeed(rawSpeed?: string): string {
  if (!rawSpeed) {
    return "unknown"
  }

  const normalized = rawSpeed.toLowerCase()
  const value = SPEED_MAP[normalized]

  if (value === undefined) {
    return rawSpeed
  }

  if (value >= 1000) {
    const gbps = value / 1000
    return `${gbps} Gbps`
  }

  return `${value} Mbps`
}
