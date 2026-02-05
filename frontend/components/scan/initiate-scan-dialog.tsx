"use client"

import React, { useState, useMemo, useCallback } from "react"
import { AlertTriangle, CheckCircle2, ChevronRight, Play, Settings } from "@/components/icons"
import { useTranslations } from "next-intl"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { LoadingSpinner } from "@/components/loading-spinner"
import { ScanConfigEditor } from "./scan-config-editor"
import { cn } from "@/lib/utils"
import { mergeEngineConfigurations } from "@/lib/engine-config"

import type { Organization } from "@/types/organization.types"

import { initiateScan } from "@/services/scan.service"
import { toast } from "sonner"
import { useEngines, usePresetEngines } from "@/hooks/use-engines"
import { useQueryClient } from "@tanstack/react-query"

interface InitiateScanDialogProps {
  organization?: Organization | null
  organizationId?: number
  targetId?: number
  targetName?: string
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess?: () => void
}

export function InitiateScanDialog({
  organization,
  organizationId,
  targetId,
  targetName,
  open,
  onOpenChange,
  onSuccess,
}: InitiateScanDialogProps) {
  const t = useTranslations("scan.initiate")
  const tToast = useTranslations("toast")
  const queryClient = useQueryClient()
  const [selectedEngineIds, setSelectedEngineIds] = useState<number[]>([])
  const [selectedPresetId, setSelectedPresetId] = useState<string | null>(null)
  const [selectMode, setSelectMode] = useState<"preset" | "custom">("preset")
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [currentStep, setCurrentStep] = useState(1)
  
  // Configuration state management
  const [configuration, setConfiguration] = useState("")
  const [isConfigEdited, setIsConfigEdited] = useState(false)
  const [isYamlValid, setIsYamlValid] = useState(true)
  const [showOverwriteConfirm, setShowOverwriteConfirm] = useState(false)
  const [pendingConfigChange, setPendingConfigChange] = useState<string | null>(null)
  const [pendingEngineIds, setPendingEngineIds] = useState<number[] | null>(null)
  const [pendingPresetId, setPendingPresetId] = useState<string | null>(null)

  const { data: engines, isLoading: isLoadingEngines, isError: isEnginesError } = useEngines()
  const { data: presetEngines, isLoading: isLoadingPresets, isError: isPresetsError } = usePresetEngines()

  const selectedEngines = useMemo(() => {
    if (!selectedEngineIds.length || !engines) return []
    return engines.filter((e) => selectedEngineIds.includes(e.id))
  }, [selectedEngineIds, engines])

  const selectedPreset = useMemo(() => {
    if (!presetEngines || !selectedPresetId) return null
    return presetEngines.find((preset) => preset.id === selectedPresetId) || null
  }, [presetEngines, selectedPresetId])

  // Handle manual config editing
  const handleManualConfigChange = useCallback((value: string) => {
    setConfiguration(value)
    setIsConfigEdited(true)
  }, [])

  const buildConfigFromEngines = useCallback((engineIds: number[]) => {
    if (!engines) return ""
    const selected = engines.filter((e) => engineIds.includes(e.id))
    return mergeEngineConfigurations(selected.map((e) => e.configuration || ""))
  }, [engines])

  const applyEngineSelection = useCallback((engineIds: number[], nextConfig: string) => {
    setSelectedEngineIds(engineIds)
    setConfiguration(nextConfig)
    setIsConfigEdited(false)
    setIsYamlValid(true)
  }, [])

  const handleEngineIdsChange = useCallback((engineIds: number[]) => {
    const nextConfig = buildConfigFromEngines(engineIds)
    if (isConfigEdited && configuration !== nextConfig) {
      setPendingEngineIds(engineIds)
      setPendingConfigChange(nextConfig)
      setPendingPresetId(null)
      setShowOverwriteConfirm(true)
      return
    }
    applyEngineSelection(engineIds, nextConfig)
    setSelectedPresetId(null)
  }, [applyEngineSelection, buildConfigFromEngines, configuration, isConfigEdited])

  const handlePresetSelect = useCallback((presetId: string, presetConfig: string) => {
    if (isConfigEdited && configuration !== presetConfig) {
      setPendingEngineIds(null)
      setPendingConfigChange(presetConfig)
      setPendingPresetId(presetId)
      setShowOverwriteConfirm(true)
      return
    }
    setSelectedPresetId(presetId)
    setConfiguration(presetConfig)
    setIsConfigEdited(false)
    setIsYamlValid(true)
  }, [configuration, isConfigEdited])

  const handleOverwriteConfirm = () => {
    if (pendingConfigChange !== null) {
      if (pendingPresetId !== null) {
        setSelectedPresetId(pendingPresetId)
        setConfiguration(pendingConfigChange)
        setIsConfigEdited(false)
        setIsYamlValid(true)
      } else {
        const nextEngineIds = pendingEngineIds ?? selectedEngineIds
        applyEngineSelection(nextEngineIds, pendingConfigChange)
      }
    }
    setShowOverwriteConfirm(false)
    setPendingConfigChange(null)
    setPendingEngineIds(null)
    setPendingPresetId(null)
  }

  const handleOverwriteCancel = () => {
    setShowOverwriteConfirm(false)
    setPendingConfigChange(null)
    setPendingEngineIds(null)
    setPendingPresetId(null)
  }

  const handleYamlValidationChange = (isValid: boolean) => {
    setIsYamlValid(isValid)
  }

  const handleInitiate = async () => {
    if (selectMode === "preset") {
      if (!selectedPresetId) {
        toast.error(tToast("noPresetSelected"))
        return
      }
    } else if (selectedEngineIds.length === 0) {
      toast.error(tToast("noEngineSelected"))
      return
    }
    if (!configuration.trim()) {
      toast.error(tToast("emptyConfig"))
      return
    }
    if (!isYamlValid) {
      toast.error(tToast("invalidConfig"))
      return
    }
    
    if (!organizationId && !targetId) {
      toast.error(tToast("paramError"), { description: tToast("paramErrorDesc") })
      return
    }
    
    setIsSubmitting(true)
    try {
      const engineIds = selectMode === "custom" ? selectedEngineIds : []
      const engineNames = selectMode === "custom"
        ? selectedEngines.slice(0, 1).map((e) => e.name)
        : [selectedPresetId as string]

      const response = await initiateScan({
        organizationId,
        targetId,
        configuration,
        engineIds,
        engineNames,
      })
      
      // 后端返回 201 说明成功创建扫描任务
      const scanCount = response.scans?.length || response.count || 0
      toast.success(tToast("scanInitiated"), {
        description: response.message || tToast("scanInitiatedDesc", { count: scanCount }),
      })
      queryClient.invalidateQueries({ queryKey: ["scans"] })
      queryClient.invalidateQueries({ queryKey: ["scan-statistics"] })
      onSuccess?.()
      onOpenChange(false)
      setSelectedEngineIds([])
      setSelectedPresetId(null)
      setConfiguration("")
      setIsConfigEdited(false)
      setIsYamlValid(true)
      setCurrentStep(1)
      setSelectMode("preset")
      setShowOverwriteConfirm(false)
      setPendingConfigChange(null)
      setPendingEngineIds(null)
      setPendingPresetId(null)
    } catch (err: unknown) {
      const error = err as { response?: { data?: { error?: { code?: string; message?: string } } } }
      toast.error(tToast("initiateScanFailed"), {
        description: error?.response?.data?.error?.message || (err instanceof Error ? err.message : tToast("unknownError")),
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!isSubmitting) {
      onOpenChange(newOpen)
      if (!newOpen) {
        setSelectedEngineIds([])
        setSelectedPresetId(null)
        setConfiguration("")
        setIsConfigEdited(false)
        setIsYamlValid(true)
        setCurrentStep(1)
        setSelectMode("preset")
        setShowOverwriteConfirm(false)
        setPendingConfigChange(null)
        setPendingEngineIds(null)
        setPendingPresetId(null)
      }
    }
  }

  const steps = [
    { id: 1, title: t("steps.selectEngine") },
    { id: 2, title: t("steps.editConfig") },
  ]

  const hasConfig = configuration.trim().length > 0
  const canProceedToReview = selectMode === "preset"
    ? !!selectedPresetId
    : selectedEngineIds.length > 0
  // Check if normal start is available
  const canStart = configuration.trim().length > 0 &&
                   isYamlValid &&
                   (selectMode === "preset" ? !!selectedPresetId : selectedEngineIds.length > 0)

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-[90vw] sm:max-w-[900px] p-0 gap-0 flex flex-col max-h-[90vh]">
        <DialogHeader className="px-6 pt-6 pb-4 shrink-0">
          <DialogTitle className="flex items-center gap-2">
            <Play className="h-5 w-5" />
            {t("title")}
          </DialogTitle>
          <DialogDescription className="mt-1">
            {targetName ? (
              <>{t("targetDesc")} <span className="font-medium text-foreground">{targetName}</span></>
            ) : (
              <>{t("orgDesc")} <span className="font-medium text-foreground">{organization?.name}</span></>
            )}
          </DialogDescription>
        </DialogHeader>

        {/* Scrollable content area */}
        <div className="flex-1 overflow-y-auto border-t">
          <>
              <div className="px-6 pt-5 pb-2">
                <div className="flex flex-wrap items-center gap-3 text-sm">
                  {steps.map((step, index) => (
                    <React.Fragment key={step.id}>
                      <div className={cn(
                        "flex items-center gap-2",
                        currentStep === step.id ? "text-foreground" : "text-muted-foreground"
                      )}>
                        <span
                          className={cn(
                            "inline-flex h-7 w-7 items-center justify-center rounded-full border text-xs font-medium",
                            currentStep === step.id
                              ? "border-primary/40 bg-primary/10 text-primary"
                              : "border-muted-foreground/30 text-muted-foreground"
                          )}
                        >
                          {step.id}
                        </span>
                        <span className={currentStep === step.id ? "font-medium" : undefined}>
                          {step.title}
                        </span>
                      </div>
                      {index < steps.length - 1 && (
                        <ChevronRight className="h-4 w-4 text-muted-foreground" />
                      )}
                    </React.Fragment>
                  ))}
                </div>
              </div>

              {currentStep === 1 && (
                <div className="p-6 space-y-6">
                  <Tabs value={selectMode} onValueChange={(value) => setSelectMode(value as "preset" | "custom")}>
                    <TabsList>
                      <TabsTrigger value="preset">{t("mode.preset")}</TabsTrigger>
                      <TabsTrigger value="custom">{t("mode.custom")}</TabsTrigger>
                    </TabsList>
                    <TabsContent value="preset" className="mt-4 space-y-3">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm font-medium">{t("presets.title")}</p>
                          <p className="text-xs text-muted-foreground">{t("presets.selectHint")}</p>
                        </div>
                        {selectedPreset && (
                          <Badge variant="secondary" className="text-xs">
                            {selectedPreset.name}
                          </Badge>
                        )}
                      </div>
                      {isLoadingPresets ? (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <LoadingSpinner />
                          {t("presets.loading")}
                        </div>
                      ) : isPresetsError ? (
                        <div className="text-sm text-destructive">{t("loadFailed")}</div>
                      ) : presetEngines && presetEngines.length > 0 ? (
                        <div className="grid gap-2 sm:grid-cols-2">
                          {presetEngines.map((preset) => {
                            const isSelected = preset.id === selectedPresetId
                            return (
                              <button
                                key={preset.id}
                                type="button"
                                onClick={() => handlePresetSelect(preset.id, preset.configuration || "")}
                                disabled={isSubmitting}
                                className={cn(
                                  "flex flex-col items-start gap-2 rounded-lg border px-3 py-2 text-left transition-all",
                                  isSelected
                                    ? "border-primary/50 bg-primary/5"
                                    : "border-border hover:border-primary/40 hover:bg-muted/30",
                                  isSubmitting && "opacity-60 cursor-not-allowed"
                                )}
                              >
                                <span className="text-sm font-medium">{preset.name}</span>
                                {preset.description && (
                                  <span className="text-xs text-muted-foreground line-clamp-2">
                                    {preset.description}
                                  </span>
                                )}
                              </button>
                            )
                          })}
                        </div>
                      ) : (
                        <div className="text-sm text-muted-foreground">{t("presets.empty")}</div>
                      )}
                    </TabsContent>
                    <TabsContent value="custom" className="mt-4 space-y-3">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm font-medium">{t("selectEngineTitle")}</p>
                          <p className="text-xs text-muted-foreground">{t("selectEngineHint")}</p>
                        </div>
                        {selectedEngineIds.length > 0 && (
                          <Badge variant="secondary" className="text-xs">
                            {t("selectedCount", { count: selectedEngineIds.length })}
                          </Badge>
                        )}
                      </div>
                      {isLoadingEngines ? (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <LoadingSpinner />
                          {t("loading")}
                        </div>
                      ) : isEnginesError ? (
                        <div className="text-sm text-destructive">{t("loadFailed")}</div>
                      ) : engines && engines.length > 0 ? (
                        <div className="grid gap-2 sm:grid-cols-2">
                          {engines.map((engine) => {
                            const isSelected = selectedEngineIds.includes(engine.id)
                            return (
                              <label
                                key={engine.id}
                                htmlFor={`initiate-engine-${engine.id}`}
                                className={cn(
                                  "flex items-center gap-3 rounded-lg border px-3 py-2 cursor-pointer transition-all",
                                  isSelected
                                    ? "border-primary/50 bg-primary/5"
                                    : "border-border hover:border-primary/40 hover:bg-muted/30"
                                )}
                              >
                                <Checkbox
                                  id={`initiate-engine-${engine.id}`}
                                  checked={isSelected}
                                  onCheckedChange={(checked) => {
                                    const nextIds = checked ? [engine.id] : []
                                    handleEngineIdsChange(nextIds)
                                  }}
                                  disabled={isSubmitting}
                                  className="h-4 w-4"
                                />
                                <div className="min-w-0">
                                  <p className="text-sm font-medium truncate">{engine.name}</p>
                                  <p className="text-xs text-muted-foreground truncate">
                                    {engine.configuration ? t("configTitle") : t("noConfig")}
                                  </p>
                                </div>
                              </label>
                            )
                          })}
                        </div>
                      ) : (
                        <div className="text-sm text-muted-foreground">{t("noEngines")}</div>
                      )}
                    </TabsContent>
                  </Tabs>
                </div>
              )}

              {currentStep === 2 && (
                <div className="grid gap-6 p-6 lg:grid-cols-[1.1fr_0.9fr]">
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2 text-sm">
                        <Settings className="h-4 w-4 text-muted-foreground" />
                        <span className="font-medium">{t("configTitle")}</span>
                      </div>
                      {isConfigEdited && (
                        <Badge variant="outline" className="text-xs">
                          {t("configEdited")}
                        </Badge>
                      )}
                    </div>
                    <div className="border rounded-lg overflow-hidden">
                      <ScanConfigEditor
                        configuration={configuration}
                        onChange={handleManualConfigChange}
                        onValidationChange={handleYamlValidationChange}
                        selectedEngines={selectedEngines}
                        isConfigEdited={isConfigEdited}
                        disabled={isSubmitting}
                        className="h-[420px]"
                      />
                    </div>
                  </div>

                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium">{t("validation.title")}</span>
                      <Badge variant={isYamlValid && hasConfig ? "secondary" : "outline"} className="text-xs">
                        {isYamlValid && hasConfig ? t("validation.yamlOk") : t("validation.yamlError")}
                      </Badge>
                    </div>
                    <div className="rounded-lg border bg-muted/20 p-4 space-y-3 text-sm">
                      <div className="flex items-start gap-2">
                        {isYamlValid && hasConfig ? (
                          <CheckCircle2 className="h-4 w-4 text-emerald-500 mt-0.5" />
                        ) : (
                          <AlertTriangle className="h-4 w-4 text-amber-500 mt-0.5" />
                        )}
                        <span>{isYamlValid && hasConfig ? t("validation.yamlOk") : t("validation.yamlError")}</span>
                      </div>
                      {selectMode === "preset" ? (
                        <div className="flex items-start gap-2">
                          {selectedPreset ? (
                            <CheckCircle2 className="h-4 w-4 text-emerald-500 mt-0.5" />
                          ) : (
                            <AlertTriangle className="h-4 w-4 text-amber-500 mt-0.5" />
                          )}
                          <span>
                            {selectedPreset
                              ? t("validation.presetOk", { name: selectedPreset.name })
                              : t("validation.presetMissing")}
                          </span>
                        </div>
                      ) : (
                        <div className="flex items-start gap-2">
                          {selectedEngineIds.length > 0 ? (
                            <CheckCircle2 className="h-4 w-4 text-emerald-500 mt-0.5" />
                          ) : (
                            <AlertTriangle className="h-4 w-4 text-amber-500 mt-0.5" />
                          )}
                          <span>
                            {selectedEngineIds.length > 0
                              ? t("validation.enginesOk", { count: selectedEngineIds.length })
                              : t("validation.enginesMissing")}
                          </span>
                        </div>
                      )}
                      <div className="flex items-start gap-2">
                        {hasConfig ? (
                          <CheckCircle2 className="h-4 w-4 text-emerald-500 mt-0.5" />
                        ) : (
                          <AlertTriangle className="h-4 w-4 text-amber-500 mt-0.5" />
                        )}
                        <span>
                          {hasConfig
                            ? (isConfigEdited
                              ? t("validation.configEdited")
                              : selectMode === "preset"
                                ? t("validation.configFromPreset")
                                : t("validation.configFromEngine"))
                            : t("validation.configMissing")}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              )}
          </>
        </div>

        {/* Sticky footer */}
        <DialogFooter className="px-6 py-4 border-t shrink-0 bg-background">
          <div className="flex items-center justify-between w-full">
            <div className="text-sm text-muted-foreground">
              {currentStep === 1 && selectMode === "custom" && selectedEngineIds.length > 0 && (
                <span className="text-primary">{t("selectedCount", { count: selectedEngineIds.length })}</span>
              )}
              {currentStep === 1 && selectMode === "preset" && selectedPreset && (
                <span className="text-primary">{selectedPreset.name}</span>
              )}
              {currentStep === 2 && (
                <span className={canStart ? "text-primary" : undefined}>
                  {canStart ? t("validation.yamlOk") : t("validation.yamlError")}
                </span>
              )}
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isSubmitting}
              >
                {t("cancel")}
              </Button>
              {currentStep > 1 && (
                <Button
                  variant="outline"
                  onClick={() => setCurrentStep((prev) => Math.max(1, prev - 1))}
                  disabled={isSubmitting}
                >
                  {t("back")}
                </Button>
              )}
              {currentStep === 1 ? (
                <>
                  <Button
                    onClick={() => setCurrentStep(2)}
                    disabled={!canProceedToReview || isSubmitting}
                  >
                    {t("next")}
                  </Button>
                </>
              ) : (
                <Button
                  onClick={() => handleInitiate()}
                  disabled={!canStart || isSubmitting}
                >
                  {isSubmitting ? (
                    <>
                      <LoadingSpinner />
                      {t("initiating")}
                    </>
                  ) : (
                    <>
                      <Play className="h-4 w-4 mr-2" />
                      {t("startScan")}
                    </>
                  )}
                </Button>
              )}
            </div>
          </div>
        </DialogFooter>
      </DialogContent>
      
      {/* Overwrite confirmation dialog */}
      <AlertDialog
        open={showOverwriteConfirm}
        onOpenChange={(open) => {
          if (!open) {
            handleOverwriteCancel()
          } else {
            setShowOverwriteConfirm(true)
          }
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("overwriteConfirm.title")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("overwriteConfirm.description")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={handleOverwriteCancel}>
              {t("overwriteConfirm.cancel")}
            </AlertDialogCancel>
            <AlertDialogAction onClick={handleOverwriteConfirm}>
              {t("overwriteConfirm.confirm")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </Dialog>
  )
}
