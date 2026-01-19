"use client"

import React from "react"
import { useRouter } from "next/navigation"
import { useTranslations } from "next-intl"
import dynamic from "next/dynamic"
import { LoginBootScreen } from "@/components/auth/login-boot-screen"
import { TerminalLogin } from "@/components/ui/terminal-login"
import { useLogin, useAuth } from "@/hooks/use-auth"
import { useRoutePrefetch } from "@/hooks/use-route-prefetch"

// Dynamic import to avoid SSR issues with WebGL
const PixelBlast = dynamic(() => import("@/components/PixelBlast"), { ssr: false })

const BOOT_SPLASH_MS = 600
const BOOT_FADE_MS = 200
const LOGIN_SUCCESS_DELAY_MS = 1200 // 登录成功后显示启动屏幕的时间
const LOGIN_SUCCESS_FADE_MS = 500 // 登录成功后淡出的时间

type BootOverlayPhase = "entering" | "visible" | "leaving" | "hidden"

export default function LoginPage() {
  // Preload all page components on login page
  useRoutePrefetch()
  const router = useRouter()
  const { data: auth, isLoading: authLoading } = useAuth()
  const { mutateAsync: login, isPending } = useLogin()
  const t = useTranslations("auth.terminal")

  // Always show a short splash on entering the login page.
  const [bootMinDone, setBootMinDone] = React.useState(false)
  const [bootPhase, setBootPhase] = React.useState<BootOverlayPhase>("entering")
  const [loginSuccess, setLoginSuccess] = React.useState(false) // 跟踪登录成功状态
  const [showSuccessSplash, setShowSuccessSplash] = React.useState(false) // 是否显示登录成功的启动屏幕

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

    setBootPhase("leaving")
  }, [authLoading, bootMinDone, bootPhase])

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

  // If already logged in, show success splash then redirect to dashboard
  React.useEffect(() => {
    if (!bootMinDone) return
    if (authLoading) return
    if (auth?.authenticated && !showSuccessSplash) {
      // 登录成功，显示成功启动屏幕
      setShowSuccessSplash(true)
      setLoginSuccess(true)

      // 延迟后开始淡出并跳转
      const successTimer = setTimeout(() => {
        router.push("/dashboard/")
      }, LOGIN_SUCCESS_DELAY_MS + LOGIN_SUCCESS_FADE_MS)

      return () => clearTimeout(successTimer)
    }
  }, [auth?.authenticated, authLoading, bootMinDone, router, showSuccessSplash])

  const handleLogin = async (username: string, password: string) => {
    await login({ username, password })
  }

  // While authenticated, keep showing the splash until redirect happens.
  if (auth?.authenticated && showSuccessSplash) {
    return <LoginBootScreen success={loginSuccess} />
  }

  const loginVisible = bootPhase === "leaving" || bootPhase === "hidden"

  return (
    <div className="relative flex min-h-svh flex-col bg-black">
      <div className={`fixed inset-0 z-0 transition-opacity duration-300 ${loginVisible ? "opacity-100" : "opacity-0"}`}>
        <PixelBlast
          className=""
          style={{}}
          pixelSize={6.5}
          patternScale={4.5}
          color="#06b6d4"
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
