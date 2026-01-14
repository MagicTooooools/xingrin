"use client"

import * as React from "react"
import dynamic from "next/dynamic"
import { cn } from "@/lib/utils"

// Dynamic import to avoid SSR issues with GSAP
const Shuffle = dynamic(() => import("@/components/Shuffle"), { ssr: false })

type LoginStep = "username" | "password" | "authenticating" | "success" | "error"

interface TerminalLoginTranslations {
  title: string
  subtitle: string
  usernamePrompt: string
  passwordPrompt: string
  authenticating: string
  processing: string
  accessGranted: string
  welcomeMessage: string
  authFailed: string
  invalidCredentials: string
  shortcuts: string
  submit: string
  cancel: string
  clear: string
  startEnd: string
}

interface TerminalLine {
  text: string
  type: "prompt" | "input" | "info" | "success" | "error" | "warning"
}

interface TerminalLoginProps {
  onLogin: (username: string, password: string) => Promise<void>
  isPending?: boolean
  className?: string
  translations: TerminalLoginTranslations
}

export function TerminalLogin({
  onLogin,
  isPending = false,
  className,
  translations: t,
}: TerminalLoginProps) {
  const [step, setStep] = React.useState<LoginStep>("username")
  const [username, setUsername] = React.useState("")
  const [password, setPassword] = React.useState("")
  const [lines, setLines] = React.useState<TerminalLine[]>([])
  const [cursorPosition, setCursorPosition] = React.useState(0)
  const [isFocused, setIsFocused] = React.useState(false)
  const inputRef = React.useRef<HTMLInputElement>(null)
  const containerRef = React.useRef<HTMLDivElement>(null)

  // Focus input on mount and when step changes
  React.useEffect(() => {
    inputRef.current?.focus()
  }, [step])

  // Click anywhere to focus input
  const handleContainerClick = () => {
    inputRef.current?.focus()
  }

  const addLine = (line: TerminalLine) => {
    setLines((prev) => [...prev, line])
  }

  const getCurrentValue = () => {
    if (step === "username") return username
    if (step === "password") return password
    return ""
  }

  const setCurrentValue = (value: string) => {
    if (step === "username") {
      setUsername(value)
      setCursorPosition(value.length)
    } else if (step === "password") {
      setPassword(value)
      setCursorPosition(value.length)
    }
  }

  const handleKeyDown = async (e: React.KeyboardEvent<HTMLInputElement>) => {
    const value = getCurrentValue()

    // Ctrl+C - Cancel/Clear current input
    if (e.ctrlKey && e.key === "c") {
      e.preventDefault()
      if (step === "username" || step === "password") {
        addLine({ text: `^C`, type: "warning" })
        setCurrentValue("")
        setCursorPosition(0)
      }
      return
    }

    // Ctrl+U - Clear line (delete from cursor to start)
    if (e.ctrlKey && e.key === "u") {
      e.preventDefault()
      setCurrentValue("")
      setCursorPosition(0)
      return
    }

    // Ctrl+A - Move cursor to start
    if (e.ctrlKey && e.key === "a") {
      e.preventDefault()
      setCursorPosition(0)
      if (inputRef.current) {
        inputRef.current.setSelectionRange(0, 0)
      }
      return
    }

    // Ctrl+E - Move cursor to end
    if (e.ctrlKey && e.key === "e") {
      e.preventDefault()
      setCursorPosition(value.length)
      if (inputRef.current) {
        inputRef.current.setSelectionRange(value.length, value.length)
      }
      return
    }

    // Ctrl+W - Delete word before cursor
    if (e.ctrlKey && e.key === "w") {
      e.preventDefault()
      const beforeCursor = value.slice(0, cursorPosition)
      const afterCursor = value.slice(cursorPosition)
      const lastSpace = beforeCursor.trimEnd().lastIndexOf(" ")
      const newBefore = lastSpace === -1 ? "" : beforeCursor.slice(0, lastSpace + 1)
      setCurrentValue(newBefore + afterCursor)
      setCursorPosition(newBefore.length)
      return
    }

    // Tab - Move to next field (username -> password)
    if (e.key === "Tab" && step === "username") {
      e.preventDefault()
      if (!username.trim()) return
      addLine({ text: `> ${t.usernamePrompt}: `, type: "prompt" })
      addLine({ text: username, type: "input" })
      setStep("password")
      setCursorPosition(0)
      return
    }

    // Enter - Submit
    if (e.key === "Enter") {
      if (step === "username") {
        if (!username.trim()) return
        addLine({ text: `> ${t.usernamePrompt}: `, type: "prompt" })
        addLine({ text: username, type: "input" })
        setStep("password")
        setCursorPosition(0)
      } else if (step === "password") {
        if (!password.trim()) return
        addLine({ text: `> ${t.passwordPrompt}: `, type: "prompt" })
        addLine({ text: "*".repeat(password.length), type: "input" })
        addLine({ text: "", type: "info" })
        addLine({ text: `> ${t.authenticating}`, type: "warning" })
        setStep("authenticating")

        try {
          await onLogin(username, password)
          addLine({ text: `> ${t.accessGranted}`, type: "success" })
          addLine({ text: `> ${t.welcomeMessage}`, type: "success" })
          setStep("success")
        } catch {
          addLine({ text: `> ${t.authFailed}`, type: "error" })
          addLine({ text: `> ${t.invalidCredentials}`, type: "error" })
          addLine({ text: "", type: "info" })
          setStep("error")
          setTimeout(() => {
            setUsername("")
            setPassword("")
            setLines([])
            setCursorPosition(0)
            setStep("username")
          }, 2000)
        }
      }
      return
    }
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setCurrentValue(value)
    setCursorPosition(e.target.selectionStart || value.length)
  }

  const handleSelect = (e: React.SyntheticEvent<HTMLInputElement>) => {
    const target = e.target as HTMLInputElement
    setCursorPosition(target.selectionStart || 0)
  }

  const isInputDisabled = step === "authenticating" || step === "success" || isPending

  const getCurrentPrompt = () => {
    if (step === "username") return `> ${t.usernamePrompt}: `
    if (step === "password") return `> ${t.passwordPrompt}: `
    return "> "
  }

  const getDisplayValue = () => {
    if (step === "username") return username
    if (step === "password") return "*".repeat(password.length)
    return ""
  }

  // Render cursor at position
  const renderInputWithCursor = () => {
    const displayValue = getDisplayValue()
    const before = displayValue.slice(0, cursorPosition)
    const after = displayValue.slice(cursorPosition)
    const cursorChar = after[0] || ""

    if (!isFocused) {
      return <span className="text-zinc-100">{displayValue}</span>
    }

    return (
      <>
        <span className="text-zinc-100">{before}</span>
        <span className="animate-blink inline-block min-w-[0.6em] bg-green-500 text-black">
          {cursorChar || "\u00A0"}
        </span>
        <span className="text-zinc-100">{after.slice(1)}</span>
      </>
    )
  }

  return (
    <div
      ref={containerRef}
      onClick={handleContainerClick}
      className={cn(
        "border-zinc-700 bg-zinc-900/80 backdrop-blur-sm z-0 w-full max-w-2xl rounded-xl border cursor-text",
        className
      )}
    >
      {/* Terminal header */}
      <div className="border-zinc-700 flex items-center gap-x-2 border-b px-4 py-3">
        <div className="flex flex-row gap-x-2">
          <div className="h-3 w-3 rounded-full bg-red-500"></div>
          <div className="h-3 w-3 rounded-full bg-yellow-500"></div>
          <div className="h-3 w-3 rounded-full bg-green-500"></div>
        </div>
        <span className="ml-2 text-xs text-zinc-400 font-mono">{t.title}</span>
      </div>

      {/* Terminal content */}
      <div className="p-4 font-mono text-sm min-h-[280px]">
        {/* Shuffle Title Banner */}
        <div className="mb-6 text-center">
          <Shuffle
            text="STAR PATROL"
            className="!text-4xl sm:!text-5xl md:!text-6xl !font-bold text-cyan-500"
            shuffleDirection="up"
            duration={0.5}
            stagger={0.04}
            shuffleTimes={2}
            triggerOnHover={true}
            triggerOnce={false}
            autoPlay={false}
          />
          <div className="text-zinc-400 text-sm mt-3">
            ─────────── {t.subtitle} ───────────
          </div>
        </div>

        {/* ========== Mobile Form ========== */}
        <div className="sm:hidden">
          {(step === "username" || step === "password" || step === "error") && (
            <form
              onSubmit={async (e) => {
                e.preventDefault()
                if (!username.trim() || !password.trim()) return
                setStep("authenticating")
                try {
                  await onLogin(username, password)
                  setStep("success")
                } catch {
                  setStep("error")
                  setTimeout(() => {
                    setUsername("")
                    setPassword("")
                    setStep("username")
                  }, 2000)
                }
              }}
              className="space-y-4"
            >
              <div>
                <label className="text-green-500 text-xs mb-1 block">{t.usernamePrompt}</label>
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  disabled={isInputDisabled}
                  className="w-full bg-zinc-800 border border-zinc-600 rounded px-3 py-2 text-zinc-100 outline-none focus:border-green-500 font-mono text-sm"
                  autoComplete="username"
                />
              </div>
              <div>
                <label className="text-green-500 text-xs mb-1 block">{t.passwordPrompt}</label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={isInputDisabled}
                  className="w-full bg-zinc-800 border border-zinc-600 rounded px-3 py-2 text-zinc-100 outline-none focus:border-green-500 font-mono text-sm"
                  autoComplete="current-password"
                />
              </div>
              {step === "error" && (
                <p className="text-red-500 text-sm">{t.invalidCredentials}</p>
              )}
              <button
                type="submit"
                disabled={isInputDisabled}
                className="w-full py-2 px-4 bg-green-600 hover:bg-green-500 disabled:opacity-50 text-black font-mono text-sm rounded transition-colors"
              >
                {t.submit}
              </button>
            </form>
          )}
          {step === "authenticating" && (
            <div className="text-yellow-500 text-center py-4">
              <span className="animate-pulse">{t.processing}</span>
            </div>
          )}
          {step === "success" && (
            <div className="text-green-500 text-center py-4">
              {t.accessGranted}
            </div>
          )}
        </div>

        {/* ========== Desktop Terminal ========== */}
        <div className="hidden sm:block">
          {/* Previous lines */}
          {lines.map((line, index) => (
            <span
              key={index}
              className={cn(
                "whitespace-pre-wrap",
                line.type === "prompt" && "text-green-500",
                line.type === "input" && "text-zinc-100",
                line.type === "info" && "text-zinc-500",
                line.type === "success" && "text-green-500",
                line.type === "error" && "text-red-500",
                line.type === "warning" && "text-yellow-500"
              )}
            >
              {line.text}
              {(line.type === "prompt" || line.text === "") ? "" : "\n"}
            </span>
          ))}

          {/* Current input line */}
          {(step === "username" || step === "password") && (
            <div className="flex items-center">
              <span className="text-green-500">{getCurrentPrompt()}</span>
              {renderInputWithCursor()}
              <input
                ref={inputRef}
                type={step === "password" ? "password" : "text"}
                value={getCurrentValue()}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown}
                onSelect={handleSelect}
                onFocus={() => setIsFocused(true)}
                onBlur={() => setIsFocused(false)}
                disabled={isInputDisabled}
                className="absolute opacity-0 pointer-events-none"
                autoComplete={step === "username" ? "username" : "current-password"}
                autoFocus
              />
            </div>
          )}

          {/* Loading indicator */}
          {step === "authenticating" && (
            <div className="flex items-center text-yellow-500">
              <span className="animate-pulse">{t.processing}</span>
            </div>
          )}

          {/* Keyboard shortcuts hint */}
          {(step === "username" || step === "password") && (
            <div className="mt-6 text-xs text-zinc-600">
              <span className="text-zinc-500">{t.shortcuts}:</span>{" "}
              <span className="text-cyan-600">Enter</span> {t.submit}{" "}
              <span className="text-cyan-600">Ctrl+C</span> {t.cancel}{" "}
              <span className="text-cyan-600">Ctrl+U</span> {t.clear}{" "}
              <span className="text-cyan-600">Ctrl+A/E</span> {t.startEnd}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
