import { useCallback, useEffect, useMemo, useRef, useState } from "react"
import {
  buildInstallCommand,
  isLocalOrigin,
  normalizeOrigin,
} from "@/lib/agent-install-helpers"
import type { RegistrationTokenResponse } from "@/types/agent.types"
import {
  useAgentVerification,
  useInstallCommandCopy,
  type VerificationState,
} from "@/components/settings/workers/agent-install-dialog-state-hooks"

const FALLBACK_SERVER_URL = "https://your-lunafox-server"
type UseAgentInstallDialogStateProps = {
  open: boolean
  token: RegistrationTokenResponse | null
  isGenerating: boolean
  tToast: (key: string, params?: Record<string, string | number | Date>) => string
}

export function useAgentInstallDialogState({
  open,
  token,
  isGenerating,
  tToast,
}: UseAgentInstallDialogStateProps) {
  const dialogRef = useRef<HTMLDivElement>(null)
  const safeServerUrl = typeof window !== "undefined" ? window.location.origin : FALLBACK_SERVER_URL
  const normalizedOrigin = normalizeOrigin(safeServerUrl)
  const defaultRegisterUrl = normalizedOrigin

  const [registerUrlInput, setRegisterUrlInput] = useState(defaultRegisterUrl)
  const [configOpen, setConfigOpen] = useState(false)
  const [step, setStep] = useState<1 | 2 | 3>(1)
  const [verificationState, setVerificationState] = useState<VerificationState>("idle")
  const { copied, handleCopy } = useInstallCommandCopy({
    open,
    dialogRef,
    tToast,
  })

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
  }, [open, defaultRegisterUrl])

  useAgentVerification({
    open,
    step,
    verificationState,
    setVerificationState,
  })

  const installCommand = useMemo(() => {
    return buildInstallCommand(token?.token ?? "", registerUrlInput)
  }, [token, registerUrlInput])

  const commandToCopy = installCommand

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

  const goNext = useCallback(() => {
    if (step === 2) {
      setVerificationState("waiting")
    }
    setStep((prev) => (prev < 3 ? ((prev + 1) as 1 | 2 | 3) : prev))
  }, [step])

  const goPrev = useCallback(() => {
    setStep((prev) => (prev > 1 ? ((prev - 1) as 1 | 2 | 3) : prev))
  }, [])

  return {
    dialogRef,
    registerUrlInput,
    setRegisterUrlInput,
    configOpen,
    setConfigOpen,
    step,
    verificationState,
    copied,
    hasToken,
    tokenExpiresAt,
    isTokenValid,
    handleCopy,
    installCommand,
    canCopyCommand,
    showLocalOnlyWarning,
    canGoNext,
    goNext,
    goPrev,
    setStep,
    setVerificationState,
    defaultRegisterUrl,
  }
}
