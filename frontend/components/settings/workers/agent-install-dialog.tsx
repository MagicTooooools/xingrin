"use client"

import { useAgentInstallDialogState } from "@/components/settings/workers/agent-install-dialog-state"
import { useTranslations } from "next-intl"
import { DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { useFormatRelativeTime } from "@/lib/i18n-format"
import type { RegistrationTokenResponse } from "@/types/agent.types"
import {
  AgentInstallCommandPanel,
  AgentInstallCommandTips,
  AgentInstallConfigPanel,
  AgentInstallFooter,
  AgentInstallStepIndicator,
  AgentInstallTokenCard,
  AgentInstallVerificationPanel,
  AgentInstallVerificationTips,
} from "@/components/settings/workers/agent-install-dialog-sections"

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
  const {
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
  } = useAgentInstallDialogState({
    open,
    token,
    isGenerating,
    tToast,
  })

  return (
    <DialogContent ref={dialogRef} className="sm:max-w-[860px]">
      <DialogHeader>
        <DialogTitle>{t("install.title")}</DialogTitle>
        <DialogDescription>{t("install.desc")}</DialogDescription>
      </DialogHeader>
      <div className="space-y-4">
        <AgentInstallStepIndicator t={t} step={step} />

        <AgentInstallTokenCard
          t={t}
          token={token}
          hasToken={hasToken}
          isTokenValid={isTokenValid}
          isGenerating={isGenerating}
          formatRelativeTime={formatRelativeTime}
          onGenerate={onGenerate}
        />

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
            <AgentInstallConfigPanel
              t={t}
              registerUrlInput={registerUrlInput}
              setRegisterUrlInput={setRegisterUrlInput}
              configOpen={configOpen}
              setConfigOpen={setConfigOpen}
              showLocalOnlyWarning={showLocalOnlyWarning}
              hasToken={hasToken}
            />
          </div>
        )}

        {step === 3 && (
          <div className="space-y-3">
            <AgentInstallVerificationPanel t={t} verificationState={verificationState} />
            <AgentInstallVerificationTips t={t} />
          </div>
        )}

        {step === 2 && (
          <div className="space-y-3">
            <AgentInstallCommandPanel
              t={t}
              tActions={tActions}
              tToast={tToast}
              token={token}
              installCommand={installCommand}
              isGenerating={isGenerating}
              canCopyCommand={canCopyCommand}
              copied={copied}
              onCopy={handleCopy}
            />
          </div>
        )}

        {step === 2 && (
          <AgentInstallCommandTips t={t} />
        )}

        <AgentInstallFooter
          tActions={tActions}
          step={step}
          canGoNext={canGoNext}
          onPrev={goPrev}
          onNext={goNext}
        />
      </div>
    </DialogContent>
  )
}
