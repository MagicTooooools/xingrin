"use client"

import React, { useState, useEffect } from "react"
import { FileCode, Save, X, AlertCircle, CheckCircle2, AlertTriangle } from "@/components/icons"
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
import { CodeEditor } from "@/components/ui/code-editor"
import { toast } from "sonner"
import type { ScanEngine } from "@/types/engine.types"

interface EngineEditDialogProps {
  engine: ScanEngine | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onSave?: (engineId: number, yamlContent: string) => Promise<void>
}

/**
 * Engine configuration edit dialog
 * Uses Monaco Editor to provide VSCode-level editing experience
 */
export function EngineEditDialog({
  engine,
  open,
  onOpenChange,
  onSave,
}: EngineEditDialogProps) {
  const t = useTranslations("scan.engine.edit")
  const tToast = useTranslations("toast")
  const tCommon = useTranslations("common.actions")
  const [yamlContent, setYamlContent] = useState("")
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [hasChanges, setHasChanges] = useState(false)
  const [yamlError, setYamlError] = useState<{ message: string; line?: number; column?: number } | null>(null)

  // Generate sample YAML configuration
  const generateSampleYaml = (engine: ScanEngine) => {
    return `# Engine name: ${engine.name}

# ==================== Subdomain Discovery ====================
subdomain_discovery:
  tools:
    subfinder:
      enabled: true
      timeout: 600      # 10 minutes (required)
      
    assetfinder:
      enabled: true
      timeout: 600      # 10 minutes (required)


# ==================== Port Scan ====================
port_scan:
  tools:
    naabu_active:
      enabled: true
      timeout: auto     # Auto calculate
      threads: 5
      top-ports: 100
      rate: 10
      
    naabu_passive:
      enabled: true
      timeout: auto


# ==================== Site Scan ====================
site_scan:
  tools:
    httpx:
      enabled: true
      timeout: auto         # Auto calculate
      # screenshot: true    # Enable site screenshot (requires Chromium)


# ==================== Directory Scan ====================
directory_scan:
  tools:
    ffuf:
      enabled: true
      timeout: auto                            # Auto calculate timeout
      wordlist: ~/Desktop/dirsearch_dicc.txt   # Wordlist file path (required)
      delay: 0.1-2.0
      threads: 10
      request_timeout: 10
      match_codes: 200,201,301,302,401,403


# ==================== URL Fetch ====================
url_fetch:
  tools:
    waymore:
      enabled: true
      timeout: auto
    
    katana:
      enabled: true
      timeout: auto
      depth: 5
      threads: 10
      rate-limit: 30
      random-delay: 1
      retry: 2
      request-timeout: 12
    
    uro:
      enabled: true
      timeout: auto
    
    httpx:
      enabled: true
      timeout: auto
`
  }

  // When engine changes, update YAML content
  useEffect(() => {
    if (engine && open) {
      // TODO: Get actual YAML configuration from backend API
      // If engine has configuration use it, otherwise use sample configuration
      const content = engine.configuration || generateSampleYaml(engine)
      setYamlContent(content)
      setHasChanges(false)
      setYamlError(null)
    }
  }, [engine, open])

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
    setHasChanges(true)
    validateYaml(value)
  }

  // Handle save
  const handleSave = async () => {
    if (!engine) return

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
        await onSave(engine.id, yamlContent)
      } else {
        // TODO: Call actual API to save YAML configuration
        await new Promise(resolve => setTimeout(resolve, 1000))
      }

      setHasChanges(false)
      onOpenChange(false)
    } catch {
      // Error toast is handled by useUpdateEngine hook
    } finally {
      setIsSubmitting(false)
    }
  }

  // Handle close
  const handleClose = () => {
    if (hasChanges) {
      const confirmed = window.confirm(t("confirmClose"))
      if (!confirmed) return
    }
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-6xl max-w-[calc(100%-2rem)] h-[90vh] flex flex-col p-0">
        <div className="flex flex-col h-full">
          <DialogHeader className="px-6 pt-6 pb-4 border-b">
            <DialogTitle className="flex items-center gap-2">
              <FileCode className="h-5 w-5" />
              {t("title", { name: engine?.name ?? "" })}
            </DialogTitle>
            <DialogDescription>
              {t("desc")}
            </DialogDescription>
          </DialogHeader>

          <div className="flex-1 overflow-hidden px-6 py-4">
            <div className="flex flex-col h-full gap-2">
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
                className={yamlError ? 'border-destructive' : ''}
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
              <p className="flex items-center gap-1 text-xs text-amber-600 dark:text-amber-400">
                <AlertTriangle className="h-3.5 w-3.5" />
                {t("unsavedChanges")}
              </p>
            </div>
          </div>

          <DialogFooter className="px-6 py-4 border-t gap-2">
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
              disabled={isSubmitting || !hasChanges || !!yamlError}
            >
              {isSubmitting ? (
                <>
                  <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                  {t("saving")}
                </>
              ) : (
                <>
                  <Save className="h-4 w-4" />
                  {t("saveConfig")}
                </>
              )}
            </Button>
          </DialogFooter>
        </div>
      </DialogContent>
    </Dialog>
  )
}
