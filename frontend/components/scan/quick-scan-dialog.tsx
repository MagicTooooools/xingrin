"use client"

import * as React from "react"
import { useTranslations } from "next-intl"
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog"
import { EnginePresetSelector } from "./engine-preset-selector"
import { ScanConfigEditor } from "./scan-config-editor"
import { useQuickScanDialogState } from "@/components/scan/quick-scan-dialog-state"
import {
  QuickScanFooter,
  QuickScanHeader,
  QuickScanOverwriteDialog,
  QuickScanTargetStep,
  QuickScanTrigger,
} from "@/components/scan/quick-scan-dialog-sections"

interface QuickScanDialogProps {
  trigger?: React.ReactNode
}

export function QuickScanDialog({ trigger }: QuickScanDialogProps) {
  const t = useTranslations("quickScan")
  const {
    open,
    handleClose,
    isSubmitting,
    step,
    targetInput,
    setTargetInput,
    selectedEngineIds,
    selectedPresetId,
    setSelectedPresetId,
    configuration,
    isConfigEdited,
    showOverwriteConfirm,
    setShowOverwriteConfirm,
    lineNumbersRef,
    handleTextareaScroll,
    validInputs,
    invalidInputs,
    hasErrors,
    engines,
    selectedEngines,
    handlePresetConfigChange,
    handleManualConfigChange,
    handleEngineIdsChange,
    handleOverwriteConfirm,
    handleOverwriteCancel,
    handleYamlValidationChange,
    canProceedToStep2,
    canProceedToStep3,
    canSubmit,
    handleNext,
    handleBack,
    handleSubmit,
    totalSteps,
  } = useQuickScanDialogState({ t })
  
  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogTrigger asChild>
        <QuickScanTrigger t={t} trigger={trigger} />
      </DialogTrigger>
      <DialogContent className="max-w-[90vw] sm:max-w-[900px] p-0 gap-0">
        <QuickScanHeader t={t} step={step} totalSteps={totalSteps} />

        <div className="border-t h-[480px] overflow-hidden">
          {/* Step 1: Target input */}
          {step === 1 && (
            <QuickScanTargetStep
              t={t}
              targetInput={targetInput}
              onTargetChange={setTargetInput}
              lineNumbersRef={lineNumbersRef}
              onScroll={handleTextareaScroll}
              invalidInputs={invalidInputs}
              hasErrors={hasErrors}
            />
          )}

          {/* Step 2: Select preset/engines */}
          {step === 2 && engines && (
            <EnginePresetSelector
              engines={engines}
              selectedEngineIds={selectedEngineIds}
              selectedPresetId={selectedPresetId}
              onPresetChange={setSelectedPresetId}
              onEngineIdsChange={handleEngineIdsChange}
              onConfigurationChange={handlePresetConfigChange}
              disabled={isSubmitting}
            />
          )}

          {/* Step 3: Edit configuration */}
          {step === 3 && (
            <ScanConfigEditor
              configuration={configuration}
              onChange={handleManualConfigChange}
              onValidationChange={handleYamlValidationChange}
              selectedEngines={selectedEngines}
              isConfigEdited={isConfigEdited}
              disabled={isSubmitting}
            />
          )}
        </div>
        <QuickScanFooter
          t={t}
          step={step}
          validCount={validInputs.length}
          invalidCount={invalidInputs.length}
          selectedEngineCount={selectedEngineIds.length}
          isSubmitting={isSubmitting}
          canProceedToStep2={canProceedToStep2}
          canProceedToStep3={canProceedToStep3}
          canSubmit={canSubmit}
          onBack={handleBack}
          onNext={handleNext}
          onSubmit={handleSubmit}
        />
      </DialogContent>
      
      {/* Overwrite confirmation dialog */}
      <QuickScanOverwriteDialog
        open={showOverwriteConfirm}
        onOpenChange={setShowOverwriteConfirm}
        onCancel={handleOverwriteCancel}
        onConfirm={handleOverwriteConfirm}
        t={t}
      />
    </Dialog>
  )
}
