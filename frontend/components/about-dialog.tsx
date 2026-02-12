"use client"

import type { ReactNode } from "react"
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog"
import { useAboutDialogState } from "@/components/about-dialog-state"
import {
  AboutDialogHeader,
  AboutDialogBranding,
  AboutDialogVersionInfo,
  AboutDialogLinks,
  AboutDialogFooter,
  Separator,
} from "@/components/about-dialog-sections"

interface AboutDialogProps {
  children: ReactNode
}

export function AboutDialog({ children }: AboutDialogProps) {
  const {
    t,
    isChecking,
    updateResult,
    checkError,
    currentVersion,
    latestVersion,
    hasUpdate,
    releaseUrl,
    logoSrc,
    handleCheckUpdate,
  } = useAboutDialogState()

  return (
    <Dialog>
      <DialogTrigger asChild>
        {children}
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <AboutDialogHeader t={t} />

        <div className="space-y-6">
          <AboutDialogBranding t={t} logoSrc={logoSrc} />

          <AboutDialogVersionInfo
            t={t}
            currentVersion={currentVersion}
            latestVersion={latestVersion}
            hasUpdate={hasUpdate}
            checkError={checkError}
            isChecking={isChecking}
            releaseUrl={releaseUrl}
            showLatest={!!updateResult}
            onCheckUpdate={handleCheckUpdate}
          />

          <Separator />

          <AboutDialogLinks t={t} />

          <AboutDialogFooter t={t} />
        </div>
      </DialogContent>
    </Dialog>
  )
}
