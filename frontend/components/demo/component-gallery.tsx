"use client"

import React from "react"
import type { ComponentGroup } from "@/components/demo/component-index"
import { Link } from "@/i18n/navigation"
import { cn } from "@/lib/utils"
import { toast } from "sonner"
import type { ColumnDef } from "@tanstack/react-table"
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts"

import { PageHeader } from "@/components/common/page-header"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Checkbox } from "@/components/ui/checkbox"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Switch } from "@/components/ui/switch"
import { Toggle } from "@/components/ui/toggle"
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { ConfirmDialog } from "@/components/ui/confirm-dialog"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer"
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
  CommandShortcut,
} from "@/components/ui/command"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Separator } from "@/components/ui/separator"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { CardGridSkeleton } from "@/components/ui/card-grid-skeleton"
import { MasterDetailSkeleton } from "@/components/ui/master-detail-skeleton"
import { DataTableSkeleton } from "@/components/ui/data-table-skeleton"
import { Skeleton } from "@/components/ui/skeleton"
import { Spinner } from "@/components/ui/spinner"
import { Progress } from "@/components/ui/progress"
import { ShieldLoader } from "@/components/ui/shield-loader"
import { WaveGrid } from "@/components/ui/wave-grid"
import { Calendar } from "@/components/ui/calendar"
import { DateTimePicker } from "@/components/ui/datetime-picker"
import { Dropzone, DropzoneContent, DropzoneEmptyState } from "@/components/ui/dropzone"
import { CopyablePopoverContent } from "@/components/ui/copyable-popover-content"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { UnifiedDataTable } from "@/components/ui/data-table"
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"
import { MermaidDiagram } from "@/components/ui/mermaid-diagram"
import { Terminal, TypingAnimation } from "@/components/ui/terminal"
import { TerminalLogin } from "@/components/ui/terminal-login"
import { YamlEditor } from "@/components/ui/yaml-editor"
import { Toaster } from "@/components/ui/sonner"
import { Banner, BannerAction, BannerContent, BannerDescription, BannerIcon, BannerTitle } from "@/components/ui/shadcn-io/banner"
import { Status, StatusIndicator, StatusLabel } from "@/components/ui/shadcn-io/status"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
} from "@/components/ui/sidebar"
import {
  AlertTriangle,
  ChevronRight,
  Info,
  Search,
  Settings,
  Wrench,
} from "@/components/icons"

type ComponentGalleryProps = {
  componentGroups: ComponentGroup[]
}

type DemoCardProps = {
  title: string
  description?: string
  children: React.ReactNode
  className?: string
}

const DemoCard = ({ title, description, children, className }: DemoCardProps) => (
  <Card className={cn("border-border/70", className)}>
    <CardHeader className="gap-1">
      <CardTitle className="text-base">{title}</CardTitle>
      {description ? (
        <CardDescription className="text-xs">{description}</CardDescription>
      ) : null}
    </CardHeader>
    <CardContent className="space-y-3">{children}</CardContent>
  </Card>
)

const DemoSection = ({
  id,
  title,
  description,
  children,
}: {
  id?: string
  title: string
  description?: string
  children: React.ReactNode
}) => (
  <section id={id} className="space-y-4">
    <div className="px-4 lg:px-6">
      <h2 className="text-lg font-semibold tracking-tight">{title}</h2>
      {description ? (
        <p className="mt-1 text-sm text-muted-foreground">{description}</p>
      ) : null}
    </div>
    {children}
  </section>
)

