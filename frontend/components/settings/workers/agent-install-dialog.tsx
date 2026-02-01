"use client"

import { useCallback, useEffect, useMemo, useRef, useState } from "react"
import { useTranslations } from "next-intl"
import { IconChevronDown, IconChevronUp, IconLoader2 } from "@tabler/icons-react"
import { toast } from "sonner"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import { agentService } from "@/services/agent.service"
import { useFormatRelativeTime } from "@/lib/i18n-format"
import type { Agent, RegistrationTokenResponse } from "@/types/agent.types"
import { cn } from "@/lib/utils"

const FALLBACK_SERVER_URL = "https://your-lunafox-server"
const COPY_TOAST_ID = "install-command-copy"
const VERIFY_POLL_INTERVAL_MS = 5000
const VERIFY_TIMEOUT_MS = 60_000
const VERIFY_PAGE_SIZE = 200

const normalizeOrigin = (value: string): string => value.replace(/\/+$/, "")

const isLocalOrigin = (origin: string): boolean => {
  try {
    const url = new URL(origin)
    return url.hostname === "localhost" || url.hostname === "127.0.0.1" || url.hostname === "::1"
  } catch {
    return false
  }
}

type InstallConfig = {
  registerUrl: string
}

type VerificationState = "idle" | "waiting" | "success" | "timeout"
type AgentSnapshot = {
  status: string
  connectedAt?: string | null
  lastHeartbeat?: string | null
  heartbeatUpdatedAt?: string | null
  createdAt?: string | null
}

const buildInstallCommands = (token: string, config: InstallConfig) => {
  if (!token) return { linux: "", windows: "" }
  const registerUrl = config.registerUrl.trim()
  if (!registerUrl) {
    return { linux: "", windows: "" }
  }
  const scriptBaseUrl = normalizeOrigin(registerUrl)
  const linux = `curl -kfsSL "${scriptBaseUrl}/api/agents/install.sh?token=${token}" | LUNAFOX_REGISTER_URL="${registerUrl}" bash`
  const windows = `[System.Net.ServicePointManager]::ServerCertificateValidationCallback = { $true }; $env:LUNAFOX_REGISTER_URL=\"${registerUrl}\"; irm "${scriptBaseUrl}/api/agents/install.ps1?token=${token}" | iex`
  return { linux, windows }
}

type AgentInstallDialogProps = {
  open: boolean
  token: RegistrationTokenResponse | null
  isGenerating: boolean
  onGenerate: () => void
}

