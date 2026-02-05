"use client"

import React, { useState, useEffect } from "react"
import { FileText, Save, X, AlertTriangle } from "@/components/icons"
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
import { useWordlistContent, useUpdateWordlistContent } from "@/hooks/use-wordlists"
import type { Wordlist } from "@/types/wordlist.types"

interface WordlistEditDialogProps {
  wordlist: Wordlist | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function WordlistEditDialog({
  wordlist,
  open,
  onOpenChange,
}: WordlistEditDialogProps) {
  const t = useTranslations("tools.wordlists.editDialog")
  
  const [content, setContent] = useState("")
  const [hasChanges, setHasChanges] = useState(false)

  const { data: originalContent, isLoading } = useWordlistContent(
    open && wordlist ? wordlist.id : null
  )
  const updateMutation = useUpdateWordlistContent()

  useEffect(() => {
    if (originalContent !== undefined && open) {
      setContent(originalContent)
      setHasChanges(false)
    }
  }, [originalContent, open])

  useEffect(() => {
    if (!open) {
      setContent("")
      setHasChanges(false)
    }
  }, [open])

  const handleEditorChange = (value: string) => {
    setContent(value)
    setHasChanges(value !== originalContent)
  }

  const handleSave = async () => {
    if (!wordlist) return

    updateMutation.mutate(
      { id: wordlist.id, content },
      {
        onSuccess: () => {
          setHasChanges(false)
          onOpenChange(false)
        },
      }
    )
  }

  const handleClose = () => {
    if (hasChanges) {
      const confirmed = window.confirm(t("confirmClose"))
      if (!confirmed) return
    }
    onOpenChange(false)
  }

  const lineCount = content.split("\n").length

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-6xl max-w-[calc(100%-2rem)] h-[90vh] flex flex-col p-0">
        <div className="flex flex-col h-full">
          <DialogHeader className="px-6 pt-6 pb-4 border-b">
            <DialogTitle className="flex items-center gap-2">
              <FileText className="h-5 w-5" />
              {t("title", { name: wordlist?.name || "" })}
            </DialogTitle>
            <DialogDescription>{t("desc")}</DialogDescription>
          </DialogHeader>

          <div className="flex-1 overflow-hidden px-6 py-4">
            <div className="flex flex-col h-full gap-2">
              <div className="flex items-center justify-between">
                <Label>{t("content")}</Label>
                <div className="flex items-center gap-4 text-xs text-muted-foreground">
                  <span>{t("lines", { count: lineCount.toLocaleString() })}</span>
                  {wordlist?.fileHash && (
                    <span title={wordlist.fileHash}>
                      {t("hash")}: {wordlist.fileHash.slice(0, 12)}...
                    </span>
                  )}
                </div>
              </div>

              {isLoading ? (
                <div className="flex items-center justify-center h-full border rounded-md">
                  <div className="flex flex-col items-center gap-2">
                    <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
                    <p className="text-sm text-muted-foreground">{t("loading")}</p>
                  </div>
                </div>
              ) : (
                <CodeEditor
                  value={content}
                  onChange={handleEditorChange}
                  language="plaintext"
                  readOnly={updateMutation.isPending}
                  showLineNumbers
                  showFoldGutter={false}
                />
              )}

              {hasChanges && (
                <p className="flex items-center gap-1 text-xs text-amber-600 dark:text-amber-400">
                  <AlertTriangle className="h-3.5 w-3.5" />
                  {t("unsavedChanges")}
                </p>
              )}
            </div>
          </div>

          <DialogFooter className="px-6 py-4 border-t gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={updateMutation.isPending}
            >
              <X className="h-4 w-4" />
              {t("cancel")}
            </Button>
            <Button
              type="button"
              onClick={handleSave}
              disabled={updateMutation.isPending || !hasChanges}
            >
              {updateMutation.isPending ? (
                <>
                  <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                  {t("saving")}
                </>
              ) : (
                <>
                  <Save className="h-4 w-4" />
                  {t("save")}
                </>
              )}
            </Button>
          </DialogFooter>
        </div>
      </DialogContent>
    </Dialog>
  )
}
