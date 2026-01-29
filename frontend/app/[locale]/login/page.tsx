"use client"

import React from "react"
import { useRouter } from "next/navigation"
import { useTranslations } from "next-intl"
import { useQueryClient } from "@tanstack/react-query"
import { TerminalLogin } from "@/components/ui/terminal-login"
import { useLogin, useAuth } from "@/hooks/use-auth"
import { vulnerabilityKeys } from "@/hooks/use-vulnerabilities"
import { useRoutePrefetch } from "@/hooks/use-route-prefetch"
import { getAssetStatistics, getStatisticsHistory } from "@/services/dashboard.service"
import { getScans } from "@/services/scan.service"
import { VulnerabilityService } from "@/services/vulnerability.service"

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

  const [isReady, setIsReady] = React.useState(false)

  // Hide the inline boot splash and show login content
  React.useEffect(() => {
    // Small delay to ensure smooth transition
    const timer = setTimeout(() => {
      const splash = document.getElementById('boot-splash')
      if (splash) {
        splash.classList.add('hidden')
        setTimeout(() => splash.remove(), 300)
      }
      setIsReady(true)
    }, 100)
    return () => clearTimeout(timer)
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

  return (
    <div className="relative flex min-h-svh flex-col bg-black">
      {/* Circuit Board Animation */}
      <div className={`fixed inset-0 z-0 transition-opacity duration-300 ${isReady ? "opacity-100" : "opacity-0"}`}>
        <div className="circuit-container">
          {/* Grid pattern */}
          <div className="circuit-grid" />
          
          {/* === Main backbone traces === */}
          {/* Horizontal main lines - 6 lines */}
          <div className="trace trace-h" style={{ top: '12%', left: 0, width: '100%' }}>
            <div className="trace-glow" style={{ animationDuration: '6s' }} />
          </div>
          <div className="trace trace-h" style={{ top: '28%', left: 0, width: '100%' }}>
            <div className="trace-glow" style={{ animationDelay: '1s', animationDuration: '5s' }} />
          </div>
          <div className="trace trace-h" style={{ top: '44%', left: 0, width: '100%' }}>
            <div className="trace-glow" style={{ animationDelay: '2s', animationDuration: '5.5s' }} />
          </div>
          <div className="trace trace-h" style={{ top: '60%', left: 0, width: '100%' }}>
            <div className="trace-glow" style={{ animationDelay: '3s', animationDuration: '4.5s' }} />
          </div>
          <div className="trace trace-h" style={{ top: '76%', left: 0, width: '100%' }}>
            <div className="trace-glow" style={{ animationDelay: '4s', animationDuration: '5s' }} />
          </div>
          <div className="trace trace-h" style={{ top: '92%', left: 0, width: '100%' }}>
            <div className="trace-glow" style={{ animationDelay: '5s', animationDuration: '6s' }} />
          </div>
          
          {/* Vertical main lines - 6 lines */}
          <div className="trace trace-v" style={{ left: '8%', top: 0, height: '100%' }}>
            <div className="trace-glow trace-glow-v" style={{ animationDelay: '0.5s', animationDuration: '7s' }} />
          </div>
          <div className="trace trace-v" style={{ left: '24%', top: 0, height: '100%' }}>
            <div className="trace-glow trace-glow-v" style={{ animationDelay: '1.5s', animationDuration: '6s' }} />
          </div>
          <div className="trace trace-v" style={{ left: '40%', top: 0, height: '100%' }}>
            <div className="trace-glow trace-glow-v" style={{ animationDelay: '2.5s', animationDuration: '5.5s' }} />
          </div>
          <div className="trace trace-v" style={{ left: '56%', top: 0, height: '100%' }}>
            <div className="trace-glow trace-glow-v" style={{ animationDelay: '3.5s', animationDuration: '6.5s' }} />
          </div>
          <div className="trace trace-v" style={{ left: '72%', top: 0, height: '100%' }}>
            <div className="trace-glow trace-glow-v" style={{ animationDelay: '4.5s', animationDuration: '5s' }} />
          </div>
          <div className="trace trace-v" style={{ left: '88%', top: 0, height: '100%' }}>
            <div className="trace-glow trace-glow-v" style={{ animationDelay: '5.5s', animationDuration: '6s' }} />
          </div>
          
        </div>
        
        <style jsx>{`
          .circuit-container {
            position: absolute;
            inset: 0;
            background: #0a0a0a;
            overflow: hidden;
          }
          
          .circuit-grid {
            position: absolute;
            inset: 0;
            background-image: 
              linear-gradient(rgba(255, 255, 255, 0.03) 1px, transparent 1px),
              linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
            background-size: 40px 40px;
          }
          
          .trace {
            position: absolute;
            background: rgba(255, 255, 255, 0.15);
            overflow: hidden;
          }
          
          .trace-h {
            height: 2px;
          }
          
          .trace-v {
            width: 2px;
          }
          
          .trace-glow {
            position: absolute;
            top: -2px;
            left: -20%;
            width: 30%;
            height: 6px;
            background: linear-gradient(90deg, transparent, #fff, #aaa, transparent);
            animation: traceFlow 3s linear infinite;
            filter: blur(2px);
          }
          
          .trace-glow-v {
            top: -20%;
            left: -2px;
            width: 6px;
            height: 30%;
            background: linear-gradient(180deg, transparent, #fff, #aaa, transparent);
            animation: traceFlowV 3s linear infinite;
          }
          
          @keyframes traceFlow {
            0% { left: -30%; }
            100% { left: 100%; }
          }
          
          @keyframes traceFlowV {
            0% { top: -30%; }
            100% { top: 100%; }
          }
        `}</style>
      </div>

      {/* Fingerprint identifier - for FOFA/Shodan and other search engines to identify */}
      <meta name="generator" content="Orbit ASM Platform" />

      {/* Main content area */}
      <div
        className={`relative z-10 flex-1 flex items-center justify-center p-6 transition-[opacity,transform] duration-300 ${
          isReady ? "opacity-100 translate-y-0" : "opacity-0 translate-y-2"
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
          isReady ? "opacity-100" : "opacity-0"
        }`}
      >
        <p className="text-xs text-muted-foreground">
          {process.env.NEXT_PUBLIC_VERSION || "dev"}
        </p>
      </div>
    </div>
  )
}
