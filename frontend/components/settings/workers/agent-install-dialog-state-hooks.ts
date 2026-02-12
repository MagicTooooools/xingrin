import { useCallback, useEffect, useRef, useState } from "react"
import { toast } from "sonner"
import { agentService } from "@/services/agent.service"
import {
  buildAgentSnapshot,
  detectNewAgent,
  type AgentSnapshot,
} from "@/lib/agent-install-helpers"
import type { Agent } from "@/types/agent.types"

const COPY_TOAST_ID = "install-command-copy"
const VERIFY_POLL_INTERVAL_MS = 5000
const VERIFY_TIMEOUT_MS = 60_000
const VERIFY_PAGE_SIZE = 200

export type VerificationState = "idle" | "waiting" | "success" | "timeout"

type UseAgentVerificationProps = {
  open: boolean
  step: 1 | 2 | 3
  verificationState: VerificationState
  setVerificationState: (state: VerificationState) => void
}

export function useAgentVerification({
  open,
  step,
  verificationState,
  setVerificationState,
}: UseAgentVerificationProps) {
  const verificationIntervalRef = useRef<number | null>(null)
  const verificationTimeoutRef = useRef<number | null>(null)
  const verificationBaselineRef = useRef<Map<number, AgentSnapshot>>(new Map())

  const clearVerificationTimers = useCallback(() => {
    if (verificationIntervalRef.current) {
      window.clearInterval(verificationIntervalRef.current)
      verificationIntervalRef.current = null
    }
    if (verificationTimeoutRef.current) {
      window.clearTimeout(verificationTimeoutRef.current)
      verificationTimeoutRef.current = null
    }
  }, [])

  useEffect(() => {
    if (!open || step !== 3 || verificationState !== "waiting") return

    let active = true
    const startAt = Date.now()

    const fetchAgents = async (): Promise<Agent[]> => {
      const response = await agentService.getAgents(1, VERIFY_PAGE_SIZE)
      return response.results ?? []
    }

    const poll = async () => {
      try {
        const agents = await fetchAgents()
        if (!active || verificationState !== "waiting") return
        if (detectNewAgent(agents, verificationBaselineRef.current, startAt)) {
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
        verificationBaselineRef.current = buildAgentSnapshot(agents)
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
  }, [clearVerificationTimers, open, step, verificationState, setVerificationState])
}

type UseInstallCommandCopyProps = {
  open: boolean
  dialogRef: React.RefObject<HTMLDivElement | null>
  tToast: (key: string, params?: Record<string, string | number | Date>) => string
}

export function useInstallCommandCopy({
  open,
  dialogRef,
  tToast,
}: UseInstallCommandCopyProps) {
  const [copied, setCopied] = useState(false)
  const copyResetRef = useRef<number | null>(null)

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
  }, [dialogRef])

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

  return {
    copied,
    handleCopy,
  }
}