const demoRoutes: Record<string, string> = {
  root: "/dashboard/",
  common: "/prototypes/header-demo-a/",
  dashboard: "/dashboard/",
  scan: "/scan/",
  target: "/target/",
  organization: "/organization/",
  vulnerabilities: "/vulnerabilities/",
  search: "/search/",
  tools: "/tools/",
  settings: "/settings/workers/",
  fingerprints: "/tools/fingerprints/",
  endpoints: "/target/",
  directories: "/target/",
  websites: "/target/",
  subdomains: "/target/",
  "ip-addresses": "/target/",
  notifications: "/settings/notifications/",
  auth: "/login/",
  providers: "/dashboard/",
  "animate-ui": "/prototypes/dashboard-demo/",
  prototypes: "/prototypes/scan-dialogs/",
  disk: "/settings/database-health/",
  screenshots: "/target/",
}

const sampleTableData = [
  { id: "1", name: "Alpha Node", status: "online", owner: "Core" },
  { id: "2", name: "Beta Node", status: "degraded", owner: "Edge" },
  { id: "3", name: "Gamma Node", status: "maintenance", owner: "Ops" },
]

const sampleTableColumns: ColumnDef<(typeof sampleTableData)[number]>[] = [
  { accessorKey: "name", header: "Name" },
  { accessorKey: "status", header: "Status" },
  { accessorKey: "owner", header: "Owner" },
]

const chartData = [
  { name: "Mon", value: 40 },
  { name: "Tue", value: 62 },
  { name: "Wed", value: 51 },
  { name: "Thu", value: 78 },
  { name: "Fri", value: 55 },
  { name: "Sat", value: 92 },
  { name: "Sun", value: 68 },
]

const chartConfig = {
  value: {
    label: "Risk",
    color: "var(--color-primary)",
  },
}

const mermaidChart = `flowchart LR
  A[Recon] --> B{Scan Engine}
  B -->|Fast| C[Quick Scan]
  B -->|Deep| D[Full Scan]
  C --> E[Assets]
  D --> E[Assets]
  E --> F[Risk Report]
`

const terminalTranslations = {
  title: "LunaFox Access Terminal",
  subtitle: "Secure access handshake",
  usernamePrompt: "USERNAME",
  passwordPrompt: "PASSWORD",
  authenticating: "AUTHENTICATING",
  processing: "PROCESSING",
  accessGranted: "ACCESS GRANTED",
  welcomeMessage: "Welcome back, operator.",
  authFailed: "AUTH FAILED",
  invalidCredentials: "Invalid credentials",
  shortcuts: "Shortcuts",
  submit: "Enter",
  cancel: "Ctrl+C",
  clear: "Ctrl+U",
  startEnd: "Ctrl+A",
}

