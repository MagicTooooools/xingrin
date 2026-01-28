"use client"

import { useEffect, useMemo, useRef, useState } from "react"
import { useTranslations } from "next-intl"
import {
  IconServer,
  IconCloud,
  IconCloudOff,
  IconHeartbeat,
  IconLoader2,
} from "@tabler/icons-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { ConfirmDialog } from "@/components/ui/confirm-dialog"
import { Skeleton } from "@/components/ui/skeleton"
import { CopyablePopoverContent } from "@/components/ui/copyable-popover-content"
import { useFormatRelativeTime } from "@/lib/i18n-format"
import {
  useAgents,
  useCreateRegistrationToken,
  useDeleteAgent,
} from "@/hooks/use-agents"
import type { Agent, RegistrationTokenResponse } from "@/types/agent.types"
import { AgentConfigDialog } from "./worker-dialog"
import { AgentCardCompact } from "./agent-card-compact"

const FALLBACK_SERVER_URL = "https://your-orbit-server"

function InstallDialog({
  token,
  onGenerate,
  isGenerating,
  open,
  onOpenChange,
}: {
  token: RegistrationTokenResponse | null
  onGenerate: () => void
  isGenerating: boolean
  open: boolean
  onOpenChange: (open: boolean) => void
}) {
  const t = useTranslations("settings.workers")
  const formatRelativeTime = useFormatRelativeTime()
  const [serverUrl, setServerUrl] = useState("")
  const autoRequestedRef = useRef(false)

  useEffect(() => {
    if (typeof window !== "undefined") {
      setServerUrl(window.location.origin)
    }
  }, [])

  const safeServerUrl = serverUrl || FALLBACK_SERVER_URL
  const registrationToken = token?.token || ""
  const hasToken = Boolean(registrationToken)

  const tokenExpiresAt = useMemo(() => {
    if (!token?.expiresAt) return 0
    const timestamp = new Date(token.expiresAt).getTime()
    return Number.isNaN(timestamp) ? 0 : timestamp
  }, [token])

  const isTokenValid = tokenExpiresAt > Date.now()

  const linuxCommand = useMemo(() => {
    if (!registrationToken) return ""
    return `curl -fsSL "${safeServerUrl}/api/agents/install.sh?token=${registrationToken}" | bash`
  }, [safeServerUrl, registrationToken])

  const windowsCommand = useMemo(() => {
    if (!registrationToken) return ""
    return `irm "${safeServerUrl}/api/agents/install.ps1?token=${registrationToken}" | iex`
  }, [safeServerUrl, registrationToken])

  useEffect(() => {
    if (!open) {
      autoRequestedRef.current = false
      return
    }
    if (isGenerating || isTokenValid) return
    if (autoRequestedRef.current) return
    autoRequestedRef.current = true
    onGenerate()
  }, [open, isGenerating, isTokenValid, onGenerate])

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <Card className="border-muted/60 bg-gradient-to-br from-muted/40 via-background to-background">
        <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <IconHeartbeat className="h-5 w-5 text-primary" />
              {t("install.title")}
            </CardTitle>
            <CardDescription>{t("install.desc")}</CardDescription>
          </div>
          <DialogTrigger asChild>
            <Button size="sm">{t("install.openDialog")}</Button>
          </DialogTrigger>
        </CardHeader>
      </Card>

      <DialogContent className="sm:max-w-[720px]">
        <DialogHeader>
          <DialogTitle>{t("install.title")}</DialogTitle>
          <DialogDescription>{t("install.desc")}</DialogDescription>
        </DialogHeader>
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <Badge variant="outline" className={hasToken ? "border-emerald-500/30 text-emerald-600" : "border-muted-foreground/30"}>
              {hasToken ? t("install.commandStatusReady") : t("install.commandStatusGenerating")}
            </Badge>
            {hasToken && token?.expiresAt && (
              <span>{t("install.commandExpires", { time: formatRelativeTime(token.expiresAt) })}</span>
            )}
          </div>
          <Button size="sm" onClick={onGenerate} disabled={isGenerating}>
            {isGenerating && <IconLoader2 className="mr-2 h-4 w-4 animate-spin" />}
            {token ? t("install.regenerateToken") : t("install.generateToken")}
          </Button>
        </div>

        <div className="grid gap-4 lg:grid-cols-[1.3fr,0.7fr]">
          <div className="space-y-3">
            <div className="rounded-lg border bg-background p-3 space-y-3">
              <div>
                <p className="text-sm font-medium">{t("install.commandTitle")}</p>
                <p className="text-xs text-muted-foreground">{t("install.commandDesc")}</p>
              </div>
              <Tabs defaultValue="linux" className="w-full">
                <TabsList className="grid w-full grid-cols-2">
                  <TabsTrigger value="linux">{t("install.linuxTab")}</TabsTrigger>
                  <TabsTrigger value="windows">{t("install.windowsTab")}</TabsTrigger>
                </TabsList>
                <TabsContent value="linux" className="space-y-2">
                  <div className="rounded-lg border bg-muted/40 p-3">
                    {hasToken ? (
                      <CopyablePopoverContent value={linuxCommand} className="font-mono text-xs whitespace-pre-wrap" />
                    ) : (
                      <div className="text-xs text-muted-foreground">{t("install.commandPlaceholder")}</div>
                    )}
                  </div>
                </TabsContent>
                <TabsContent value="windows" className="space-y-2">
                  <div className="rounded-lg border bg-muted/40 p-3">
                    {hasToken ? (
                      <CopyablePopoverContent value={windowsCommand} className="font-mono text-xs whitespace-pre-wrap" />
                    ) : (
                      <div className="text-xs text-muted-foreground">{t("install.commandPlaceholder")}</div>
                    )}
                  </div>
                </TabsContent>
              </Tabs>
            </div>
          </div>

          <div className="space-y-3">
            <div className="rounded-lg border bg-muted/30 p-3">
              <p className="text-sm font-medium mb-2">{t("install.stepsTitle")}</p>
              <div className="grid gap-2 text-xs text-muted-foreground">
                <div className="flex items-center gap-2">
                  <span className="rounded-full border bg-background px-2 py-0.5 text-[10px] font-medium text-foreground">
                    {t("install.step1Label")}
                  </span>
                  <span>{t("install.step1Desc")}</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="rounded-full border bg-background px-2 py-0.5 text-[10px] font-medium text-foreground">
                    {t("install.step2Label")}
                  </span>
                  <span>{t("install.step2Desc")}</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="rounded-full border bg-background px-2 py-0.5 text-[10px] font-medium text-foreground">
                    {t("install.step3Label")}
                  </span>
                  <span>{t("install.step3Desc")}</span>
                </div>
              </div>
            </div>

            <div className="rounded-lg border bg-muted/30 p-3">
              <p className="text-sm font-medium mb-2">{t("install.requirementsTitle")}</p>
              <div className="grid gap-2 text-xs text-muted-foreground">
                <div className="flex items-center gap-2">
                  <span className="h-1.5 w-1.5 rounded-full bg-muted-foreground/60" />
                  <span>{t("install.requirementsDocker")}</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="h-1.5 w-1.5 rounded-full bg-muted-foreground/60" />
                  <span>{t("install.requirementsAccess")}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}

