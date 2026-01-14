"use client"

import React from "react"
import { useRouter } from "next/navigation"
import { useTranslations } from "next-intl"
import dynamic from "next/dynamic"
import { Spinner } from "@/components/ui/spinner"
import { TerminalLogin } from "@/components/ui/terminal-login"
import { useLogin, useAuth } from "@/hooks/use-auth"
import { useRoutePrefetch } from "@/hooks/use-route-prefetch"

// Dynamic import to avoid SSR issues with WebGL
const PixelBlast = dynamic(() => import("@/components/PixelBlast"), { ssr: false })

export default function LoginPage() {
  // Preload all page components on login page
  useRoutePrefetch()
  const router = useRouter()
  const { data: auth, isLoading: authLoading } = useAuth()
  const { mutateAsync: login, isPending } = useLogin()
  const t = useTranslations("auth.terminal")

  // If already logged in, redirect to dashboard
  React.useEffect(() => {
    if (auth?.authenticated) {
      router.push("/dashboard/")
    }
  }, [auth, router])

  const handleLogin = async (username: string, password: string) => {
    await login({ username, password })
  }

  // Show spinner while loading
  if (authLoading) {
    return (
      <div className="flex min-h-svh w-full flex-col items-center justify-center gap-4 bg-background">
        <Spinner className="size-8 text-primary" />
        <p className="text-muted-foreground text-sm" suppressHydrationWarning>loading...</p>
      </div>
    )
  }

  // Don't show login page if already logged in
  if (auth?.authenticated) {
    return null
  }

  return (
    <div className="relative flex min-h-svh flex-col bg-black">
      <div className="fixed inset-0 z-0">
        <PixelBlast
          style={{}}
          pixelSize={6}
          patternScale={4.5}
          color="#06b6d4"
        />
      </div>

      {/* Fingerprint identifier - for FOFA/Shodan and other search engines to identify */}
      <meta name="generator" content="Star Patrol ASM Platform" />

      {/* Main content area */}
      <div className="relative z-10 flex-1 flex items-center justify-center p-6">
        <TerminalLogin
          onLogin={handleLogin}
          isPending={isPending}
          translations={{
            title: t("title"),
            subtitle: t("subtitle"),
            usernamePrompt: t("usernamePrompt"),
            passwordPrompt: t("passwordPrompt"),
            authenticating: t("authenticating"),
            processing: t("processing"),
            accessGranted: t("accessGranted"),
            welcomeMessage: t("welcomeMessage"),
            authFailed: t("authFailed"),
            invalidCredentials: t("invalidCredentials"),
          }}
        />
      </div>

      {/* Version number - fixed at the bottom of the page */}
      <div className="relative z-10 flex-shrink-0 text-center py-4">
        <p className="text-xs text-muted-foreground">
          {process.env.NEXT_PUBLIC_VERSION || "dev"}
        </p>
      </div>
    </div>
  )
}
