"use client"

import * as React from "react"

import { cn } from "@/lib/utils"

type BootLine = {
  text: string
  className?: string
}

const BOOT_LINES: BootLine[] = [
  { text: "> booting ORBIT...", className: "text-yellow-500" },
  { text: "> initializing secure terminal...", className: "text-zinc-200" },
  { text: "> loading modules: auth, i18n, ui...", className: "text-zinc-200" },
  { text: "> checking session...", className: "text-yellow-500" },
  { text: "> ready.", className: "text-green-500" },
]

const SUCCESS_LINES: BootLine[] = [
  { text: "> authentication successful", className: "text-green-500" },
  { text: "> loading user profile...", className: "text-zinc-200" },
  { text: "> initializing dashboard...", className: "text-zinc-200" },
  { text: "> preparing workspace...", className: "text-yellow-500" },
  { text: "> access granted.", className: "text-green-500" },
]

// Keep the log animation snappy so it can complete within the 0.6s splash.
const STEP_DELAYS_MS = [70, 90, 110, 130, 150]

const GLITCH_MS = 600

export function LoginBootScreen({ className, success = false }: { className?: string; success?: boolean }) {
  const [visible, setVisible] = React.useState(0)
  const [entered, setEntered] = React.useState(false)
  const [glitchOn, setGlitchOn] = React.useState(true)

  // 根据 success 状态选择显示的行
  const displayLines = success ? SUCCESS_LINES : BOOT_LINES

  React.useEffect(() => {
    const raf = requestAnimationFrame(() => setEntered(true))
    return () => cancelAnimationFrame(raf)
  }, [])

  React.useEffect(() => {
    setGlitchOn(true)
    const timer = setTimeout(() => setGlitchOn(false), GLITCH_MS)
    return () => clearTimeout(timer)
  }, [])

  React.useEffect(() => {
    setVisible(0)

    const timers: Array<ReturnType<typeof setTimeout>> = []
    let acc = 0

    for (let i = 0; i < displayLines.length; i++) {
      acc += STEP_DELAYS_MS[i] ?? 160
      timers.push(
        setTimeout(() => {
          setVisible((prev) => Math.max(prev, i + 1))
        }, acc)
      )
    }

    return () => {
      timers.forEach(clearTimeout)
    }
  }, [displayLines])

  const progress = Math.round((Math.min(visible, displayLines.length) / displayLines.length) * 100)

  return (
    <div className={cn("relative flex min-h-svh flex-col bg-black", glitchOn && "orbit-splash-glitch", className)}>
      {/* Main content area */}
      <div className="relative z-10 flex-1 flex items-center justify-center p-6">
        <div
          className={cn(
            "border-zinc-700 bg-zinc-900/80 backdrop-blur-sm z-0 w-full max-w-xl rounded-xl border transition-opacity duration-200 ease-out motion-reduce:transition-none",
            entered ? "opacity-100" : "opacity-0"
          )}
        >
          {/* Terminal header */}
          <div className="border-zinc-700 flex items-center gap-x-2 border-b px-4 py-3">
            <div className="flex flex-row gap-x-2">
              <div className="h-3 w-3 rounded-full bg-red-500" />
              <div className="h-3 w-3 rounded-full bg-yellow-500" />
              <div className="h-3 w-3 rounded-full bg-green-500" />
            </div>
            <span className="ml-2 text-xs text-zinc-400 font-mono">ORBIT · boot</span>
            <span className="ml-auto text-xs text-zinc-500 font-mono">{progress}%</span>
          </div>

          {/* Terminal body */}
          <div className="p-4 font-mono text-sm min-h-[280px]">
            <div className="mb-6 text-center">
              <div
                className={cn(
                  "text-3xl sm:text-4xl !font-bold tracking-wide",
                  "bg-gradient-to-r from-[#FF10F0] via-[#B026FF] to-[#FF10F0] bg-clip-text text-transparent",
                  glitchOn && "orbit-glitch-text"
                )}
                data-text="ORBIT"
                style={{
                  filter: "drop-shadow(0 0 20px rgba(255, 16, 240, 0.5)) drop-shadow(0 0 40px rgba(176, 38, 255, 0.3))"
                }}
              >
                ORBIT
              </div>
              <div className="mt-3 flex items-center gap-3 text-zinc-400 text-xs">
                <span className="h-px flex-1 bg-gradient-to-r from-transparent via-[#B026FF] to-transparent" />
                <span className="whitespace-nowrap">system bootstrap</span>
                <span className="h-px flex-1 bg-gradient-to-r from-transparent via-[#B026FF] to-transparent" />
              </div>
            </div>

            <div className="space-y-1">
              {displayLines.slice(0, visible).map((line, idx) => (
                <div key={idx} className={cn("whitespace-pre-wrap", line.className)}>
                  {line.text}
                </div>
              ))}

              {/* Cursor */}
              <div className="text-green-500">
                <span className="inline-block h-4 w-2 align-middle bg-green-500 animate-pulse" />
              </div>
            </div>

            {/* Progress bar */}
            <div className="mt-6">
              <div className="h-1.5 w-full rounded bg-zinc-800 overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-[#FF10F0] to-[#B026FF]"
                  style={{
                    width: `${progress}%`,
                    boxShadow: "0 0 10px rgba(255, 16, 240, 0.5), 0 0 20px rgba(176, 38, 255, 0.3)"
                  }}
                />
              </div>
              <div className="mt-2 text-xs text-zinc-500">
                Checking session…
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