export function ComponentGallery({ componentGroups }: ComponentGalleryProps) {
  const [confirmOpen, setConfirmOpen] = React.useState(false)
  const [dialogOpen, setDialogOpen] = React.useState(false)
  const [drawerOpen, setDrawerOpen] = React.useState(false)
  const [sheetOpen, setSheetOpen] = React.useState(false)
  const [calendarDate, setCalendarDate] = React.useState<Date | undefined>(new Date())
  const [dateTime, setDateTime] = React.useState<Date | undefined>(new Date())
  const [yamlValue, setYamlValue] = React.useState("targets:\\n  - example.com\\nscan:\\n  mode: quick\\n")
  const [dropFiles, setDropFiles] = React.useState<File[]>([])
  const [checkboxValue, setCheckboxValue] = React.useState(false)
  const [radioValue, setRadioValue] = React.useState("a")
  const [switchValue, setSwitchValue] = React.useState(true)
  const [toggleValue, setToggleValue] = React.useState(false)
  const [toggleGroupValue, setToggleGroupValue] = React.useState<string[]>(["bold"])
  const [selectValue, setSelectValue] = React.useState("alpha")

  return (
    <div className="flex flex-col gap-8 py-6">
      <Toaster />
      <PageHeader
        code="COMP-LAB"
        title="组件总览"
        description="覆盖 UI 基础组件与业务模块组件的统一演示入口。业务模块默认建议在 Mock 模式查看完整数据展示。"
      />

      <DemoSection
        id="ui-basics"
        title="UI 基础组件"
        description="按钮、表单、选择器与基础元素。"
      >
        <div className="grid gap-4 px-4 lg:px-6 md:grid-cols-2 xl:grid-cols-3">
          <DemoCard title="Button" description="基础按钮与状态">
            <div className="flex flex-wrap gap-2">
              <Button>Primary</Button>
              <Button variant="outline">Outline</Button>
              <Button variant="secondary">Secondary</Button>
              <Button variant="destructive">Destructive</Button>
              <Button size="sm">Small</Button>
            </div>
          </DemoCard>

          <DemoCard title="Badge" description="状态与标签">
            <div className="flex flex-wrap gap-2">
              <Badge>Default</Badge>
              <Badge variant="secondary">Secondary</Badge>
              <Badge variant="outline">Outline</Badge>
              <Badge className="bg-[var(--success)]/10 text-[var(--success)] border-[var(--success)]/30" variant="outline">
                Healthy
              </Badge>
            </div>
          </DemoCard>

          <DemoCard title="Input + Label" description="基础输入">
            <div className="space-y-2">
              <Label htmlFor="demo-input">资产名称</Label>
              <Input id="demo-input" placeholder="example.com" />
            </div>
          </DemoCard>

          <DemoCard title="Textarea" description="多行输入">
            <Textarea placeholder="输入说明..." rows={4} />
          </DemoCard>

          <DemoCard title="Checkbox" description="多选">
            <div className="flex items-center gap-2">
              <Checkbox
                id="demo-checkbox"
                checked={checkboxValue}
                onCheckedChange={(value) => setCheckboxValue(Boolean(value))}
              />
              <Label htmlFor="demo-checkbox">启用深度扫描</Label>
            </div>
          </DemoCard>

          <DemoCard title="Radio Group" description="单选">
            <RadioGroup value={radioValue} onValueChange={setRadioValue} className="gap-2">
              <div className="flex items-center gap-2">
                <RadioGroupItem value="a" id="radio-a" />
                <Label htmlFor="radio-a">Quick</Label>
              </div>
              <div className="flex items-center gap-2">
                <RadioGroupItem value="b" id="radio-b" />
                <Label htmlFor="radio-b">Deep</Label>
              </div>
            </RadioGroup>
          </DemoCard>

          <DemoCard title="Switch" description="开关">
            <div className="flex items-center gap-3">
              <Switch checked={switchValue} onCheckedChange={setSwitchValue} />
              <span className="text-sm text-muted-foreground">
                {switchValue ? "实时监控开启" : "实时监控关闭"}
              </span>
            </div>
          </DemoCard>

          <DemoCard title="Toggle / ToggleGroup" description="强调型切换">
            <div className="flex flex-col gap-3">
              <Toggle pressed={toggleValue} onPressedChange={setToggleValue}>
                单个 Toggle
              </Toggle>
              <ToggleGroup
                type="multiple"
                value={toggleGroupValue}
                onValueChange={setToggleGroupValue}
                className="flex flex-wrap gap-2"
              >
                <ToggleGroupItem value="bold">Bold</ToggleGroupItem>
                <ToggleGroupItem value="italic">Italic</ToggleGroupItem>
                <ToggleGroupItem value="mono">Mono</ToggleGroupItem>
              </ToggleGroup>
            </div>
          </DemoCard>

          <DemoCard title="Select" description="下拉选择">
            <Select value={selectValue} onValueChange={setSelectValue}>
              <SelectTrigger>
                <SelectValue placeholder="选择扫描策略" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="alpha">Alpha</SelectItem>
                <SelectItem value="beta">Beta</SelectItem>
                <SelectItem value="gamma">Gamma</SelectItem>
              </SelectContent>
            </Select>
          </DemoCard>

          <DemoCard title="DateTimePicker" description="日期时间选择">
            <DateTimePicker value={dateTime} onChange={setDateTime} />
          </DemoCard>

          <DemoCard title="Calendar" description="日期日历">
            <Calendar mode="single" selected={calendarDate} onSelect={setCalendarDate} />
          </DemoCard>

          <DemoCard title="Dropzone" description="上传拖拽区">
            <Dropzone
              src={dropFiles}
              maxFiles={3}
              onDrop={(files) => setDropFiles(files)}
              className="border-dashed"
            >
              <DropzoneEmptyState />
              <DropzoneContent />
            </Dropzone>
          </DemoCard>

          <DemoCard title="Avatar" description="用户头像">
            <div className="flex items-center gap-3">
              <Avatar>
                <AvatarImage src="/images/icon-64.png" alt="User" />
                <AvatarFallback>LF</AvatarFallback>
              </Avatar>
              <span className="text-sm text-muted-foreground">Operator</span>
            </div>
          </DemoCard>
        </div>
      </DemoSection>

      <DemoSection
        id="ui-overlays"
        title="反馈与浮层"
        description="弹窗、提示、反馈与状态组件。"
      >
        <div className="grid gap-4 px-4 lg:px-6 md:grid-cols-2 xl:grid-cols-3">
          <DemoCard title="Alert" description="信息提示">
            <Alert>
              <AlertTriangle className="size-4" />
              <AlertTitle>扫描警告</AlertTitle>
              <AlertDescription>目标存在高危端口，建议开启深度扫描。</AlertDescription>
            </Alert>
          </DemoCard>

          <DemoCard title="AlertDialog" description="确认型弹窗">
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button variant="outline">打开 AlertDialog</Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>确定执行删除？</AlertDialogTitle>
                  <AlertDialogDescription>该操作不可撤销。</AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>取消</AlertDialogCancel>
                  <AlertDialogAction>确认</AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </DemoCard>

          <DemoCard title="ConfirmDialog" description="业务确认弹窗">
            <div className="space-y-2">
              <Button onClick={() => setConfirmOpen(true)}>打开 ConfirmDialog</Button>
              <ConfirmDialog
                open={confirmOpen}
                onOpenChange={setConfirmOpen}
                title="提交扫描任务"
                description="提交后将进入队列。"
                onConfirm={() => setConfirmOpen(false)}
              />
            </div>
          </DemoCard>

          <DemoCard title="Dialog" description="通用对话框">
            <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
              <DialogTrigger asChild>
                <Button variant="outline">打开 Dialog</Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>新建目标</DialogTitle>
                  <DialogDescription>快速添加一个新的资产目标。</DialogDescription>
                </DialogHeader>
                <Input placeholder="example.com" />
                <DialogFooter>
                  <Button onClick={() => setDialogOpen(false)}>保存</Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </DemoCard>

          <DemoCard title="Sheet" description="侧边面板">
            <Sheet open={sheetOpen} onOpenChange={setSheetOpen}>
              <SheetTrigger asChild>
                <Button variant="outline">打开 Sheet</Button>
              </SheetTrigger>
              <SheetContent side="right">
                <SheetHeader>
                  <SheetTitle>配置面板</SheetTitle>
                  <SheetDescription>快速调整扫描策略。</SheetDescription>
                </SheetHeader>
                <div className="mt-4 space-y-2">
                  <Label>策略名称</Label>
                  <Input placeholder="Default" />
                </div>
              </SheetContent>
            </Sheet>
          </DemoCard>

          <DemoCard title="Drawer" description="抽屉式交互">
            <Drawer open={drawerOpen} onOpenChange={setDrawerOpen}>
              <DrawerTrigger asChild>
                <Button variant="outline">打开 Drawer</Button>
              </DrawerTrigger>
              <DrawerContent>
                <DrawerHeader>
                  <DrawerTitle>快速执行</DrawerTitle>
                  <DrawerDescription>启动一个快速扫描任务。</DrawerDescription>
                </DrawerHeader>
                <DrawerFooter className="pb-6">
                  <Button onClick={() => setDrawerOpen(false)}>启动</Button>
                </DrawerFooter>
              </DrawerContent>
            </Drawer>
          </DemoCard>

          <DemoCard title="Popover" description="浮层内容">
            <Popover>
              <PopoverTrigger asChild>
                <Button variant="outline">打开 Popover</Button>
              </PopoverTrigger>
              <PopoverContent className="w-64">
                <CopyablePopoverContent value="https://example.com/asset/alpha" />
              </PopoverContent>
            </Popover>
          </DemoCard>

          <DemoCard title="HoverCard" description="悬浮卡片">
            <HoverCard>
              <HoverCardTrigger asChild>
                <Button variant="ghost">Hover 预览</Button>
              </HoverCardTrigger>
              <HoverCardContent className="w-64">
                <div className="space-y-1">
                  <p className="text-sm font-medium">资产状态</p>
                  <p className="text-xs text-muted-foreground">最近一次扫描：3 分钟前</p>
                </div>
              </HoverCardContent>
            </HoverCard>
          </DemoCard>

          <DemoCard title="Tooltip" description="轻提示">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="outline">提示</Button>
                </TooltipTrigger>
                <TooltipContent>这是一个 Tooltip</TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </DemoCard>

          <DemoCard title="DropdownMenu" description="下拉菜单">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline">操作</Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuLabel>快速操作</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem>开始扫描</DropdownMenuItem>
                <DropdownMenuItem>查看详情</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </DemoCard>

          <DemoCard title="Toast (Sonner)" description="轻量通知">
            <div className="flex gap-2">
              <Button
                variant="outline"
                onClick={() => toast.success("任务已提交")}
              >
                Success
              </Button>
              <Button
                variant="outline"
                onClick={() => toast.error("请求失败")}
              >
                Error
              </Button>
            </div>
          </DemoCard>

          <DemoCard title="Progress + Spinner" description="进度指示">
            <div className="space-y-3">
              <Progress value={65} />
              <div className="flex items-center gap-2">
                <Spinner />
                <span className="text-xs text-muted-foreground">处理中…</span>
              </div>
            </div>
          </DemoCard>

          <DemoCard title="Skeletons" description="骨架屏">
            <div className="space-y-3">
              <Skeleton className="h-6 w-1/2" />
              <CardGridSkeleton />
            </div>
          </DemoCard>

          <DemoCard title="Master Detail Skeleton" description="主从骨架">
            <MasterDetailSkeleton />
          </DemoCard>

          <DemoCard title="Data Table Skeleton" description="表格骨架">
            <DataTableSkeleton />
          </DemoCard>

          <DemoCard title="Shield Loader" description="工业风加载">
            <div className="flex justify-center">
              <ShieldLoader />
            </div>
          </DemoCard>

          <DemoCard title="Wave Grid" description="网格波纹">
            <WaveGrid />
          </DemoCard>
        </div>
      </DemoSection>

      <DemoSection
        id="ui-layout"
        title="布局与数据展示"
        description="容器、表格、命令面板与侧边栏。"
      >
        <div className="grid gap-4 px-4 lg:px-6 md:grid-cols-2 xl:grid-cols-3">
          <DemoCard title="Card" description="内容容器">
            <Card className="border-dashed">
              <CardHeader>
                <CardTitle className="text-sm">资产概览</CardTitle>
                <CardDescription>示例描述文本</CardDescription>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                这里是卡片内容区域。
              </CardContent>
            </Card>
          </DemoCard>

          <DemoCard title="Tabs" description="标签页">
            <Tabs defaultValue="overview">
              <TabsList>
                <TabsTrigger value="overview">Overview</TabsTrigger>
                <TabsTrigger value="assets">Assets</TabsTrigger>
              </TabsList>
              <TabsContent value="overview" className="text-sm text-muted-foreground">
                Overview 内容示例。
              </TabsContent>
              <TabsContent value="assets" className="text-sm text-muted-foreground">
                Assets 内容示例。
              </TabsContent>
            </Tabs>
          </DemoCard>

          <DemoCard title="ScrollArea" description="滚动容器">
            <ScrollArea className="h-28 rounded-md border p-2">
              <div className="space-y-2 text-sm">
                {Array.from({ length: 6 }).map((_, index) => (
                  <div key={index} className="flex items-center justify-between">
                    <span>任务 #{index + 1}</span>
                    <Badge variant="outline">Queued</Badge>
                  </div>
                ))}
              </div>
            </ScrollArea>
          </DemoCard>

          <DemoCard title="Separator" description="分割线">
            <div className="space-y-2">
              <div className="text-sm">上方内容</div>
              <Separator />
              <div className="text-sm text-muted-foreground">下方内容</div>
            </div>
          </DemoCard>

          <DemoCard title="Collapsible" description="折叠内容">
            <Collapsible>
              <CollapsibleTrigger asChild>
                <Button variant="outline">
                  展开更多
                  <ChevronRight className="ml-2 size-4" />
                </Button>
              </CollapsibleTrigger>
              <CollapsibleContent className="mt-2 text-sm text-muted-foreground">
                这里是折叠内容区域。
              </CollapsibleContent>
            </Collapsible>
          </DemoCard>

          <DemoCard title="Command" description="命令面板">
            <Command>
              <CommandInput placeholder="搜索指令..." />
              <CommandList>
                <CommandEmpty>无结果</CommandEmpty>
                <CommandGroup heading="导航">
                  <CommandItem>
                    <Search className="mr-2 size-4" />
                    搜索资产
                    <CommandShortcut>⌘K</CommandShortcut>
                  </CommandItem>
                  <CommandItem>
                    <Settings className="mr-2 size-4" />
                    系统设置
                  </CommandItem>
                </CommandGroup>
                <CommandSeparator />
                <CommandGroup heading="工具">
                  <CommandItem>
                    <Tool className="mr-2 size-4" />
                    快速扫描
                  </CommandItem>
                </CommandGroup>
              </CommandList>
            </Command>
          </DemoCard>

          <DemoCard title="Table" description="基础表格">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>任务</TableHead>
                  <TableHead>状态</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow>
                  <TableCell>Recon</TableCell>
                  <TableCell>Running</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>Scan</TableCell>
                  <TableCell>Queued</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </DemoCard>

          <DemoCard title="UnifiedDataTable" description="统一数据表格">
            <UnifiedDataTable
              data={sampleTableData}
              columns={sampleTableColumns}
              hideToolbar
              hidePagination
              enableRowSelection={false}
            />
          </DemoCard>

          <DemoCard title="Sidebar" description="侧边栏骨架">
            <SidebarProvider>
              <div className="flex h-48 border rounded-md overflow-hidden">
                <Sidebar className="w-40">
                  <SidebarContent>
                    <SidebarGroup>
                      <SidebarGroupLabel>导航</SidebarGroupLabel>
                      <SidebarGroupContent>
                        <SidebarMenu>
                          <SidebarMenuItem>
                            <SidebarMenuButton>Dashboard</SidebarMenuButton>
                          </SidebarMenuItem>
                          <SidebarMenuItem>
                            <SidebarMenuButton>Scan</SidebarMenuButton>
                          </SidebarMenuItem>
                        </SidebarMenu>
                      </SidebarGroupContent>
                    </SidebarGroup>
                  </SidebarContent>
                </Sidebar>
                <div className="flex-1 p-3 text-xs text-muted-foreground">
                  Sidebar 预览区域
                </div>
              </div>
            </SidebarProvider>
          </DemoCard>
        </div>
      </DemoSection>

      <DemoSection
        id="ui-visual"
        title="可视化与高级组件"
        description="图表、终端、编辑器等高级组件。"
      >
        <div className="grid gap-4 px-4 lg:px-6 md:grid-cols-2 xl:grid-cols-3">
          <DemoCard title="Chart" description="Recharts 组合">
            <ChartContainer config={chartConfig} className="h-40">
              <AreaChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <ChartTooltip content={<ChartTooltipContent />} />
                <Area type="monotone" dataKey="value" stroke="var(--color-primary)" fill="var(--color-primary)" fillOpacity={0.2} />
                <ChartLegend content={<ChartLegendContent />} />
              </AreaChart>
            </ChartContainer>
          </DemoCard>

          <DemoCard title="MermaidDiagram" description="流程图渲染">
            <MermaidDiagram chart={mermaidChart} />
          </DemoCard>

          <DemoCard title="Terminal" description="终端动效">
            <Terminal className="max-w-full">
              <TypingAnimation>lunafox init --mode stealth</TypingAnimation>
              <TypingAnimation delay={800}>fetching assets...</TypingAnimation>
              <TypingAnimation delay={1400}>scan running...</TypingAnimation>
            </Terminal>
          </DemoCard>

          <DemoCard title="TerminalLogin" description="登录终端">
            <TerminalLogin
              translations={terminalTranslations}
              onLogin={async () => toast.success("认证完成")}
            />
          </DemoCard>

          <DemoCard title="YamlEditor" description="YAML 编辑器">
            <div className="h-48">
              <YamlEditor value={yamlValue} onChange={setYamlValue} />
            </div>
          </DemoCard>
        </div>
      </DemoSection>

      <DemoSection
        id="ui-branding"
        title="品牌与状态组件"
        description="Banner 与状态标识。"
      >
        <div className="grid gap-4 px-4 lg:px-6 md:grid-cols-2 xl:grid-cols-3">
          <DemoCard title="Banner" description="顶部提示横条">
            <Banner>
              <BannerIcon icon={Info} />
              <BannerContent>
                <BannerTitle>系统维护</BannerTitle>
                <BannerDescription>预计 5 分钟后恢复。</BannerDescription>
              </BannerContent>
              <BannerAction>查看详情</BannerAction>
            </Banner>
          </DemoCard>

          <DemoCard title="Status" description="状态指示">
            <div className="flex flex-col gap-2">
              <Status status="online">
                <StatusIndicator />
                <StatusLabel />
              </Status>
              <Status status="degraded">
                <StatusIndicator />
                <StatusLabel />
              </Status>
            </div>
          </DemoCard>
        </div>
      </DemoSection>

      <DemoSection
        id="business-index"
        title="业务组件索引"
        description="以下为业务模块组件清单与入口。建议在 Mock 模式查看完整数据交互。"
      >
        <div className="grid gap-4 px-4 lg:px-6 md:grid-cols-2">
          {componentGroups.map((group) => {
            const items = group.items
            const preview = items.slice(0, 14)
            const rest = items.length - preview.length
            const route = demoRoutes[group.key]

            return (
              <Card key={group.key} className="border-border/60">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2 text-base">
                    <span>{group.title}</span>
                    {route ? (
                      <Badge variant="outline" className="text-[10px]">
                        {route}
                      </Badge>
                    ) : null}
                  </CardTitle>
                  {group.description ? (
                    <CardDescription className="text-xs">{group.description}</CardDescription>
                  ) : null}
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex flex-wrap gap-2">
                    {preview.map((item) => (
                      <Badge key={item} variant="outline" className="text-[10px]">
                        {item}
                      </Badge>
                    ))}
                    {rest > 0 ? (
                      <Badge variant="secondary" className="text-[10px]">
                        +{rest}
                      </Badge>
                    ) : null}
                  </div>
                  {route ? (
                    <Button asChild variant="outline" size="sm">
                      <Link href={route}>进入模块</Link>
                    </Button>
                  ) : null}
                </CardContent>
              </Card>
            )
          })}
        </div>
      </DemoSection>
    </div>
  )
}
