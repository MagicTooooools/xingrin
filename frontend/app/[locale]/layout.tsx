import type React from "react"
import type { Metadata } from "next"
import { NextIntlClientProvider } from 'next-intl'
import { getMessages, setRequestLocale, getTranslations } from 'next-intl/server'
import { notFound } from 'next/navigation'
import { locales, localeHtmlLang, type Locale } from '@/i18n/config'

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

  return (
    <html lang={localeHtmlLang[locale as Locale]} suppressHydrationWarning>
      <head>
        {/* Inline critical CSS for instant boot splash - matches LoginBootScreen exactly */}
        <style suppressHydrationWarning dangerouslySetInnerHTML={{ __html: `
          #boot-splash {
            position: fixed;
            inset: 0;
            z-index: 9999;
            display: flex;
            flex-direction: column;
            min-height: 100svh;
            background: #0a0a0f;
            overflow: hidden;
            transition: opacity 0.2s ease-out;
          }
          #boot-splash.hidden {
            opacity: 0;
            pointer-events: none;
          }
          /* Animated gradient background */
          #boot-splash .bg-gradient {
            position: fixed;
            inset: 0;
            overflow: hidden;
          }
          #boot-splash .bg-blob {
            position: absolute;
            top: -50%;
            left: -50%;
            width: 200%;
            height: 200%;
            opacity: 0.3;
            background: conic-gradient(from 0deg at 50% 50%, #0a0a0f 0deg, #3f3f46 60deg, #0a0a0f 120deg, #52525b 180deg, #0a0a0f 240deg, #3f3f46 300deg, #0a0a0f 360deg);
            filter: blur(80px);
            animation: boot-blob 20s linear infinite;
          }
          #boot-splash .bg-overlay {
            position: absolute;
            inset: 0;
            background: rgba(10,10,15,0.8);
          }
          /* Grid background */
          #boot-splash .bg-grid {
            position: fixed;
            inset: 0;
            opacity: 0.4;
            background-image:
              linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px),
              linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px);
            background-size: 50px 50px;
            -webkit-mask-image: radial-gradient(circle at center, black, transparent 80%);
            mask-image: radial-gradient(circle at center, black, transparent 80%);
          }
          /* Main content */
          #boot-splash .content {
            position: relative;
            z-index: 10;
            flex: 1;
            display: flex;
            align-items: center;
            justify-content: center;
          }
          #boot-splash .center {
            text-align: center;
          }
          /* Logo container */
          #boot-splash .logo-container {
            position: relative;
            width: 200px;
            height: 200px;
            margin: 0 auto 40px;
          }
          #boot-splash .logo-spinner {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            border-radius: 50%;
            border: 4px solid transparent;
            border-top-color: rgba(255,255,255,0.8);
            animation: boot-spin 1s linear infinite;
          }
          #boot-splash .logo-spinner::before {
            content: '';
            position: absolute;
            top: -4px;
            left: -4px;
            right: -4px;
            bottom: -4px;
            border-radius: 50%;
            border: 4px solid transparent;
            border-top-color: rgba(200,200,200,0.6);
            animation: boot-spin 1.5s linear infinite reverse;
          }
          #boot-splash .logo-spinner::after {
            content: '';
            position: absolute;
            top: 8px;
            left: 8px;
            right: 8px;
            bottom: 8px;
            border-radius: 50%;
            border: 2px solid transparent;
            border-bottom-color: rgba(255,255,255,0.3);
            animation: boot-spin 2s linear infinite;
          }
          #boot-splash .logo {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 120px;
            height: 120px;
            z-index: 10;
          }
          /* Title */
          #boot-splash .title {
            font-size: 32px;
            font-weight: bold;
            letter-spacing: -0.025em;
            margin-bottom: 8px;
          }
          #boot-splash .title-luna {
            background: linear-gradient(to bottom right, #d4d4d8, #f4f4f5);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
          }
          #boot-splash .title-fox {
            background: linear-gradient(to bottom right, #a1a1aa, #e4e4e7);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
          }
          /* Loading status */
          #boot-splash .status {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 12px;
            margin-top: 24px;
          }
          #boot-splash .status-spinner {
            width: 16px;
            height: 16px;
            border: 2px solid rgba(255,255,255,0.2);
            border-top-color: rgba(255,255,255,0.8);
            border-radius: 50%;
            animation: boot-spin 0.8s linear infinite;
          }
          #boot-splash .status-text {
            font-size: 14px;
            color: #6b7280;
            font-weight: 500;
          }
          /* Progress bar */
          #boot-splash .progress {
            width: 240px;
            height: 4px;
            background: rgba(255,255,255,0.1);
            border-radius: 4px;
            margin: 24px auto 0;
            overflow: hidden;
          }
          #boot-splash .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #a1a1aa, #f4f4f5);
            border-radius: 4px;
            animation: boot-progress 3s ease-in-out infinite;
          }
          /* Dots */
          #boot-splash .dots {
            display: flex;
            justify-content: center;
            gap: 6px;
            margin-top: 32px;
          }
          #boot-splash .dot {
            width: 6px;
            height: 6px;
            background: rgba(255,255,255,0.3);
            border-radius: 50%;
            animation: boot-dot 1.5s ease-in-out infinite;
          }
          #boot-splash .dot:nth-child(2) { animation-delay: 0.2s; }
          #boot-splash .dot:nth-child(3) { animation-delay: 0.4s; }
          #boot-splash .dot:nth-child(4) { animation-delay: 0.6s; }
          #boot-splash .dot:nth-child(5) { animation-delay: 0.8s; }
          @keyframes boot-spin {
            to { transform: rotate(360deg); }
          }
          @keyframes boot-progress {
            0% { width: 0%; }
            50% { width: 70%; }
            100% { width: 100%; }
          }
          @keyframes boot-dot {
            0%, 100% { background: rgba(255,255,255,0.3); transform: scale(1); }
            50% { background: rgba(255,255,255,0.8); transform: scale(1.3); }
          }
          @keyframes boot-blob {
            0% { transform: translate(0,0) rotate(0deg); }
            33% { transform: translate(2%,2%) rotate(120deg); }
            66% { transform: translate(-2%,2%) rotate(240deg); }
            100% { transform: translate(0,0) rotate(360deg); }
          }
        `}} />
      </head>
      <body className={fontConfig.className} style={fontConfig.style}>
        {/* Inline boot splash - identical to LoginBootScreen, shows immediately before JS loads */}
        <div id="boot-splash" suppressHydrationWarning>
          <div className="bg-gradient">
            <div className="bg-blob"></div>
            <div className="bg-overlay"></div>
          </div>
          <div className="bg-grid"></div>
          <div className="content">
            <div className="center">
              <div className="logo-container">
                <div className="logo-spinner"></div>
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img className="logo" src="/images/icon-256.png" alt="" />
              </div>
              <h1 className="title">
                <span className="title-luna">Luna</span>
                <span className="title-fox">Fox</span>
              </h1>
              <div className="status">
                <div className="status-spinner"></div>
                <span className="status-text">Initializing...</span>
              </div>
              <div className="progress">
                <div className="progress-fill"></div>
              </div>
              <div className="dots">
                <div className="dot"></div>
                <div className="dot"></div>
                <div className="dot"></div>
                <div className="dot"></div>
                <div className="dot"></div>
              </div>
            </div>
          </div>
        </div>
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
          defaultTheme="dark"
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
