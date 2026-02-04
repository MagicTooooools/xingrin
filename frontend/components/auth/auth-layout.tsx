"use client"

import React from "react"
import { usePathname } from "next/navigation"
import { useLocale, useTranslations } from "next-intl"
import { AppSidebar } from "@/components/app-sidebar"
import { UnifiedHeader } from "@/components/unified-header"
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar"
import { Toaster } from "@/components/ui/sonner"
import { LoadingState } from "@/components/loading-spinner"
import { Suspense } from "react"
import { useAuth } from "@/hooks/use-auth"
import { useRouter } from "next/navigation"

// Public routes that don't require authentication (without locale prefix)
const PUBLIC_ROUTES = ["/login"]

interface AuthLayoutProps {
  children: React.ReactNode
}

/**
 * Check if the current path is a public route
 * Handles internationalized paths like /en/login, /zh/login
 */
function isPublicPath(pathname: string): boolean {
  // Remove locale prefix (e.g., /en/login -> /login, /zh/login -> /login)
  const pathWithoutLocale = pathname.replace(/^\/[a-z]{2}(?=\/|$)/, '')
  return PUBLIC_ROUTES.some((route) => 
    pathWithoutLocale === route || pathWithoutLocale.startsWith(`${route}/`)
  )
}

/**
 * Authentication layout component
 * Decides whether to show sidebar based on login status and route
 * 
 * 新布局结构：
 * ┌─────────────────────────────────────────────────────────┐
 * │  Logo区域 (固定宽度)  │  顶栏内容 (搜索/通知/语言等)     │
 * ├──────────────────────┼──────────────────────────────────┤
 * │  侧边栏菜单          │  主内容区域                       │
 * │  (无Logo)            │                                  │
 * └──────────────────────┴──────────────────────────────────┘
 */
export function AuthLayout({ children }: AuthLayoutProps) {
  const pathname = usePathname()
  const router = useRouter()
  const { data: auth, isLoading } = useAuth()
  const tCommon = useTranslations("common")
  const locale = useLocale()

  // Check if it's a public route (login page)
  const isPublicRoute = isPublicPath(pathname)

  // Redirect to login page if not authenticated (useEffect must be before all conditional returns)
  React.useEffect(() => {
    if (!isLoading && !auth?.authenticated && !isPublicRoute) {
      const normalized = "/login/"
      const loginPath = `/${locale}${normalized}`
      router.push(loginPath)
    }
  }, [auth, isLoading, isPublicRoute, router, locale])

  // If it's login page, render content directly (without sidebar)
  if (isPublicRoute) {
    return (
      <>
        {children}
        <Toaster />
      </>
    )
  }

  const showLoading = isLoading || !auth?.authenticated
  const canRenderApp = !isLoading && !!auth?.authenticated

  // Authenticated - show full layout with unified header
  // 布局结构：
  // ┌─────────────────────────────────────────────────────────┐
  // │  Logo区域 (固定宽度)  │  顶栏内容 (搜索/通知/语言等)     │
  // ├──────────────────────┼──────────────────────────────────┤
  // │  侧边栏菜单          │  主内容区域                       │
  // │  (无Logo)            │                                  │
  // └──────────────────────┴──────────────────────────────────┘
  return (
    <>
      <LoadingState active={showLoading} message="loading..." />
      {canRenderApp ? (
        <SidebarProvider
          className="animate-app-fade-in !min-h-0 flex flex-col h-svh"
          style={
            {
              "--sidebar-width": "calc(var(--spacing) * 62)",
              "--header-height": "calc(var(--spacing) * 12)",
            } as React.CSSProperties
          }
        >
          {/* 统一顶栏 - 横跨整个页面，包含 Logo */}
          <UnifiedHeader />
          
          {/* 下方内容区：侧边栏 + 主内容 */}
          <div className="flex flex-1 min-h-0">
            <AppSidebar />
            <SidebarInset className="flex min-h-0 flex-col flex-1">
              <div className="flex flex-col flex-1 min-h-0 overflow-y-auto">
                <div className="@container/main flex-1 min-h-0 flex flex-col gap-2">
                  <Suspense fallback={<LoadingState message={tCommon("status.pageLoading")} />}>
                    {children}
                  </Suspense>
                </div>
              </div>
            </SidebarInset>
          </div>
        </SidebarProvider>
      ) : null}
      <Toaster />
    </>
  )
}
