export const COLOR_THEMES = [
  { id: "vercel", name: "Vercel", color: "#000000", colors: ["#ffffff", "#000000", "#666666", "#999999"], isDark: false },
  { id: "vercel-dark", name: "Vercel Dark", color: "#000000", colors: ["#000000", "#ffffff", "#333333", "#666666"], isDark: true },
  { id: "violet-bloom", name: "Violet Bloom", color: "#7c3aed", colors: ["#7c3aed", "#8b5cf6", "#a78bfa", "#c4b5fd"], isDark: false },
  { id: "bubblegum", name: "Bubblegum", color: "#d946a8", colors: ["#d946a8", "#ec4899", "#f472b6", "#f9a8d4"], isDark: false },
  { id: "quantum-rose", name: "Quantum Rose", color: "#e11d48", colors: ["#e11d48", "#f43f5e", "#fb7185", "#fda4af"], isDark: false },
  { id: "clean-slate", name: "Clean Slate", color: "#3b82f6", colors: ["#3b82f6", "#60a5fa", "#93c5fd", "#bfdbfe"], isDark: false },
  { id: "cosmic-night", name: "Cosmic Night", color: "#6366f1", colors: ["#1e1b4b", "#6366f1", "#818cf8", "#a5b4fc"], isDark: true },
  { id: "cyberpunk-1", name: "Cyberpunk", color: "#00ffff", colors: ["#0f172a", "#00ffff", "#a855f7", "#ec4899"], isDark: true },
  { id: "eva-01", name: "EVA Unit-01", color: "#9333ea", colors: ["#1a0a2e", "#9333ea", "#22c55e", "#84cc16"], isDark: true },
] as const

export type ColorThemeId = typeof COLOR_THEMES[number]["id"]

export const DEFAULT_COLOR_THEME_ID: ColorThemeId = "vercel-dark"
export const COLOR_THEME_COOKIE_KEY = "color-theme"

const COLOR_THEME_IDS = new Set<string>(COLOR_THEMES.map((theme) => theme.id))

export function isColorThemeId(value: string | null | undefined): value is ColorThemeId {
  return typeof value === "string" && COLOR_THEME_IDS.has(value)
}

export function resolveColorThemeId(value: string | null | undefined): ColorThemeId {
  return isColorThemeId(value) ? value : DEFAULT_COLOR_THEME_ID
}

export function isDarkColorTheme(themeId: ColorThemeId): boolean {
  return COLOR_THEMES.find((theme) => theme.id === themeId)?.isDark ?? false
}

