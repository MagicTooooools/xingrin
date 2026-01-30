"use client"

import React from "react"
import { usePathname } from "next/navigation"
import { useTranslations } from "next-intl"
import { AppSidebar } from "@/components/app-sidebar"
import { SiteHeader } from "@/components/site-header"
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
 */
export function AuthLayout({ children }: AuthLayoutProps) {
  const pathname = usePathname()
  const router = useRouter()
  const { data: auth, isLoading } = useAuth()
  const tCommon = useTranslations("common")

  // Check if it's a public route (login page)
  const isPublicRoute = isPublicPath(pathname)

  // Hide the inline boot splash once app is ready (non-login routes only)
  React.useEffect(() => {
    if (isPublicRoute) return
    const splash = document.getElementById('boot-splash')
    if (splash) {
      splash.classList.add('hidden')
    }
  }, [isPublicRoute])

  // Redirect to login page if not authenticated (useEffect must be before all conditional returns)
  React.useEffect(() => {
    if (!isLoading && !auth?.authenticated && !isPublicRoute) {
      router.push("/login/")
    }
  }, [auth, isLoading, isPublicRoute, router])

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

  // Authenticated - show full layout (with sidebar)
  return (
    <>
      <LoadingState active={showLoading} message="loading..." />
      {canRenderApp ? (
        <SidebarProvider
          style={
            {
              "--sidebar-width": "calc(var(--spacing) * 62)",
              "--header-height": "calc(var(--spacing) * 11)",
            } as React.CSSProperties
          }
        >
          <AppSidebar />
          <SidebarInset className="flex min-h-0 flex-col h-svh">
            <SiteHeader />
            <div className="flex flex-col flex-1 min-h-0 overflow-y-auto">
              <div className="@container/main flex-1 min-h-0 flex flex-col gap-2">
                <Suspense fallback={<LoadingState message={tCommon("status.pageLoading")} />}>
                  {children}
                </Suspense>
                <Toaster />
              </div>
            </div>
          </SidebarInset>
        </SidebarProvider>
      ) : null}
    </>
  )
}
