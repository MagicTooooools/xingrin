"use client"

import React, { useState } from "react"
import { FileCode, Save, X, AlertCircle, CheckCircle2, ArrowLeft, ArrowRight, Lock, Check } from "@/components/icons"
import * as yaml from "js-yaml"
import { useTranslations } from "next-intl"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Badge } from "@/components/ui/badge"
import { CodeEditor } from "@/components/ui/code-editor"
import { toast } from "sonner"
import { usePresetEngines } from "@/hooks/use-engines"
import { cn } from "@/lib/utils"
import { parseEngineCapabilities } from "@/lib/engine-config"
import type { PresetEngine } from "@/types/engine.types"

interface EngineCreateDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSave?: (name: string, yamlContent: string) => Promise<void>
  /** Pre-selected preset (skip step 1) */
  preSelectedPreset?: PresetEngine
}

/**
 * Create new engine dialog - requires selecting a preset template
 */
export function EngineCreateDialog({
  open,
  onOpenChange,
  onSave,
  preSelectedPreset,
}: EngineCreateDialogProps) {
  const t = useTranslations("scan.engine.create")
  const tEngine = useTranslations("scan.engine")
  const tToast = useTranslations("toast")
  const tCommon = useTranslations("common.actions")
  
  // Step: 1 = select preset, 2 = edit config
  const [step, setStep] = useState<1 | 2>(1)
  const [selectedPreset, setSelectedPreset] = useState<PresetEngine | null>(null)
  const [engineName, setEngineName] = useState("")
  const [yamlContent, setYamlContent] = useState("")
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [yamlError, setYamlError] = useState<{ message: string; line?: number; column?: number } | null>(null)
  
  const { data: presetEngines = [] } = usePresetEngines()

  // Reset form when dialog opens
  React.useEffect(() => {
    if (open) {
      if (preSelectedPreset) {
        // Skip to step 2 if preset is pre-selected
        setSelectedPreset(preSelectedPreset)
        setEngineName(`${preSelectedPreset.name} (副本)`)
        setYamlContent(preSelectedPreset.configuration)
        setStep(2)
      } else {
        setStep(1)
        setSelectedPreset(null)
        setEngineName("")
        setYamlContent("")
      }
      setYamlError(null)
    }
  }, [open, preSelectedPreset])

  // Validate YAML syntax
  const validateYaml = (content: string) => {
    if (!content.trim()) {
      setYamlError(null)
      return true
    }

    try {
      yaml.load(content)
      setYamlError(null)
      return true
    } catch (error) {
      const yamlError = error as yaml.YAMLException
      setYamlError({
        message: yamlError.message,
        line: yamlError.mark?.line ? yamlError.mark.line + 1 : undefined,
        column: yamlError.mark?.column ? yamlError.mark.column + 1 : undefined,
      })
      return false
    }
  }

  // Handle editor content change
  const handleEditorChange = (value: string) => {
    setYamlContent(value)
    validateYaml(value)
  }

  // Handle save
  const handleSave = async () => {
    // Validate engine name
    if (!engineName.trim()) {
      toast.error(tToast("engineNameRequired"))
      return
    }

    // YAML validation
    if (!yamlContent.trim()) {
      toast.error(tToast("configRequired"))
      return
    }

    if (!validateYaml(yamlContent)) {
      toast.error(tToast("yamlSyntaxError"), {
        description: yamlError?.message,
      })
      return
    }

    setIsSubmitting(true)
    try {
      if (onSave) {
        await onSave(engineName, yamlContent)
      } else {
        // TODO: Call actual API to create engine
        await new Promise(resolve => setTimeout(resolve, 1000))
      }
      
      toast.success(tToast("engineCreateSuccess"), {
        description: tToast("engineCreateSuccessDesc", { name: engineName }),
      })
      onOpenChange(false)
    } catch (error) {
      toast.error(tToast("engineCreateFailed"), {
        description: error instanceof Error ? error.message : tToast("unknownError"),
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  // Handle preset selection and go to step 2
  const handleSelectPreset = (preset: PresetEngine) => {
    setSelectedPreset(preset)
  }

  const handleNextStep = () => {
    if (!selectedPreset) return
    setEngineName(`${selectedPreset.name} (副本)`)
    setYamlContent(selectedPreset.configuration)
    setStep(2)
  }

  const handleBackToStep1 = () => {
    // Only go back if not pre-selected
    if (!preSelectedPreset) {
      setStep(1)
    }
  }

  // Handle close
  const handleClose = () => {
    if (step === 2 && (engineName.trim() || yamlContent)) {
      const confirmed = window.confirm(t("confirmClose"))
      if (!confirmed) return
    }
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-4xl max-w-[calc(100%-2rem)] h-[85vh] flex flex-col p-0">
        <div className="flex flex-col h-full">
          <DialogHeader className="px-6 pt-6 pb-4 border-b">
            <DialogTitle className="flex items-center gap-2">
              <FileCode className="h-5 w-5" />
              {t("title")}
            </DialogTitle>
            <DialogDescription>
              {step === 1 ? t("selectPresetDesc") : t("editConfigDesc")}
            </DialogDescription>
            {/* Step indicator */}
            <div className="flex items-center gap-2 mt-3">
              <div className={cn(
                "flex items-center justify-center w-6 h-6 rounded-full text-xs font-medium",
                step === 1 ? "bg-primary text-primary-foreground" : "bg-muted text-muted-foreground"
              )}>
                1
              </div>
              <div className="w-8 h-px bg-border" />
              <div className={cn(
                "flex items-center justify-center w-6 h-6 rounded-full text-xs font-medium",
                step === 2 ? "bg-primary text-primary-foreground" : "bg-muted text-muted-foreground"
              )}>
                2
              </div>
            </div>
          </DialogHeader>

          {step === 1 ? (
            /* Step 1: Select preset template */
            <>
              <ScrollArea className="flex-1 px-6 py-4">
                <div className="grid grid-cols-2 gap-3">
                  {presetEngines.map((preset) => (
                    <button
                      key={preset.id}
                      onClick={() => handleSelectPreset(preset)}
                      className={cn(
                        "text-left p-4 rounded-lg border-2 transition-all",
                        selectedPreset?.id === preset.id
                          ? "border-primary bg-primary/5"
                          : "border-border hover:border-primary/50 hover:bg-muted/50"
                      )}
                    >
                      <div className="flex items-start gap-3">
                        <div className={cn(
                          "flex h-8 w-8 items-center justify-center rounded-lg shrink-0",
                          selectedPreset?.id === preset.id ? "bg-primary/10" : "bg-muted"
                        )}>
                          {selectedPreset?.id === preset.id ? (
                            <Check className="h-4 w-4 text-primary" />
                          ) : (
                            <Lock className="h-4 w-4 text-muted-foreground" />
                          )}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium text-sm">{preset.name}</div>
                          <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                            {preset.description}
                          </p>
                          <div className="flex flex-wrap gap-1 mt-2">
                            {(() => {
                              const features = parseEngineCapabilities(preset.configuration)
                              return (
                                <>
                                  {features.slice(0, 3).map((feature) => (
                                    <Badge key={feature} variant="secondary" className="text-[10px] px-1.5 py-0">
                                      {tEngine(`features.${feature}`)}
                                    </Badge>
                                  ))}
                                  {features.length > 3 && (
                                    <Badge variant="secondary" className="text-[10px] px-1.5 py-0">
                                      +{features.length - 3}
                                    </Badge>
                                  )}
                                </>
                              )
                            })()}
                          </div>
                        </div>
                      </div>
                    </button>
                  ))}
                </div>
              </ScrollArea>

              <DialogFooter className="px-6 py-4 border-t gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleClose}
                >
                  <X className="h-4 w-4" />
                  {tCommon("cancel")}
                </Button>
                <Button
                  type="button"
                  onClick={handleNextStep}
                  disabled={!selectedPreset}
                >
                  {t("nextStep")}
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </DialogFooter>
            </>
          ) : (
            /* Step 2: Edit configuration */
            <>
              <div className="flex-1 overflow-hidden px-6 py-4">
                <div className="flex flex-col h-full gap-4">
                  {/* Engine name input */}
                  <div className="space-y-2">
                    <Label htmlFor="engine-name">
                      {t("engineName")} <span className="text-destructive">*</span>
                    </Label>
                    <Input
                      id="engine-name"
                      value={engineName}
                      onChange={(e) => setEngineName(e.target.value)}
                      placeholder={t("engineNamePlaceholder")}
                      disabled={isSubmitting}
                      className="max-w-md"
                    />
                  </div>

                  {/* YAML editor */}
                  <div className="flex flex-col flex-1 min-h-0 gap-2">
                    <div className="flex items-center justify-between">
                      <Label>{t("yamlConfig")}</Label>
                      {/* Syntax validation status */}
                      <div className="flex items-center gap-2">
                        {yamlContent.trim() && (
                          yamlError ? (
                            <div className="flex items-center gap-1 text-xs text-destructive">
                              <AlertCircle className="h-3.5 w-3.5" />
                              <span>{t("syntaxError")}</span>
                            </div>
                          ) : (
                            <div className="flex items-center gap-1 text-xs text-green-600 dark:text-green-400">
                              <CheckCircle2 className="h-3.5 w-3.5" />
                              <span>{t("syntaxValid")}</span>
                            </div>
                          )
                        )}
                      </div>
                    </div>

                    {/* CodeMirror Editor */}
                    <CodeEditor
                      value={yamlContent}
                      onChange={handleEditorChange}
                      language="yaml"
                      readOnly={isSubmitting}
                      className={`flex-1 ${yamlError ? 'border-destructive' : ''}`}
                      showLineNumbers
                      showFoldGutter
                    />

                    {/* Error message display */}
                    {yamlError && (
                      <div className="flex items-start gap-2 p-3 bg-destructive/10 border border-destructive/20 rounded-md">
                        <AlertCircle className="h-4 w-4 text-destructive mt-0.5 flex-shrink-0" />
                        <div className="flex-1 text-xs">
                          <p className="font-semibold text-destructive mb-1">
                            {yamlError.line && yamlError.column
                              ? t("errorLocation", { line: yamlError.line, column: yamlError.column })
                              : tToast("yamlSyntaxError")}
                          </p>
                          <p className="text-muted-foreground">{yamlError.message}</p>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              <DialogFooter className="px-6 py-4 border-t">
                <div className="flex items-center justify-between w-full">
                  {!preSelectedPreset ? (
                    <Button
                      type="button"
                      variant="ghost"
                      onClick={handleBackToStep1}
                      disabled={isSubmitting}
                    >
                      <ArrowLeft className="h-4 w-4" />
                      {t("prevStep")}
                    </Button>
                  ) : (
                    <div />
                  )}
                  <div className="flex items-center gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleClose}
                      disabled={isSubmitting}
                    >
                      <X className="h-4 w-4" />
                      {tCommon("cancel")}
                    </Button>
                    <Button
                      type="button"
                      onClick={handleSave}
                      disabled={isSubmitting || !engineName.trim() || !!yamlError}
                    >
                      {isSubmitting ? (
                        <>
                          <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                          {t("creating")}
                        </>
                      ) : (
                        <>
                          <Save className="h-4 w-4" />
                          {t("createEngine")}
                        </>
                      )}
                    </Button>
                  </div>
                </div>
              </DialogFooter>
            </>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
