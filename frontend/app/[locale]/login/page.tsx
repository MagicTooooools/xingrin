"use client"

import React from "react"
import { useRouter } from "next/navigation"
import { useTranslations } from "next-intl"
import { useQueryClient } from "@tanstack/react-query"
import dynamic from "next/dynamic"
import { LoginBootScreen } from "@/components/auth/login-boot-screen"
import { TerminalLogin } from "@/components/ui/terminal-login"
import { useLogin, useAuth } from "@/hooks/use-auth"
import { vulnerabilityKeys } from "@/hooks/use-vulnerabilities"
import { useRoutePrefetch } from "@/hooks/use-route-prefetch"
import { getAssetStatistics, getStatisticsHistory } from "@/services/dashboard.service"
import { getScans } from "@/services/scan.service"
import { VulnerabilityService } from "@/services/vulnerability.service"

// Dynamic import to avoid SSR issues with WebGL
const PixelBlast = dynamic(() => import("@/components/PixelBlast"), { ssr: false })

const BOOT_SPLASH_MS = 600
const BOOT_FADE_MS = 200

type BootOverlayPhase = "entering" | "visible" | "leaving" | "hidden"

export default function LoginPage() {
  // Preload all page components on login page
  useRoutePrefetch()
  const router = useRouter()
  const queryClient = useQueryClient()
  const { data: auth, isLoading: authLoading } = useAuth()
  const { mutateAsync: login, isPending } = useLogin()
  const t = useTranslations("auth.terminal")

  const loginStartedRef = React.useRef(false)
  const [loginReady, setLoginReady] = React.useState(false)

  const [pixelFirstFrame, setPixelFirstFrame] = React.useState(false)
  const handlePixelFirstFrame = React.useCallback(() => {
    setPixelFirstFrame(true)
  }, [])

  // 提取预加载逻辑为可复用函数
  const prefetchDashboardData = React.useCallback(async () => {
    const scansParams = { page: 1, pageSize: 10 }
    const vulnsParams = { page: 1, pageSize: 10 }

    return Promise.allSettled([
      queryClient.prefetchQuery({
        queryKey: ["asset", "statistics"],
        queryFn: getAssetStatistics,
      }),
      queryClient.prefetchQuery({
        queryKey: ["asset", "statistics", "history", 7],
        queryFn: () => getStatisticsHistory(7),
      }),
      queryClient.prefetchQuery({
        queryKey: ["scans", scansParams],
        queryFn: () => getScans(scansParams),
      }),
      queryClient.prefetchQuery({
        queryKey: vulnerabilityKeys.list(vulnsParams),
        queryFn: () => VulnerabilityService.getAllVulnerabilities(vulnsParams),
      }),
    ])
  }, [queryClient])

  // Always show a short splash on entering the login page.
  const [bootMinDone, setBootMinDone] = React.useState(false)
  const [bootPhase, setBootPhase] = React.useState<BootOverlayPhase>("entering")

  React.useEffect(() => {
    setBootMinDone(false)
    setBootPhase("entering")

    const bootTimer = setTimeout(() => setBootMinDone(true), BOOT_SPLASH_MS)
    const raf = requestAnimationFrame(() => setBootPhase("visible"))

    return () => {
      clearTimeout(bootTimer)
      cancelAnimationFrame(raf)
    }
  }, [])


  // Start hiding the splash after the minimum time AND auth check completes.
  // Note: don't schedule the fade-out timer in the same effect where we set `bootPhase`,
  // otherwise the effect cleanup will cancel the timer when `bootPhase` changes.
  React.useEffect(() => {
    if (bootPhase !== "visible") return
    if (!bootMinDone) return
    if (authLoading) return
    if (!pixelFirstFrame) return

    setBootPhase("leaving")
  }, [authLoading, bootMinDone, bootPhase, pixelFirstFrame])

  React.useEffect(() => {
    if (bootPhase !== "leaving") return

    const timer = setTimeout(() => setBootPhase("hidden"), BOOT_FADE_MS)
    return () => clearTimeout(timer)
  }, [bootPhase])

  // Memoize translations object to avoid recreating on every render
  const translations = React.useMemo(() => ({
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
    shortcuts: t("shortcuts"),
    submit: t("submit"),
    cancel: t("cancel"),
    clear: t("clear"),
    startEnd: t("startEnd"),
  }), [t])

  // If already logged in, warm up the dashboard, then redirect.
  React.useEffect(() => {
    if (authLoading) return
    if (!auth?.authenticated) return
    if (loginStartedRef.current) return

    let cancelled = false

    void (async () => {
      await prefetchDashboardData()

      if (cancelled) return
      router.replace("/dashboard/")
    })()

    return () => {
      cancelled = true
    }
  }, [auth?.authenticated, authLoading, prefetchDashboardData, router])

  React.useEffect(() => {
    if (!loginReady) return
    router.replace("/dashboard/")
  }, [loginReady, router])

  const handleLogin = async (username: string, password: string) => {
    loginStartedRef.current = true
    setLoginReady(false)

    // 并行执行独立操作：登录验证 + 预加载 dashboard bundle
    const [loginRes] = await Promise.all([
      login({ username, password }),
      router.prefetch("/dashboard/"),
    ])

    // 预加载 dashboard 数据
    await prefetchDashboardData()

    // Prime auth cache so AuthLayout doesn't flash a full-screen loading state.
    queryClient.setQueryData(["auth", "me"], {
      authenticated: true,
      user: loginRes.user,
    })

    setLoginReady(true)
  }

  const loginVisible = bootPhase === "leaving" || bootPhase === "hidden"

  return (
    <div className="relative flex min-h-svh flex-col bg-black">
      <div className={`fixed inset-0 z-0 transition-opacity duration-300 ${loginVisible ? "opacity-100" : "opacity-0"}`}>
        <PixelBlast
          onFirstFrame={handlePixelFirstFrame}
          className=""
          style={{}}
          pixelSize={6.5}
          patternScale={4.5}
          color="#FF10F0"
          speed={0.35}
          enableRipples={false}
        />
      </div>

      {/* Fingerprint identifier - for FOFA/Shodan and other search engines to identify */}
      <meta name="generator" content="Orbit ASM Platform" />

      {/* Main content area */}
      <div
        className={`relative z-10 flex-1 flex items-center justify-center p-6 transition-[opacity,transform] duration-300 ${
          loginVisible ? "opacity-100 translate-y-0" : "opacity-0 translate-y-2"
        }`}
      >
        <TerminalLogin
          onLogin={handleLogin}
          authDone={loginReady}
          isPending={isPending}
          translations={translations}
        />
      </div>

      {/* Version number - fixed at the bottom of the page */}
      <div
        className={`relative z-10 flex-shrink-0 text-center py-4 transition-opacity duration-300 ${
          loginVisible ? "opacity-100" : "opacity-0"
        }`}
      >
        <p className="text-xs text-muted-foreground">
          {process.env.NEXT_PUBLIC_VERSION || "dev"}
        </p>
      </div>

      {/* Full-page splash overlay */}
      {bootPhase !== "hidden" && (
        <div
          className={`fixed inset-0 z-50 transition-opacity ease-out ${
            bootPhase === "visible" ? "opacity-100" : "opacity-0 pointer-events-none"
          }`}
          style={{ transitionDuration: `${BOOT_FADE_MS}ms` }}
        >
          <LoginBootScreen />
        </div>
      )}
    </div>
  )
}
