import React from "react"
import { IconChevronDown, IconChevronUp, IconLoader2 } from "@/components/icons"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { DialogClose } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible"
import { cn } from "@/lib/utils"
import type { RegistrationTokenResponse } from "@/types/agent.types"

type TranslationFn = (key: string, params?: Record<string, string | number | Date>) => string

type StepValue = 1 | 2 | 3

interface AgentInstallStepIndicatorProps {
  t: TranslationFn
  step: StepValue
}

export function AgentInstallStepIndicator({ t, step }: AgentInstallStepIndicatorProps) {
  return (
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
  )
}

interface AgentInstallTokenCardProps {
  t: TranslationFn
  token: RegistrationTokenResponse | null
  hasToken: boolean
  isTokenValid: boolean
  isGenerating: boolean
  formatRelativeTime: (value: string | Date) => string
  onGenerate: () => void
}

export function AgentInstallTokenCard({
  t,
  token,
  hasToken,
  isTokenValid,
  isGenerating,
  formatRelativeTime,
  onGenerate,
}: AgentInstallTokenCardProps) {
  return (
    <div className={cn("rounded-lg border p-3 transition-[background-color,border-color]", isTokenValid ? "bg-background" : "bg-muted/20")}>
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
  )
}

interface AgentInstallConfigPanelProps {
  t: TranslationFn
  registerUrlInput: string
  setRegisterUrlInput: (value: string) => void
  configOpen: boolean
  setConfigOpen: (open: boolean) => void
  showLocalOnlyWarning: boolean
  hasToken: boolean
}

export function AgentInstallConfigPanel({
  t,
  registerUrlInput,
  setRegisterUrlInput,
  configOpen,
  setConfigOpen,
  showLocalOnlyWarning,
  hasToken,
}: AgentInstallConfigPanelProps) {
  return (
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
                type="url"
                name="serverUrl"
                autoComplete="url"
                inputMode="url"
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
  )
}

interface AgentInstallCommandPanelProps {
  t: TranslationFn
  tActions: TranslationFn
  tToast: TranslationFn
  token: RegistrationTokenResponse | null
  installCommand: string
  isGenerating: boolean
  canCopyCommand: boolean
  copied: boolean
  onCopy: (value: string) => void
}

export function AgentInstallCommandPanel({
  t,
  tActions,
  tToast,
  token,
  installCommand,
  isGenerating,
  canCopyCommand,
  copied,
  onCopy,
}: AgentInstallCommandPanelProps) {
  return (
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
          onClick={() => onCopy(installCommand)}
          disabled={!canCopyCommand}
          aria-label={`${tActions("copy")} ${t("install.commandTitle")}`}
        >
          {copied ? tToast("copied") : tActions("copy")}
        </Button>
      </div>
      <div className="rounded-lg border bg-muted/30 p-3">
        {token ? (
          <pre className="max-h-48 overflow-y-auto font-mono text-xs whitespace-pre-wrap break-all">{installCommand}</pre>
        ) : (
          <div className="text-xs text-muted-foreground">
            {isGenerating ? t("install.commandStatusGenerating") : t("install.commandPlaceholder")}
          </div>
        )}
      </div>
    </div>
  )
}

interface AgentInstallCommandTipsProps {
  t: TranslationFn
}

export function AgentInstallCommandTips({ t }: AgentInstallCommandTipsProps) {
  return (
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
  )
}

interface AgentInstallVerificationPanelProps {
  t: TranslationFn
  verificationState: "idle" | "waiting" | "success" | "timeout"
}

export function AgentInstallVerificationPanel({
  t,
  verificationState,
}: AgentInstallVerificationPanelProps) {
  return (
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
  )
}

interface AgentInstallVerificationTipsProps {
  t: TranslationFn
}

export function AgentInstallVerificationTips({ t }: AgentInstallVerificationTipsProps) {
  return (
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
  )
}

interface AgentInstallFooterProps {
  tActions: TranslationFn
  step: StepValue
  canGoNext: boolean
  onPrev: () => void
  onNext: () => void
}

export function AgentInstallFooter({
  tActions,
  step,
  canGoNext,
  onPrev,
  onNext,
}: AgentInstallFooterProps) {
  return (
    <div className="flex items-center justify-between pt-1">
      <Button variant="ghost" size="sm" onClick={onPrev} disabled={step === 1}>
        {tActions("previous")}
      </Button>
      {step < 3 ? (
        <Button size="sm" onClick={onNext} disabled={!canGoNext}>
          {tActions("next")}
        </Button>
      ) : (
        <DialogClose asChild>
          <Button size="sm">{tActions("close")}</Button>
        </DialogClose>
      )}
    </div>
  )
}
