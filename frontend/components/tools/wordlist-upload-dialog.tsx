"use client"

import { useState, type FormEvent } from "react"
import { Upload, X, FileText } from "@/components/icons"
import { useTranslations } from "next-intl"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useUploadWordlist } from "@/hooks/use-wordlists"
import { cn } from "@/lib/utils"

interface WordlistUploadDialogProps {
  trigger?: React.ReactNode
}

export function WordlistUploadDialog({ trigger }: WordlistUploadDialogProps) {
  const t = useTranslations("tools.wordlists.uploadDialog")
  const tWordlists = useTranslations("tools.wordlists")
  
  const [open, setOpen] = useState(false)
  const [name, setName] = useState("")
  const [description, setDescription] = useState("")
  const [file, setFile] = useState<File | null>(null)
  const [isDragActive, setIsDragActive] = useState(false)

  const uploadMutation = useUploadWordlist()

  const resetForm = () => {
    setName("")
    setDescription("")
    setFile(null)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragActive(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragActive(false)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragActive(false)
    const droppedFile = e.dataTransfer.files[0]
    if (droppedFile && droppedFile.name.endsWith(".txt")) {
      setFile(droppedFile)
      if (!name) {
        setName(droppedFile.name.replace(/\.[^/.]+$/, ""))
      }
    }
  }

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0]
    if (selectedFile) {
      setFile(selectedFile)
      if (!name) {
        setName(selectedFile.name.replace(/\.[^/.]+$/, ""))
      }
    }
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    if (!name || !file) return

    uploadMutation.mutate(
      { name, description: description || undefined, file },
      {
        onSuccess: () => {
          resetForm()
          setOpen(false)
        },
      }
    )
  }

  const removeFile = () => {
    setFile(null)
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  return (
    <Dialog open={open} onOpenChange={(v) => { setOpen(v); if (!v) resetForm() }}>
      <DialogTrigger asChild>
        {trigger || (
          <Button>
            <Upload className="mr-2 h-4 w-4" />
            {tWordlists("upload")}
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("desc")}</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            className={cn(
              "relative flex flex-col items-center justify-center rounded-lg border-2 border-dashed p-6 transition-colors",
              isDragActive
                ? "border-primary bg-primary/5"
                : "border-muted-foreground/25 hover:border-muted-foreground/50",
              file && "border-solid border-muted-foreground/25"
            )}
          >
            {file ? (
              <div className="flex w-full items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                  <FileText className="h-5 w-5 text-primary" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="truncate font-medium text-sm">{file.name}</p>
                  <p className="text-xs text-muted-foreground">
                    {formatFileSize(file.size)}
                  </p>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="h-8 w-8 shrink-0"
                  onClick={removeFile}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            ) : (
              <>
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted">
                  <Upload className="h-6 w-6 text-muted-foreground" />
                </div>
                <div className="mt-3 text-center">
                  <p className="text-sm font-medium">{t("dragHint")}</p>
                  <p className="mt-1 text-xs text-muted-foreground">
                    {" "}
                    <label className="cursor-pointer text-primary hover:underline">
                      {t("selectFile")}
                      <input
                        type="file"
                        accept=".txt"
                        className="hidden"
                        onChange={handleFileSelect}
                      />
                    </label>
                  </p>
                  <p className="mt-2 text-xs text-muted-foreground">
                    {t("fileHint")}
                  </p>
                </div>
              </>
            )}
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="name">
                {tWordlists("name")} <span className="text-destructive">*</span>
              </Label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder={t("namePlaceholder")}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">{t("descLabel")}</Label>
              <Input
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder={t("descPlaceholder")}
              />
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
              disabled={uploadMutation.isPending}
            >
              {t("cancel")}
            </Button>
            <Button
              type="submit"
              disabled={uploadMutation.isPending || !file || !name}
            >
              {uploadMutation.isPending ? t("uploading") : t("uploadButton")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
