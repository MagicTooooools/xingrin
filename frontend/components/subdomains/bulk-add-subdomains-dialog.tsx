"use client"

import React, { useState, useRef } from "react"
import { Plus, Globe } from "lucide-react"
import { useTranslations } from "next-intl"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { LoadingSpinner } from "@/components/loading-spinner"
import { SubdomainValidator } from "@/lib/subdomain-validator"
import { useBulkCreateSubdomains } from "@/hooks/use-subdomains"

interface BulkAddSubdomainsDialogProps {
  targetId: number
  targetName?: string
  open?: boolean
  onOpenChange?: (open: boolean) => void
  onSuccess?: () => void
}

/**
 * Bulk add subdomains dialog component
 * 
 * Following the design pattern of AddTargetDialog, provides a text input with line numbers,
 * supporting real-time validation and error prompts.
 */
export function BulkAddSubdomainsDialog({
  targetId,
  targetName,
  open: externalOpen,
  onOpenChange: externalOnOpenChange,
  onSuccess,
}: BulkAddSubdomainsDialogProps) {
  const t = useTranslations("bulkAdd.subdomain")
  const tCommon = useTranslations("common.actions")
  
  // Dialog open/close state
  const [internalOpen, setInternalOpen] = useState(false)
  const open = externalOpen !== undefined ? externalOpen : internalOpen
  const setOpen = externalOnOpenChange || setInternalOpen

  // Form data state
  const [inputText, setInputText] = useState("")

  // Validation result state
  const [validationResult, setValidationResult] = useState<{
    validCount: number
    invalidCount: number
    duplicateCount: number
    firstError?: { index: number; subdomain: string; error: string }
  } | null>(null)

  // Refs for line numbers column and textarea (for synchronized scrolling)
  const lineNumbersRef = useRef<HTMLDivElement | null>(null)
  const textareaRef = useRef<HTMLTextAreaElement | null>(null)

  // Use bulk create mutation
  const bulkCreateSubdomains = useBulkCreateSubdomains()

  // Handle input change
  const handleInputChange = (value: string) => {
    setInputText(value)

    // Parse and validate
    const parsed = SubdomainValidator.parse(value)
    if (parsed.length === 0) {
      setValidationResult(null)
      return
    }

    const result = SubdomainValidator.validateBatch(parsed)
    setValidationResult({
      validCount: result.validCount,
      invalidCount: result.invalidCount,
      duplicateCount: result.duplicateCount,
      firstError: result.invalidItems[0]
        ? {
            index: result.invalidItems[0].index,
            subdomain: result.invalidItems[0].subdomain,
            error: result.invalidItems[0].error || t("formatInvalid"),
          }
        : undefined,
    })
  }

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!inputText.trim()) return
    if (!validationResult || validationResult.validCount === 0) return

    // Parse valid subdomains
    const parsed = SubdomainValidator.parse(inputText)
    const result = SubdomainValidator.validateBatch(parsed)

    bulkCreateSubdomains.mutate(
      { targetId, subdomains: result.subdomains },
      {
        onSuccess: () => {
          // Reset form
          setInputText("")
          setValidationResult(null)
          // Close dialog
          setOpen(false)
          // Call external callback
          onSuccess?.()
        },
      }
    )
  }

  // Handle dialog close
  const handleOpenChange = (newOpen: boolean) => {
    if (!bulkCreateSubdomains.isPending) {
      setOpen(newOpen)
      if (!newOpen) {
        setInputText("")
        setValidationResult(null)
      }
    }
  }

  // Synchronized scrolling
  const handleTextareaScroll = (e: React.UIEvent<HTMLTextAreaElement>) => {
    if (lineNumbersRef.current) {
      lineNumbersRef.current.scrollTop = e.currentTarget.scrollTop
    }
  }

  // Calculate line count
  const lineCount = Math.max(inputText.split("\n").length, 8)

  // Form validation
  const isFormValid =
    inputText.trim().length > 0 &&
    validationResult !== null &&
    validationResult.validCount > 0

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {externalOpen === undefined && (
        <DialogTrigger asChild>
          <Button size="sm" variant="outline">
            <Plus className="h-4 w-4" />
            {t("bulkAdd")}
          </Button>
        </DialogTrigger>
      )}

      <DialogContent className="sm:max-w-[650px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Globe className="h-5 w-5" />
            <span>{t("title")}</span>
          </DialogTitle>
          <DialogDescription>
            {t("description")}
            {targetName && (
              <span className="block mt-1">
                {t("belongsTo")} <code className="bg-muted px-1 rounded">{targetName}</code>
              </span>
            )}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="subdomains">
                {t("label")} <span className="text-destructive">*</span>
              </Label>
              <div className="flex border rounded-md overflow-hidden h-[220px]">
                {/* Line numbers column */}
                <div className="flex-shrink-0 w-12 border-r bg-muted/50">
                  <div
                    ref={lineNumbersRef}
                    className="py-3 px-2 text-right font-mono text-xs text-muted-foreground leading-[1.4] h-full overflow-y-auto scrollbar-hide"
                  >
                    {Array.from({ length: lineCount }, (_, i) => (
                      <div key={i + 1} className="h-[20px]">
                        {i + 1}
                      </div>
                    ))}
                  </div>
                </div>
                {/* Input textarea */}
                <div className="flex-1 overflow-hidden">
                  <Textarea
                    ref={textareaRef}
                    id="subdomains"
                    value={inputText}
                    onChange={(e) => handleInputChange(e.target.value)}
                    onScroll={handleTextareaScroll}
                    placeholder={t("placeholder")}
                    disabled={bulkCreateSubdomains.isPending}
                    className="font-mono h-full overflow-y-auto resize-none border-0 focus-visible:ring-0 focus-visible:ring-offset-0 leading-[1.4] text-sm py-3"
                    style={{ lineHeight: "20px" }}
                  />
                </div>
              </div>

              {/* Validation summary */}
              {validationResult && (
                <div className="text-xs space-y-1">
                  <div className="text-muted-foreground">
                    {t("valid", { count: validationResult.validCount })}
                    {validationResult.duplicateCount > 0 && (
                      <span className="text-yellow-600 ml-2">
                        {t("duplicate", { count: validationResult.duplicateCount })}
                      </span>
                    )}
                    {validationResult.invalidCount > 0 && (
                      <span className="text-destructive ml-2">
                        {t("invalid", { count: validationResult.invalidCount })}
                      </span>
                    )}
                  </div>
                  {validationResult.firstError && (
                    <div className="text-destructive">
                      {t("lineError", {
                        line: validationResult.firstError.index + 1,
                        value: validationResult.firstError.subdomain,
                        error: validationResult.firstError.error,
                      })}
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => handleOpenChange(false)}
              disabled={bulkCreateSubdomains.isPending}
            >
              {tCommon("cancel")}
            </Button>
            <Button
              type="submit"
              disabled={bulkCreateSubdomains.isPending || !isFormValid}
            >
              {bulkCreateSubdomains.isPending ? (
                <>
                  <LoadingSpinner />
                  {t("creating")}
                </>
              ) : (
                <>
                  <Plus className="h-4 w-4" />
                  {t("bulkAdd")}
                </>
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
