/**
 * 颜色主题切换 hook
 * 管理主题色（不是亮暗模式）
 */
import { useEffect, useState, useCallback } from "react"
import { useTheme } from "next-themes"
import {
  COLOR_THEMES,
  COLOR_THEME_COOKIE_KEY,
  DEFAULT_COLOR_THEME_ID,
  isColorThemeId,
  type ColorThemeId,
  isDarkColorTheme,
} from "@/lib/color-themes"

export { COLOR_THEMES }
export type { ColorThemeId }

const STORAGE_KEY = COLOR_THEME_COOKIE_KEY
const COOKIE_MAX_AGE_SECONDS = 60 * 60 * 24 * 365 * 2

// Cache for localStorage reads to avoid expensive I/O operations
let themeCache: ColorThemeId | null = null

function getCookieTheme(): ColorThemeId | null {
  if (typeof document === "undefined") return null
  const match = document.cookie.match(new RegExp(`(?:^|; )${STORAGE_KEY}=([^;]*)`))
  if (!match) return null
  const value = decodeURIComponent(match[1])
  return isColorThemeId(value) ? value : null
}

function setThemeCookie(themeId: ColorThemeId) {
  if (typeof document === "undefined") return
  try {
    document.cookie = `${STORAGE_KEY}=${encodeURIComponent(themeId)}; Path=/; Max-Age=${COOKIE_MAX_AGE_SECONDS}; SameSite=Lax`
  } catch {
    // Ignore cookie write failures (e.g. blocked)
  }
}

/**
 * 获取当前颜色主题（带缓存和错误处理）
 */
function getStoredTheme(): ColorThemeId {
  if (typeof window === "undefined") return DEFAULT_COLOR_THEME_ID

  // Check cache first
  if (themeCache !== null) {
    return themeCache
  }

  const cookieTheme = getCookieTheme()
  if (cookieTheme) {
    themeCache = cookieTheme
    return cookieTheme
  }

  // Read from localStorage with error handling
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    const theme = isColorThemeId(stored) ? stored : DEFAULT_COLOR_THEME_ID
    themeCache = theme
    return theme
  } catch {
    // localStorage throws in incognito mode or when disabled
    return DEFAULT_COLOR_THEME_ID
  }
}

/**
 * 应用颜色主题到 DOM（仅设置 data-theme）
 */
function applyThemeAttribute(themeId: ColorThemeId) {
  if (typeof window === "undefined") return
  const root = document.documentElement
  root.setAttribute("data-theme", themeId)
  if (isDarkColorTheme(themeId)) {
    root.classList.add("dark")
  } else {
    root.classList.remove("dark")
  }
}

/**
 * 颜色主题 hook
 */
export function useColorTheme() {
  const [theme, setThemeState] = useState<ColorThemeId>(DEFAULT_COLOR_THEME_ID)
  const [mounted, setMounted] = useState(false)
  const { setTheme: setNextTheme } = useTheme()

  // 初始化
  useEffect(() => {
    const stored = getStoredTheme()
    setThemeState(stored)
    applyThemeAttribute(stored)
    // 同步 next-themes 亮暗模式
    const themeConfig = COLOR_THEMES.find(t => t.id === stored)
    setNextTheme(themeConfig?.isDark ? 'dark' : 'light')
    setMounted(true)
  }, [setNextTheme])

  // 切换主题
  const setTheme = useCallback((newTheme: ColorThemeId) => {
    setThemeState(newTheme)
    // Save to localStorage with error handling
    try {
      localStorage.setItem(STORAGE_KEY, newTheme)
    } catch {
      // localStorage throws when quota exceeded or disabled
    }
    // Update cache even if storage fails
    themeCache = newTheme
    setThemeCookie(newTheme)
    applyThemeAttribute(newTheme)
    // 同步 next-themes 亮暗模式
    const themeConfig = COLOR_THEMES.find(t => t.id === newTheme)
    setNextTheme(themeConfig?.isDark ? 'dark' : 'light')
  }, [setNextTheme])

  // 获取当前主题信息
  const currentTheme = COLOR_THEMES.find(t => t.id === theme) || COLOR_THEMES[0]

  return {
    theme,
    setTheme,
    themes: COLOR_THEMES,
    currentTheme,
    mounted,
  }
}