function EmptyState({ onOpenInstall }: { onOpenInstall: () => void }) {
  const t = useTranslations("settings.workers")

  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="p-4 rounded-full bg-muted mb-4">
        <IconServer className="h-12 w-12 text-muted-foreground" />
      </div>
      <h3 className="text-lg font-semibold mb-2">{t("empty.title")}</h3>
      <p className="text-sm text-muted-foreground mb-6 max-w-md">{t("empty.desc")}</p>
      <Button onClick={onOpenInstall}>{t("empty.cta")}</Button>
    </div>
  )
}

export function AgentList() {
  const t = useTranslations("settings.workers")
  const [page] = useState(1)
  const [pageSize] = useState(10)
  const [installOpen, setInstallOpen] = useState(false)
  const [selectedAgent, setSelectedAgent] = useState<Agent | null>(null)
  const [configDialogOpen, setConfigDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [agentToDelete, setAgentToDelete] = useState<Agent | null>(null)
  const [token, setToken] = useState<RegistrationTokenResponse | null>(null)

  const { data, isLoading } = useAgents(page, pageSize)
  const createToken = useCreateRegistrationToken()
  const deleteAgent = useDeleteAgent()

  const agents = data?.results || []
  const hasAgents = agents.length > 0

  const stats = useMemo(() => {
    const total = agents.length
    const online = agents.filter((agent) => agent.status === "online").length
    const offline = agents.filter((agent) => agent.status === "offline").length
    const unhealthy = agents.filter((agent) => {
      const state = agent.health?.state?.toLowerCase()
      return state && state !== "ok"
    }).length

    return [
      { label: t("stats.total"), value: total, icon: IconServer, color: "text-foreground" },
      { label: t("stats.online"), value: online, icon: IconCloud, color: "text-emerald-600" },
      { label: t("stats.offline"), value: offline, icon: IconCloudOff, color: "text-red-500" },
      { label: t("stats.unhealthy"), value: unhealthy, icon: IconHeartbeat, color: "text-amber-500" },
    ]
  }, [agents, t])

  const handleGenerateToken = async () => {
    try {
      const response = await createToken.mutateAsync()
      setToken(response)
    } catch {
      // handled by hook
    }
  }

  const handleConfigure = (agent: Agent) => {
    setSelectedAgent(agent)
    setConfigDialogOpen(true)
  }


  const handleDelete = (agent: Agent) => {
    setAgentToDelete(agent)
    setDeleteDialogOpen(true)
  }

  const confirmDelete = async () => {
    if (!agentToDelete) return
    try {
      await deleteAgent.mutateAsync(agentToDelete.id)
      setDeleteDialogOpen(false)
      setAgentToDelete(null)
    } catch {
      // handled by hook
    }
  }


  return (
    <div className="space-y-6">
      {hasAgents && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {stats.map((stat) => (
            <Card key={stat.label} className="p-4">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg bg-muted ${stat.color}`}>
                  <stat.icon className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-2xl font-bold">{stat.value}</p>
                  <p className="text-xs text-muted-foreground">{stat.label}</p>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      <Card>
        <CardHeader>
          <div className="flex items-start justify-between gap-4">
            <div>
              <CardTitle className="flex items-center gap-2">
                <IconServer className="h-5 w-5" />
                {t("agents.title")}
              </CardTitle>
              <CardDescription>{t("agents.desc")}</CardDescription>
            </div>
            <Dialog open={installOpen} onOpenChange={setInstallOpen}>
              <DialogTrigger asChild>
                <Button size="sm">
                  <IconHeartbeat className="h-4 w-4 mr-2" />
                  {t("install.openDialog")}
                </Button>
              </DialogTrigger>
              <DialogContent className="sm:max-w-[720px]">
                <DialogHeader>
                  <DialogTitle>{t("install.title")}</DialogTitle>
                  <DialogDescription>{t("install.desc")}</DialogDescription>
                </DialogHeader>
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div className="flex items-center gap-2 text-xs text-muted-foreground">
                    <Badge variant="outline" className={token ? "border-emerald-500/30 text-emerald-600" : "border-muted-foreground/30"}>
                      {token ? t("install.commandStatusReady") : t("install.commandStatusGenerating")}
                    </Badge>
                    {token?.expiresAt && (
                      <span>{t("install.commandExpires", { time: formatRelativeTime(token.expiresAt) })}</span>
                    )}
                  </div>
                  <Button size="sm" onClick={handleGenerateToken} disabled={createToken.isPending}>
                    {createToken.isPending && <IconLoader2 className="mr-2 h-4 w-4 animate-spin" />}
                    {token ? t("install.regenerateToken") : t("install.generateToken")}
                  </Button>
                </div>

                <div className="grid gap-4 lg:grid-cols-[1.3fr,0.7fr]">
                  <div className="space-y-3">
                    <div className="rounded-lg border bg-background p-3 space-y-3">
                      <div>
                        <p className="text-sm font-medium">{t("install.commandTitle")}</p>
                        <p className="text-xs text-muted-foreground">{t("install.commandDesc")}</p>
                      </div>
                      <Tabs defaultValue="linux" className="w-full">
                        <TabsList className="grid w-full grid-cols-2">
                          <TabsTrigger value="linux">{t("install.linuxTab")}</TabsTrigger>
                          <TabsTrigger value="windows">{t("install.windowsTab")}</TabsTrigger>
                        </TabsList>
                        <TabsContent value="linux" className="space-y-2">
                          <div className="rounded-lg border bg-muted/40 p-3">
                            {token ? (
                              <CopyablePopoverContent value={`curl -fsSL "${typeof window !== "undefined" ? window.location.origin : FALLBACK_SERVER_URL}/api/agents/install.sh?token=${token.token}" | bash`} className="font-mono text-xs whitespace-pre-wrap" />
                            ) : (
                              <div className="text-xs text-muted-foreground">{t("install.commandPlaceholder")}</div>
                            )}
                          </div>
                        </TabsContent>
                        <TabsContent value="windows" className="space-y-2">
                          <div className="rounded-lg border bg-muted/40 p-3">
                            {token ? (
                              <CopyablePopoverContent value={`irm "${typeof window !== "undefined" ? window.location.origin : FALLBACK_SERVER_URL}/api/agents/install.ps1?token=${token.token}" | iex`} className="font-mono text-xs whitespace-pre-wrap" />
                            ) : (
                              <div className="text-xs text-muted-foreground">{t("install.commandPlaceholder")}</div>
                            )}
                          </div>
                        </TabsContent>
                      </Tabs>
                    </div>
                  </div>

                  <div className="space-y-3">
                    <div className="rounded-lg border bg-muted/30 p-3">
                      <p className="text-sm font-medium mb-2">{t("install.stepsTitle")}</p>
                      <div className="grid gap-2 text-xs text-muted-foreground">
                        <div className="flex items-center gap-2">
                          <span className="rounded-full border bg-background px-2 py-0.5 text-[10px] font-medium text-foreground">
                            {t("install.step1Label")}
                          </span>
                          <span>{t("install.step1Desc")}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="rounded-full border bg-background px-2 py-0.5 text-[10px] font-medium text-foreground">
                            {t("install.step2Label")}
                          </span>
                          <span>{t("install.step2Desc")}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="rounded-full border bg-background px-2 py-0.5 text-[10px] font-medium text-foreground">
                            {t("install.step3Label")}
                          </span>
                          <span>{t("install.step3Desc")}</span>
                        </div>
                      </div>
                    </div>

                    <div className="rounded-lg border bg-muted/30 p-3">
                      <p className="text-sm font-medium mb-2">{t("install.requirementsTitle")}</p>
                      <div className="grid gap-2 text-xs text-muted-foreground">
                        <div className="flex items-center gap-2">
                          <span className="h-1.5 w-1.5 rounded-full bg-muted-foreground/60" />
                          <span>{t("install.requirementsDocker")}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="h-1.5 w-1.5 rounded-full bg-muted-foreground/60" />
                          <span>{t("install.requirementsAccess")}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </DialogContent>
            </Dialog>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[...Array(3)].map((_, i) => (
                <Skeleton key={i} className="h-52 w-full rounded-lg" />
              ))}
            </div>
          ) : !hasAgents ? (
            <EmptyState onOpenInstall={() => setInstallOpen(true)} />
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {agents.map((agent) => (
                <AgentCardCompact
                  key={agent.id}
                  agent={agent}
                  onConfig={handleConfigure}
                  onDelete={handleDelete}
                />
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <AgentConfigDialog
        open={configDialogOpen}
        onOpenChange={setConfigDialogOpen}
        agent={selectedAgent}
      />

      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title={t("actions.deleteTitle")}
        description={t("actions.deleteDesc", { name: agentToDelete?.name ?? "" })}
        onConfirm={confirmDelete}
        variant="destructive"
        loading={deleteAgent.isPending}
      />

    </div>
  )
}
