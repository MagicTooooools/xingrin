import type React from "react"
import type { Metadata } from "next"
import { NextIntlClientProvider } from 'next-intl'
import { getMessages, setRequestLocale, getTranslations } from 'next-intl/server'
import { notFound } from 'next/navigation'
import { cookies } from "next/headers"
import { locales, localeHtmlLang, type Locale } from '@/i18n/config'
import { COLOR_THEME_COOKIE_KEY, isColorThemeId, DEFAULT_COLOR_THEME_ID, isDarkColorTheme } from "@/lib/color-themes"

// Import global style files
import "../globals.css"
// Import Noto Sans SC local font
import "@fontsource/noto-sans-sc/400.css"
import "@fontsource/noto-sans-sc/500.css"
import "@fontsource/noto-sans-sc/700.css"
// Import color themes
import "@/styles/themes/bubblegum.css"
import "@/styles/themes/quantum-rose.css"
import "@/styles/themes/clean-slate.css"
import "@/styles/themes/cosmic-night.css"
import "@/styles/themes/vercel.css"
import "@/styles/themes/vercel-dark.css"
import "@/styles/themes/violet-bloom.css"
import "@/styles/themes/cyberpunk-1.css"
import { Suspense } from "react"
import Script from "next/script"
import { QueryProvider } from "@/components/providers/query-provider"
import { ThemeProvider } from "@/components/providers/theme-provider"
import { UiI18nProvider } from "@/components/providers/ui-i18n-provider"
import { ColorThemeInit } from "@/components/color-theme-init"

// Import common layout components
import { RoutePrefetch } from "@/components/route-prefetch"
import { RouteProgress } from "@/components/route-progress"
import { AuthLayout } from "@/components/auth/auth-layout"

// Dynamically generate metadata
export async function generateMetadata({ params }: { params: Promise<{ locale: string }> }): Promise<Metadata> {
  const { locale } = await params
  const t = await getTranslations({ locale, namespace: 'metadata' })
  
  return {
    title: t('title'),
    description: t('description'),
    keywords: t('keywords').split(',').map(k => k.trim()),
    generator: "LunaFox ASM Platform",
    authors: [{ name: "yyhuni" }],
    icons: {
      icon: [
        { url: "/images/icon-64.png", sizes: "64x64", type: "image/png" },
        { url: "/images/icon-256.png", sizes: "256x256", type: "image/png" },
      ],
      apple: [{ url: "/images/icon-256.png", sizes: "256x256", type: "image/png" }],
    },
    openGraph: {
      title: t('ogTitle'),
      description: t('ogDescription'),
      type: "website",
      locale: locale === 'zh' ? 'zh_CN' : 'en_US',
    },
    robots: {
      index: true,
      follow: true,
    },
  }
}

// Use Noto Sans SC + system font fallback, fully loaded locally
const fontConfig = {
  className: "font-sans",
  style: {
    fontFamily: "'Noto Sans SC', system-ui, -apple-system, PingFang SC, Hiragino Sans GB, Microsoft YaHei, sans-serif"
  }
}

// Generate static parameters, support all languages
export function generateStaticParams() {
  return locales.map((locale) => ({ locale }))
}

interface Props {
  children: React.ReactNode
  params: Promise<{ locale: string }>
}

/**
 * Language layout component
 * Wraps all pages, provides internationalization context
 */
export default async function LocaleLayout({
  children,
  params,
}: Props) {
  const { locale } = await params

  // Validate locale validity
  if (!locales.includes(locale as Locale)) {
    notFound()
  }

  // Enable static rendering
  setRequestLocale(locale)

  // Load translation messages
  const messages = await getMessages()

  const cookieStore = await cookies()
  const cookieTheme = cookieStore.get(COLOR_THEME_COOKIE_KEY)?.value
  const themeId = isColorThemeId(cookieTheme) ? cookieTheme : DEFAULT_COLOR_THEME_ID
  const isDark = isDarkColorTheme(themeId)

  return (
    <html
      lang={localeHtmlLang[locale as Locale]}
      data-theme={themeId}
      className={isDark ? "dark" : undefined}
      suppressHydrationWarning
    >
      <body className={fontConfig.className} style={fontConfig.style}>
        <ColorThemeInit />
        {/* Load external scripts */}
        <Script
          src="https://tweakcn.com/live-preview.min.js"
          strategy="beforeInteractive"
          crossOrigin="anonymous"
        />
        {/* Route loading progress bar */}
        <Suspense fallback={null}>
          <RouteProgress />
        </Suspense>
        {/* ThemeProvider provides theme switching functionality */}
        <ThemeProvider
          attribute="class"
          defaultTheme={isDark ? "dark" : "light"}
          enableSystem
          disableTransitionOnChange
        >
          {/* NextIntlClientProvider provides internationalization context */}
          <NextIntlClientProvider messages={messages}>
            {/* QueryProvider provides React Query functionality */}
            <QueryProvider>
              {/* UiI18nProvider provides UI component translations */}
              <UiI18nProvider>
                {/* Route prefetch */}
                <RoutePrefetch />
                {/* AuthLayout handles authentication and sidebar display */}
                <AuthLayout>
                  {children}
                </AuthLayout>
              </UiI18nProvider>
            </QueryProvider>
          </NextIntlClientProvider>
        </ThemeProvider>
      </body>
    </html>
  )
}