export function AgentInstallDialog({
  open,
  token,
  isGenerating,
  onGenerate,
}: AgentInstallDialogProps) {
  const t = useTranslations("settings.workers")
  const tActions = useTranslations("common.actions")
  const tToast = useTranslations("toast")
  const formatRelativeTime = useFormatRelativeTime()
  const dialogRef = useRef<HTMLDivElement>(null)
  const safeServerUrl = typeof window !== "undefined" ? window.location.origin : FALLBACK_SERVER_URL
  const normalizedOrigin = normalizeOrigin(safeServerUrl)
  const defaultRegisterUrl = normalizedOrigin

  const [registerUrlInput, setRegisterUrlInput] = useState(defaultRegisterUrl)
  const [configOpen, setConfigOpen] = useState(false)
  const [step, setStep] = useState<1 | 2 | 3>(1)
  const [verificationState, setVerificationState] = useState<VerificationState>("idle")
  const [commandTab, setCommandTab] = useState<"linux" | "windows">("linux")
  const [copied, setCopied] = useState(false)
  const copyResetRef = useRef<number | null>(null)
  const verificationIntervalRef = useRef<number | null>(null)
  const verificationTimeoutRef = useRef<number | null>(null)
  const verificationBaselineRef = useRef<Map<number, AgentSnapshot>>(new Map())

  const hasToken = Boolean(token?.token)
  const tokenExpiresAt = useMemo(() => {
    if (!token?.expiresAt) return 0
    const timestamp = new Date(token.expiresAt).getTime()
    return Number.isNaN(timestamp) ? 0 : timestamp
  }, [token])
  const isTokenValid = hasToken && tokenExpiresAt > Date.now()

  useEffect(() => {
    if (!open) return
    setRegisterUrlInput(defaultRegisterUrl)
    setConfigOpen(false)
    setStep(1)
    setVerificationState("idle")
    setCommandTab("linux")
  }, [open, defaultRegisterUrl])

  const clearVerificationTimers = () => {
    if (verificationIntervalRef.current) {
      window.clearInterval(verificationIntervalRef.current)
      verificationIntervalRef.current = null
    }
    if (verificationTimeoutRef.current) {
      window.clearTimeout(verificationTimeoutRef.current)
      verificationTimeoutRef.current = null
    }
  }

  useEffect(() => {
    if (!open || step !== 3 || verificationState !== "waiting") return

    let active = true
    const startAt = Date.now()

    const toTimestamp = (value?: string | null) => {
      if (!value) return 0
      const timestamp = new Date(value).getTime()
      return Number.isNaN(timestamp) ? 0 : timestamp
    }

    const isAfterStart = (value?: string | null) => {
      const timestamp = toTimestamp(value)
      return timestamp > startAt
    }

    const buildSnapshot = (agents: Agent[]) => {
      const snapshot = new Map<number, AgentSnapshot>()
      agents.forEach((agent) => {
        snapshot.set(agent.id, {
          status: agent.status,
          connectedAt: agent.connectedAt ?? null,
          lastHeartbeat: agent.lastHeartbeat ?? null,
          heartbeatUpdatedAt: agent.heartbeat?.updatedAt ?? null,
          createdAt: agent.createdAt ?? null,
        })
      })
      return snapshot
    }

    const detectNewAgent = (agents: Agent[]) => {
      const baseline = verificationBaselineRef.current
      for (const agent of agents) {
        if (baseline.has(agent.id)) continue

        const recent =
          isAfterStart(agent.createdAt) ||
          isAfterStart(agent.connectedAt) ||
          isAfterStart(agent.lastHeartbeat) ||
          isAfterStart(agent.heartbeat?.updatedAt)
        if (!recent) continue

        if (
          agent.status !== "online" &&
          !isAfterStart(agent.lastHeartbeat) &&
          !isAfterStart(agent.heartbeat?.updatedAt)
        ) {
          continue
        }

        return true
      }
      return false
    }

    const fetchAgents = async (): Promise<Agent[]> => {
      const response = await agentService.getAgents(1, VERIFY_PAGE_SIZE)
      return response.results ?? []
    }

    const poll = async () => {
      try {
        const agents = await fetchAgents()
        if (!active || verificationState !== "waiting") return
        if (detectNewAgent(agents)) {
          setVerificationState("success")
          clearVerificationTimers()
        }
      } catch {
        // ignore network errors, keep polling until timeout
      }
    }

    const init = async () => {
      try {
        const agents = await fetchAgents()
        if (!active) return
        verificationBaselineRef.current = buildSnapshot(agents)
        await poll()
      } catch {
        verificationBaselineRef.current = new Map()
      }
    }

    clearVerificationTimers()
    init()
    verificationIntervalRef.current = window.setInterval(poll, VERIFY_POLL_INTERVAL_MS)
    verificationTimeoutRef.current = window.setTimeout(() => {
      if (!active) return
      setVerificationState("timeout")
      clearVerificationTimers()
    }, VERIFY_TIMEOUT_MS)

    return () => {
      active = false
      clearVerificationTimers()
    }
  }, [open, step, verificationState])

  useEffect(() => {
    if (copyResetRef.current) {
      window.clearTimeout(copyResetRef.current)
      copyResetRef.current = null
    }
    setCopied(false)
  }, [open])

  const copyToClipboard = useCallback(async (text: string): Promise<boolean> => {
    const value = text ?? ""
    if (navigator.clipboard?.writeText) {
      try {
        await navigator.clipboard.writeText(value)
        return true
      } catch {
        // fallback below
      }
    }

    try {
      const container = dialogRef.current || document.body
      const textArea = document.createElement("textarea")
      textArea.value = value
      textArea.style.position = "fixed"
      textArea.style.left = "-9999px"
      textArea.style.top = "-9999px"
      textArea.style.opacity = "0"
      container.appendChild(textArea)
      textArea.focus()
      textArea.select()
      textArea.setSelectionRange(0, textArea.value.length)
      const success = document.execCommand("copy")
      textArea.remove()
      return success
    } catch {
      return false
    }
  }, [])

  const handleCopy = useCallback(async (text: string) => {
    if (!text) {
      toast.error(tToast("copyFailed"), { id: COPY_TOAST_ID })
      return
    }
    toast.dismiss("agent-token")
    const success = await copyToClipboard(text)
    if (success) {
      if (copyResetRef.current) {
        window.clearTimeout(copyResetRef.current)
      }
      setCopied(true)
      copyResetRef.current = window.setTimeout(() => {
        setCopied(false)
        copyResetRef.current = null
      }, 2000)
      toast.success(tToast("copied"), { id: COPY_TOAST_ID })
    } else {
      setCopied(false)
      toast.error(tToast("copyFailed"), { id: COPY_TOAST_ID })
    }
  }, [copyToClipboard, tToast])

  const installCommands = useMemo(() => {
    return buildInstallCommands(token?.token ?? "", {
      registerUrl: registerUrlInput,
    })
  }, [token, registerUrlInput])

  const commandToCopy = useMemo(() => {
    return commandTab === "linux" ? installCommands.linux : installCommands.windows
  }, [commandTab, installCommands.linux, installCommands.windows])

  const canCopyCommand = useMemo(() => {
    return Boolean(commandToCopy)
  }, [commandToCopy])

  const showLocalOnlyWarning = useMemo(() => {
    const registerIsLocal = isLocalOrigin(normalizeOrigin(registerUrlInput))
    const registerUsesDockerHost = registerUrlInput.includes("server:8080")
    return registerIsLocal || registerUsesDockerHost
  }, [registerUrlInput])

  const canGoNext = useMemo(() => {
    if (step === 1) return isTokenValid && !isGenerating
    if (step === 2) return isTokenValid
    return true
  }, [step, isTokenValid, isGenerating])

  const goNext = () => {
    if (step === 2) {
      setVerificationState("waiting")
    }
    setStep((prev) => (prev < 3 ? ((prev + 1) as 1 | 2 | 3) : prev))
  }

  const goPrev = () => {
    setStep((prev) => (prev > 1 ? ((prev - 1) as 1 | 2 | 3) : prev))
  }

  return (
    <DialogContent ref={dialogRef} className="sm:max-w-[860px]">
      <DialogHeader>
        <DialogTitle>{t("install.title")}</DialogTitle>
        <DialogDescription>{t("install.desc")}</DialogDescription>
      </DialogHeader>
      <div className="space-y-4">
        <div className="flex flex-col gap-2 text-xs sm:flex-row sm:items-center sm:gap-3">
          <div className={cn("flex items-center gap-2", step >= 1 ? "text-foreground" : "text-muted-foreground")}>
            <span className={cn("w-5 h-5 rounded-full border flex items-center justify-center text-[10px] font-medium", step > 1 ? "border-emerald-500/40 text-emerald-600 bg-emerald-500/10" : step === 1 ? "border-blue-500/40 text-blue-600 bg-blue-500/10" : "border-border text-muted-foreground")}>
              {step > 1 ? "✓" : "1"}
            </span>
            <span>{t("install.step1Label")}</span>
          </div>
          <div className="hidden sm:block h-px flex-1 bg-border/70" />
          <div className={cn("flex items-center gap-2", step >= 2 ? "text-foreground" : "text-muted-foreground")}>
            <span className={cn("w-5 h-5 rounded-full border flex items-center justify-center text-[10px] font-medium", step > 2 ? "border-emerald-500/40 text-emerald-600 bg-emerald-500/10" : step === 2 ? "border-blue-500/40 text-blue-600 bg-blue-500/10" : "border-border text-muted-foreground")}>
              {step > 2 ? "✓" : "2"}
            </span>
            <span>{t("install.step2Label")}</span>
          </div>
          <div className="hidden sm:block h-px flex-1 bg-border/70" />
          <div className={cn("flex items-center gap-2", step >= 3 ? "text-foreground" : "text-muted-foreground")}>
            <span className={cn("w-5 h-5 rounded-full border flex items-center justify-center text-[10px] font-medium", step === 3 ? "border-blue-500/40 text-blue-600 bg-blue-500/10" : "border-border text-muted-foreground")}>
              3
            </span>
            <span>{t("install.step3Label")}</span>
          </div>
        </div>

        <div className={cn("rounded-lg border p-3 transition-all", isTokenValid ? "bg-background" : "bg-muted/20")}>
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex flex-col gap-1.5">
              <div className="flex items-center gap-2 text-xs">
                <Badge
                  variant="outline"
                  className={
                    !hasToken
                      ? "border-muted-foreground/30 text-muted-foreground"
                      : isTokenValid
                        ? "border-emerald-500/50 text-emerald-700 bg-emerald-500/5 dark:text-emerald-400"
                        : "border-rose-500/50 text-rose-700 bg-rose-500/5 dark:text-rose-400"
                  }
                >
                  <span
                    className={cn(
                      "mr-1.5 inline-block h-1.5 w-1.5 rounded-full",
                      isTokenValid ? "bg-emerald-500" : "bg-muted-foreground/40"
                    )}
                  />
                  {!hasToken
                    ? t("install.commandStatusIdle")
                    : isTokenValid
                      ? t("install.commandStatusReady")
                      : t("install.commandStatusExpired")}
                </Badge>
                {token?.expiresAt && isTokenValid && (
                  <span className="text-xs text-muted-foreground tabular-nums">
                    • {t("install.commandExpires", { time: formatRelativeTime(token.expiresAt) })}
                  </span>
                )}
              </div>
              {token?.expiresAt && isTokenValid && (
                <span className="text-[11px] text-muted-foreground">
                  {t("install.tokenUsage")}
                </span>
              )}
            </div>
            <Button size="sm" onClick={onGenerate} disabled={isGenerating} className="shrink-0">
              {isGenerating && <IconLoader2 className="mr-2 h-4 w-4 animate-spin" />}
              {token ? t("install.regenerateToken") : t("install.generateToken")}
            </Button>
          </div>
        </div>

        {step === 1 && (
          <div className="space-y-3">
            {!hasToken && (
              <div className="rounded-lg border bg-muted/20 p-3 space-y-1">
                <p className="text-xs font-medium text-muted-foreground">{t("install.tokenGuideTitle")}</p>
                <p className="text-xs text-muted-foreground">
                  {t("install.tokenGuideDesc")}
                </p>
              </div>
            )}

            <Collapsible open={configOpen} onOpenChange={setConfigOpen}>
              <div className="rounded-lg border bg-background p-3 space-y-3">
                <div className="flex items-center gap-2">
                  <CollapsibleTrigger asChild>
                    <Button variant="ghost" size="sm" className="justify-between text-xs h-8 flex-1">
                      <span className="text-sm font-medium">{t("install.configTitle")}</span>
                      <span className="flex items-center gap-1 text-muted-foreground">
                        {configOpen ? t("install.configHide") : t("install.configShow")}
                        {configOpen ? (
                          <IconChevronUp className="h-4 w-4" />
                        ) : (
                          <IconChevronDown className="h-4 w-4" />
                        )}
                      </span>
                    </Button>
                  </CollapsibleTrigger>
                </div>

                <div className="grid gap-2 text-xs text-muted-foreground">
                  <div className="flex flex-col gap-1 sm:flex-row sm:items-start sm:justify-between sm:gap-2">
                    <span className="text-[11px]">{t("install.serverUrl")}</span>
                    <span className="font-mono text-[11px] break-all sm:text-right">{registerUrlInput}</span>
                  </div>
                </div>

                <CollapsibleContent className="space-y-3">
                  <div className="grid gap-3 text-xs">
                    <div className="space-y-1">
                      <p className="text-[11px] text-muted-foreground">{t("install.presetDesc")}</p>
                      {showLocalOnlyWarning && (
                        <p className="text-[11px] text-amber-600">{t("install.localOnlyWarning")}</p>
                      )}
                    </div>
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <span className="text-muted-foreground">{t("install.serverUrl")}</span>
                      </div>
                      <Input
                        value={registerUrlInput}
                        onChange={(event) => {
                          setRegisterUrlInput(event.target.value)
                        }}
                        placeholder={t("install.serverUrlPlaceholder")}
                        className="h-8 font-mono text-xs"
                      />
                      <p className="text-[11px] text-muted-foreground">{t("install.serverUrlHint")}</p>
                    </div>
                  </div>

                  {!hasToken && (
                    <p className="text-xs text-muted-foreground">{t("install.tokenHint")}</p>
                  )}
                </CollapsibleContent>
              </div>
            </Collapsible>
          </div>
        )}

        {step === 3 && (
          <div className="space-y-3">
            <div className="rounded-lg border bg-background p-4 space-y-3">
              <div>
                <p className="text-sm font-medium">{t("install.verificationTitle")}</p>
                <p className="text-xs text-muted-foreground mt-1">{t("install.verificationDesc")}</p>
              </div>

              <div className="space-y-2">
                {verificationState === "idle" && (
                  <div className="rounded-lg border bg-muted/20 p-4 text-center">
                    <p className="text-sm text-muted-foreground">{t("install.verificationIdle")}</p>
                  </div>
                )}

                {verificationState === "waiting" && (
                  <div className="rounded-lg border bg-blue-500/5 p-4">
                    <div className="flex items-start gap-3">
                      <IconLoader2 className="h-5 w-5 animate-spin text-blue-600 shrink-0 mt-0.5" />
                      <div className="flex-1 space-y-2">
                        <p className="text-sm font-medium">{t("install.verificationWaiting")}</p>
                        <p className="text-xs text-muted-foreground">
                          {t("install.verificationWaitingDesc")}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {t("install.verificationWaitingTime")}
                        </p>
                      </div>
                    </div>
                  </div>
                )}

                {verificationState === "success" && (
                  <div className="rounded-lg border border-emerald-500/30 bg-emerald-500/5 p-4">
                    <div className="flex items-start gap-3">
                      <span className="mt-2 h-2 w-2 rounded-full bg-emerald-500 shrink-0" />
                      <div className="flex-1 space-y-2">
                        <p className="text-sm font-medium text-emerald-700 dark:text-emerald-400">{t("install.verificationSuccess")}</p>
                        <p className="text-xs text-muted-foreground">
                          {t("install.verificationSuccessDesc")}
                        </p>
                      </div>
                    </div>
                  </div>
                )}

                {verificationState === "timeout" && (
                  <div className="rounded-lg border border-amber-500/30 bg-amber-500/5 p-4">
                    <div className="flex items-start gap-3">
                      <span className="mt-2 h-2 w-2 rounded-full bg-amber-500 shrink-0" />
                      <div className="flex-1 space-y-2">
                        <p className="text-sm font-medium text-amber-700 dark:text-amber-400">{t("install.verificationTimeout")}</p>
                        <p className="text-xs text-muted-foreground">
                          {t("install.verificationTimeoutDesc")}
                        </p>
                        <ul className="text-xs text-muted-foreground space-y-1 list-disc list-inside">
                          <li>{t("install.verificationTimeoutCheck1")}</li>
                          <li>{t("install.verificationTimeoutCheck2")}</li>
                          <li>{t("install.verificationTimeoutCheck3")}</li>
                          <li>{t("install.verificationTimeoutCheck4")}</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>

            <div className="rounded-lg border bg-muted/20 p-4 text-xs">
              <div className="space-y-3">
                <div>
                  <div className="font-semibold text-foreground mb-2 flex items-center gap-1.5">
                    {t("install.verificationTips")}
                  </div>
                  <div className="space-y-1.5 text-muted-foreground">
                    <p>• {t("install.verificationTip1")}</p>
                    <p>• {t("install.verificationTip2")}</p>
                    <p>• {t("install.verificationTip3")}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {step === 2 && (
          <div className="space-y-3">
            <div className="rounded-lg border bg-background p-3 space-y-3">
              <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <p className="text-sm font-medium">{t("install.commandTitle")}</p>
                  <p className="text-xs text-muted-foreground">{t("install.commandDesc")}</p>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  className="h-8 text-xs shrink-0"
                  onClick={() => void handleCopy(commandToCopy)}
                  disabled={!canCopyCommand}
                  aria-label={`${tActions("copy")} ${t("install.commandTitle")}`}
                >
                  {copied ? tToast("copied") : tActions("copy")}
                </Button>
              </div>
              <Tabs
                value={commandTab}
                onValueChange={(value) => setCommandTab(value as "linux" | "windows")}
                className="w-full"
              >
                <TabsList className="grid w-full grid-cols-2">
                  <TabsTrigger value="linux">{t("install.linuxTab")}</TabsTrigger>
                  <TabsTrigger value="windows">{t("install.windowsTab")}</TabsTrigger>
                </TabsList>
                <TabsContent value="linux" className="space-y-2">
                  <div className="rounded-lg border bg-muted/30 p-3">
                    {token ? (
                      <pre className="max-h-48 overflow-y-auto font-mono text-xs whitespace-pre-wrap break-all">{installCommands.linux}</pre>
                    ) : (
                      <div className="text-xs text-muted-foreground">
                        {isGenerating ? t("install.commandStatusGenerating") : t("install.commandPlaceholder")}
                      </div>
                    )}
                  </div>
                </TabsContent>
                <TabsContent value="windows" className="space-y-2">
                  <div className="rounded-lg border bg-muted/30 p-3">
                    {token ? (
                      <pre className="max-h-48 overflow-y-auto font-mono text-xs whitespace-pre-wrap break-all">{installCommands.windows}</pre>
                    ) : (
                      <div className="text-xs text-muted-foreground">
                        {isGenerating ? t("install.commandStatusGenerating") : t("install.commandPlaceholder")}
                      </div>
                    )}
                  </div>
                </TabsContent>
              </Tabs>
            </div>
          </div>
        )}

        {step === 2 && (
          <div className="rounded-lg border bg-muted/20 p-4 text-xs">
            <div className="space-y-3">
              <div>
                <div className="font-semibold text-foreground mb-2 flex items-center gap-1.5">
                  {t("install.stepsTitle")}
                </div>
                <div className="space-y-1.5 text-muted-foreground">
                  <div className="flex items-start gap-2">
                    <span className="font-medium">1.</span>
                    <span>{t("install.step1Desc")}</span>
                  </div>
                  <div className="flex items-start gap-2">
                    <span className="font-medium">2.</span>
                    <span>{t("install.step2Desc")}</span>
                  </div>
                  <div className="flex items-start gap-2">
                    <span className="font-medium">3.</span>
                    <span>{t("install.step3Desc")}</span>
                  </div>
                </div>
              </div>
              <div className="pt-2 border-t border-border/50">
                <div className="font-semibold text-foreground mb-2 flex items-center gap-1.5">
                  {t("install.requirementsTitle")}
                </div>
                <div className="space-y-1 text-muted-foreground">
                  <div className="flex items-center gap-2">
                    <span>•</span>
                    <span>{t("install.requirementsDocker")}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span>•</span>
                    <span>{t("install.requirementsAccess")}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        <div className="flex items-center justify-between pt-1">
          <Button variant="ghost" size="sm" onClick={goPrev} disabled={step === 1}>
            {tActions("previous")}
          </Button>
          {step < 3 ? (
            <Button size="sm" onClick={goNext} disabled={!canGoNext}>
              {tActions("next")}
            </Button>
          ) : (
            <DialogClose asChild>
              <Button size="sm">{tActions("close")}</Button>
            </DialogClose>
          )}
        </div>
      </div>
    </DialogContent>
  )
}
