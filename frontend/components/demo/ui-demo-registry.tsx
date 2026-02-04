"use client"

import React from "react"
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts"
import type { ColumnDef } from "@tanstack/react-table"
import { useForm } from "react-hook-form"

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
import { Banner, BannerAction, BannerIcon, BannerTitle } from "@/components/ui/shadcn-io/banner"
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
import * as Icons from "@/components/icons"
import { iconNames } from "@/components/demo/icons-list"
import { toast } from "sonner"
import { Field, FieldContent, FieldDescription, FieldGroup, FieldLabel, FieldLegend, FieldSet } from "@/components/ui/field"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"

export type DemoItem = {
  slug: string
  title: string
  description?: string
  group: string
  Demo: React.ComponentType
}

const DemoShell = ({ title, description, children }: { title: string; description?: string; children: React.ReactNode }) => (
  <div className="flex flex-col gap-4 py-6">
    <PageHeader code="UI-DEMO" title={title} description={description} />
    <div className="px-4 lg:px-6">{children}</div>
  </div>
)

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

const IconDemo = () => {
  const [query, setQuery] = React.useState("")
  const items = iconNames.filter((name) => name.toLowerCase().includes(query.toLowerCase()))
  return (
    <DemoShell title="Icon" description="Carbon 图标集合（支持搜索）">
      <div className="mb-4 max-w-md">
        <Input placeholder="搜索图标..." value={query} onChange={(e) => setQuery(e.target.value)} />
      </div>
      <div className="grid gap-3 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        {items.map((name) => {
          const IconComp = (Icons as Record<string, React.ComponentType<any>>)[name]
          return (
            <div key={name} className="flex items-center gap-2 rounded-md border p-2 text-xs">
              {IconComp ? <IconComp className="size-4" /> : null}
              <span className="truncate">{name}</span>
            </div>
          )
        })}
      </div>
    </DemoShell>
  )
}

const ButtonDemo = () => (
  <DemoShell title="Button" description="按钮样式与状态">
    <div className="flex flex-wrap gap-2">
      <Button>Primary</Button>
      <Button variant="secondary">Secondary</Button>
      <Button variant="outline">Outline</Button>
      <Button variant="destructive">Destructive</Button>
      <Button size="sm">Small</Button>
      <Button size="lg">Large</Button>
    </div>
  </DemoShell>
)

const BadgeDemo = () => (
  <DemoShell title="Badge" description="标签与状态">
    <div className="flex flex-wrap gap-2">
      <Badge>Default</Badge>
      <Badge variant="secondary">Secondary</Badge>
      <Badge variant="outline">Outline</Badge>
      <Badge className="bg-[var(--success)]/10 text-[var(--success)] border-[var(--success)]/30" variant="outline">
        Healthy
      </Badge>
    </div>
  </DemoShell>
)

const InputDemo = () => (
  <DemoShell title="Input" description="基础输入">
    <div className="space-y-2 max-w-sm">
      <Label htmlFor="demo-input">资产名称</Label>
      <Input id="demo-input" placeholder="example.com" />
    </div>
  </DemoShell>
)

const LabelDemo = () => (
  <DemoShell title="Label" description="文本标签">
    <div className="flex flex-col gap-2">
      <Label>默认标签</Label>
      <Label className="text-muted-foreground">次级标签</Label>
    </div>
  </DemoShell>
)

const TextareaDemo = () => (
  <DemoShell title="Textarea" description="多行输入">
    <div className="max-w-md">
      <Textarea placeholder="输入说明..." rows={5} />
    </div>
  </DemoShell>
)

const CheckboxDemo = () => {
  const [checked, setChecked] = React.useState(false)
  return (
    <DemoShell title="Checkbox" description="多选">
      <div className="flex items-center gap-2">
        <Checkbox id="demo-checkbox" checked={checked} onCheckedChange={(value) => setChecked(Boolean(value))} />
        <Label htmlFor="demo-checkbox">启用深度扫描</Label>
      </div>
    </DemoShell>
  )
}

const RadioDemo = () => {
  const [value, setValue] = React.useState("a")
  return (
    <DemoShell title="RadioGroup" description="单选">
      <RadioGroup value={value} onValueChange={setValue} className="gap-2">
        <div className="flex items-center gap-2">
          <RadioGroupItem value="a" id="radio-a" />
          <Label htmlFor="radio-a">Quick</Label>
        </div>
        <div className="flex items-center gap-2">
          <RadioGroupItem value="b" id="radio-b" />
          <Label htmlFor="radio-b">Deep</Label>
        </div>
      </RadioGroup>
    </DemoShell>
  )
}

const SwitchDemo = () => {
  const [value, setValue] = React.useState(true)
  return (
    <DemoShell title="Switch" description="开关">
      <div className="flex items-center gap-3">
        <Switch checked={value} onCheckedChange={setValue} />
        <span className="text-sm text-muted-foreground">{value ? "实时监控开启" : "实时监控关闭"}</span>
      </div>
    </DemoShell>
  )
}

const ToggleDemo = () => {
  const [value, setValue] = React.useState(false)
  return (
    <DemoShell title="Toggle" description="单个切换">
      <Toggle pressed={value} onPressedChange={setValue}>
        单个 Toggle
      </Toggle>
    </DemoShell>
  )
}

const ToggleGroupDemo = () => {
  const [value, setValue] = React.useState<string[]>(["bold"])
  return (
    <DemoShell title="ToggleGroup" description="多选切换">
      <ToggleGroup type="multiple" value={value} onValueChange={setValue} className="flex flex-wrap gap-2">
        <ToggleGroupItem value="bold">Bold</ToggleGroupItem>
        <ToggleGroupItem value="italic">Italic</ToggleGroupItem>
        <ToggleGroupItem value="mono">Mono</ToggleGroupItem>
      </ToggleGroup>
    </DemoShell>
  )
}

const SelectDemo = () => {
  const [value, setValue] = React.useState("alpha")
  return (
    <DemoShell title="Select" description="下拉选择">
      <div className="max-w-sm">
        <Select value={value} onValueChange={setValue}>
          <SelectTrigger>
            <SelectValue placeholder="选择扫描策略" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="alpha">Alpha</SelectItem>
            <SelectItem value="beta">Beta</SelectItem>
            <SelectItem value="gamma">Gamma</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </DemoShell>
  )
}

const CalendarDemo = () => {
  const [date, setDate] = React.useState<Date | undefined>(new Date())
  return (
    <DemoShell title="Calendar" description="日期选择">
      <Calendar mode="single" selected={date} onSelect={setDate} />
    </DemoShell>
  )
}

const DateTimePickerDemo = () => {
  const [date, setDate] = React.useState<Date | undefined>(new Date())
  return (
    <DemoShell title="DateTimePicker" description="日期时间选择">
      <DateTimePicker value={date} onChange={setDate} />
    </DemoShell>
  )
}

const DropzoneDemo = () => {
  const [files, setFiles] = React.useState<File[]>([])
  return (
    <DemoShell title="Dropzone" description="上传拖拽区">
      <Dropzone src={files} maxFiles={3} onDrop={(next) => setFiles(next)} className="border-dashed">
        <DropzoneEmptyState />
        <DropzoneContent />
      </Dropzone>
    </DemoShell>
  )
}

const AvatarDemo = () => (
  <DemoShell title="Avatar" description="头像">
    <div className="flex items-center gap-3">
      <Avatar>
        <AvatarImage src="/images/icon-64.png" alt="User" />
        <AvatarFallback>LF</AvatarFallback>
      </Avatar>
      <span className="text-sm text-muted-foreground">Operator</span>
    </div>
  </DemoShell>
)

const AlertDemo = () => (
  <DemoShell title="Alert" description="提示块">
    <Alert>
      <AlertTriangle className="size-4" />
      <AlertTitle>扫描警告</AlertTitle>
      <AlertDescription>目标存在高危端口，建议开启深度扫描。</AlertDescription>
    </Alert>
  </DemoShell>
)

const AlertDialogDemo = () => (
  <DemoShell title="AlertDialog" description="确认弹窗">
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
  </DemoShell>
)

const ConfirmDialogDemo = () => {
  const [open, setOpen] = React.useState(false)
  return (
    <DemoShell title="ConfirmDialog" description="业务确认弹窗">
      <Button onClick={() => setOpen(true)}>打开 ConfirmDialog</Button>
      <ConfirmDialog
        open={open}
        onOpenChange={setOpen}
        title="提交扫描任务"
        description="提交后将进入队列。"
        onConfirm={() => setOpen(false)}
      />
    </DemoShell>
  )
}

const DialogDemo = () => (
  <DemoShell title="Dialog" description="通用对话框">
    <Dialog>
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
          <Button>保存</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </DemoShell>
)

const DrawerDemo = () => (
  <DemoShell title="Drawer" description="抽屉式交互">
    <Drawer>
      <DrawerTrigger asChild>
        <Button variant="outline">打开 Drawer</Button>
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>快速执行</DrawerTitle>
          <DrawerDescription>启动一个快速扫描任务。</DrawerDescription>
        </DrawerHeader>
        <DrawerFooter className="pb-6">
          <Button>启动</Button>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  </DemoShell>
)

const SheetDemo = () => (
  <DemoShell title="Sheet" description="侧边面板">
    <Sheet>
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
  </DemoShell>
)

const PopoverDemo = () => (
  <DemoShell title="Popover" description="浮层内容">
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline">打开 Popover</Button>
      </PopoverTrigger>
      <PopoverContent className="w-64">
        <CopyablePopoverContent value="https://example.com/asset/alpha" />
      </PopoverContent>
    </Popover>
  </DemoShell>
)

const CopyablePopoverContentDemo = () => (
  <DemoShell title="CopyablePopoverContent" description="可复制内容">
    <div className="max-w-md rounded-md border p-4">
      <CopyablePopoverContent value="https://example.com/asset/alpha" />
    </div>
  </DemoShell>
)

const HoverCardDemo = () => (
  <DemoShell title="HoverCard" description="悬浮卡片">
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
  </DemoShell>
)

const TooltipDemo = () => (
  <DemoShell title="Tooltip" description="轻提示">
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button variant="outline">提示</Button>
        </TooltipTrigger>
        <TooltipContent>这是一个 Tooltip</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  </DemoShell>
)

const DropdownMenuDemo = () => (
  <DemoShell title="DropdownMenu" description="下拉菜单">
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
  </DemoShell>
)

const CollapsibleDemo = () => (
  <DemoShell title="Collapsible" description="折叠内容">
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
  </DemoShell>
)

const CommandDemo = () => (
  <DemoShell title="Command" description="命令面板">
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
            <Wrench className="mr-2 size-4" />
            快速扫描
          </CommandItem>
        </CommandGroup>
      </CommandList>
    </Command>
  </DemoShell>
)

const TabsDemo = () => (
  <DemoShell title="Tabs" description="标签页">
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
  </DemoShell>
)

const ScrollAreaDemo = () => (
  <DemoShell title="ScrollArea" description="滚动容器">
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
  </DemoShell>
)

const SeparatorDemo = () => (
  <DemoShell title="Separator" description="分割线">
    <div className="space-y-2">
      <div className="text-sm">上方内容</div>
      <Separator />
      <div className="text-sm text-muted-foreground">下方内容</div>
    </div>
  </DemoShell>
)

const TableDemo = () => (
  <DemoShell title="Table" description="基础表格">
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
  </DemoShell>
)

const UnifiedDataTableDemo = () => (
  <DemoShell title="UnifiedDataTable" description="统一数据表格">
    <UnifiedDataTable
      data={sampleTableData}
      columns={sampleTableColumns}
      hideToolbar
      hidePagination
      enableRowSelection={false}
    />
  </DemoShell>
)

const SidebarDemo = () => (
  <DemoShell title="Sidebar" description="侧边栏骨架">
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
  </DemoShell>
)

const CardDemo = () => (
  <DemoShell title="Card" description="卡片容器">
    <Card className="border-dashed max-w-sm">
      <CardHeader>
        <CardTitle className="text-sm">资产概览</CardTitle>
        <CardDescription>示例描述文本</CardDescription>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">
        这里是卡片内容区域。
      </CardContent>
    </Card>
  </DemoShell>
)

const SkeletonDemo = () => (
  <DemoShell title="Skeleton" description="骨架屏">
    <div className="space-y-3 max-w-sm">
      <Skeleton className="h-6 w-1/2" />
      <CardGridSkeleton />
    </div>
  </DemoShell>
)

const MasterDetailSkeletonDemo = () => (
  <DemoShell title="MasterDetailSkeleton" description="主从骨架">
    <MasterDetailSkeleton />
  </DemoShell>
)

const DataTableSkeletonDemo = () => (
  <DemoShell title="DataTableSkeleton" description="表格骨架">
    <DataTableSkeleton />
  </DemoShell>
)

const SpinnerDemo = () => (
  <DemoShell title="Spinner" description="加载指示">
    <div className="flex items-center gap-2">
      <Spinner />
      <span className="text-xs text-muted-foreground">处理中…</span>
    </div>
  </DemoShell>
)

const ProgressDemo = () => (
  <DemoShell title="Progress" description="进度条">
    <Progress value={65} />
  </DemoShell>
)

const ShieldLoaderDemo = () => (
  <DemoShell title="ShieldLoader" description="工业风加载">
    <div className="flex justify-center">
      <ShieldLoader />
    </div>
  </DemoShell>
)

const WaveGridDemo = () => (
  <DemoShell title="WaveGrid" description="网格波纹">
    <WaveGrid />
  </DemoShell>
)

const ChartDemo = () => (
  <DemoShell title="Chart" description="Recharts 组合">
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
  </DemoShell>
)

const MermaidDemo = () => (
  <DemoShell title="MermaidDiagram" description="流程图渲染">
    <MermaidDiagram chart={mermaidChart} />
  </DemoShell>
)

const TerminalDemo = () => (
  <DemoShell title="Terminal" description="终端动效">
    <Terminal className="max-w-full">
      <TypingAnimation>lunafox init --mode stealth</TypingAnimation>
      <TypingAnimation delay={800}>fetching assets...</TypingAnimation>
      <TypingAnimation delay={1400}>scan running...</TypingAnimation>
    </Terminal>
  </DemoShell>
)

const TerminalLoginDemo = () => (
  <DemoShell title="TerminalLogin" description="登录终端">
    <TerminalLogin translations={terminalTranslations} onLogin={async () => toast.success("认证完成")} />
  </DemoShell>
)

const YamlEditorDemo = () => {
  const [value, setValue] = React.useState("targets:\\n  - example.com\\nscan:\\n  mode: quick\\n")
  return (
    <DemoShell title="YamlEditor" description="YAML 编辑器">
      <div className="h-60">
        <YamlEditor value={value} onChange={setValue} />
      </div>
    </DemoShell>
  )
}

const ToastDemo = () => (
  <DemoShell title="Sonner" description="通知">
    <Toaster />
    <div className="flex gap-2">
      <Button variant="outline" onClick={() => toast.success("任务已提交")}>Success</Button>
      <Button variant="outline" onClick={() => toast.error("请求失败")}>Error</Button>
    </div>
  </DemoShell>
)

const BannerDemo = () => (
  <DemoShell title="Banner" description="顶部提示横条">
    <Banner>
      <BannerIcon icon={Info} />
      <div className="flex flex-1 flex-col gap-0.5">
        <BannerTitle>系统维护</BannerTitle>
        <p className="text-xs text-primary-foreground/80">预计 5 分钟后恢复。</p>
      </div>
      <BannerAction>查看详情</BannerAction>
    </Banner>
  </DemoShell>
)

const StatusDemo = () => (
  <DemoShell title="Status" description="状态指示">
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
  </DemoShell>
)

const FieldDemo = () => (
  <DemoShell title="Field" description="表单字段组合">
    <FieldSet>
      <FieldLegend>扫描配置</FieldLegend>
      <FieldGroup>
        <Field>
          <FieldLabel>目标域名</FieldLabel>
          <FieldContent>
            <Input placeholder="example.com" />
            <FieldDescription>支持多个域名</FieldDescription>
          </FieldContent>
        </Field>
      </FieldGroup>
    </FieldSet>
  </DemoShell>
)

const FormDemo = () => {
  const form = useForm<{ name: string }>({ defaultValues: { name: "" } })
  return (
    <DemoShell title="Form" description="React Hook Form 封装">
      <Form {...form}>
        <form className="space-y-3 max-w-sm" onSubmit={form.handleSubmit(() => toast.success("已提交"))}>
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>名称</FormLabel>
                <FormControl>
                  <Input placeholder="输入名称" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit">提交</Button>
        </form>
      </Form>
    </DemoShell>
  )
}

export const uiDemoItems: DemoItem[] = [
  { slug: "icon", title: "Icon", description: "图标展示与搜索", group: "基础", Demo: IconDemo },
  { slug: "button", title: "Button", description: "按钮样式", group: "基础", Demo: ButtonDemo },
  { slug: "badge", title: "Badge", description: "标签", group: "基础", Demo: BadgeDemo },
  { slug: "input", title: "Input", description: "文本输入", group: "表单", Demo: InputDemo },
  { slug: "label", title: "Label", description: "文本标签", group: "表单", Demo: LabelDemo },
  { slug: "textarea", title: "Textarea", description: "多行输入", group: "表单", Demo: TextareaDemo },
  { slug: "checkbox", title: "Checkbox", description: "多选", group: "表单", Demo: CheckboxDemo },
  { slug: "radio-group", title: "RadioGroup", description: "单选", group: "表单", Demo: RadioDemo },
  { slug: "switch", title: "Switch", description: "开关", group: "表单", Demo: SwitchDemo },
  { slug: "toggle", title: "Toggle", description: "单个切换", group: "表单", Demo: ToggleDemo },
  { slug: "toggle-group", title: "ToggleGroup", description: "多选切换", group: "表单", Demo: ToggleGroupDemo },
  { slug: "select", title: "Select", description: "下拉选择", group: "表单", Demo: SelectDemo },
  { slug: "calendar", title: "Calendar", description: "日期选择", group: "表单", Demo: CalendarDemo },
  { slug: "datetime-picker", title: "DateTimePicker", description: "日期时间", group: "表单", Demo: DateTimePickerDemo },
  { slug: "dropzone", title: "Dropzone", description: "文件上传", group: "表单", Demo: DropzoneDemo },
  { slug: "avatar", title: "Avatar", description: "头像", group: "基础", Demo: AvatarDemo },
  { slug: "alert", title: "Alert", description: "提示", group: "反馈", Demo: AlertDemo },
  { slug: "alert-dialog", title: "AlertDialog", description: "确认弹窗", group: "反馈", Demo: AlertDialogDemo },
  { slug: "confirm-dialog", title: "ConfirmDialog", description: "业务确认", group: "反馈", Demo: ConfirmDialogDemo },
  { slug: "dialog", title: "Dialog", description: "通用对话框", group: "反馈", Demo: DialogDemo },
  { slug: "drawer", title: "Drawer", description: "抽屉", group: "反馈", Demo: DrawerDemo },
  { slug: "sheet", title: "Sheet", description: "侧边面板", group: "反馈", Demo: SheetDemo },
  { slug: "popover", title: "Popover", description: "浮层", group: "反馈", Demo: PopoverDemo },
  { slug: "copyable-popover-content", title: "CopyablePopoverContent", description: "可复制内容", group: "反馈", Demo: CopyablePopoverContentDemo },
  { slug: "hover-card", title: "HoverCard", description: "悬浮卡", group: "反馈", Demo: HoverCardDemo },
  { slug: "tooltip", title: "Tooltip", description: "提示", group: "反馈", Demo: TooltipDemo },
  { slug: "dropdown-menu", title: "DropdownMenu", description: "菜单", group: "反馈", Demo: DropdownMenuDemo },
  { slug: "collapsible", title: "Collapsible", description: "折叠", group: "布局", Demo: CollapsibleDemo },
  { slug: "command", title: "Command", description: "命令面板", group: "布局", Demo: CommandDemo },
  { slug: "tabs", title: "Tabs", description: "标签页", group: "布局", Demo: TabsDemo },
  { slug: "scroll-area", title: "ScrollArea", description: "滚动区域", group: "布局", Demo: ScrollAreaDemo },
  { slug: "separator", title: "Separator", description: "分割线", group: "布局", Demo: SeparatorDemo },
  { slug: "table", title: "Table", description: "基础表格", group: "数据", Demo: TableDemo },
  { slug: "unified-data-table", title: "UnifiedDataTable", description: "统一表格", group: "数据", Demo: UnifiedDataTableDemo },
  { slug: "sidebar", title: "Sidebar", description: "侧边栏", group: "布局", Demo: SidebarDemo },
  { slug: "card", title: "Card", description: "卡片", group: "布局", Demo: CardDemo },
  { slug: "skeleton", title: "Skeleton", description: "骨架屏", group: "反馈", Demo: SkeletonDemo },
  { slug: "card-grid-skeleton", title: "CardGridSkeleton", description: "卡片骨架", group: "反馈", Demo: SkeletonDemo },
  { slug: "master-detail-skeleton", title: "MasterDetailSkeleton", description: "主从骨架", group: "反馈", Demo: MasterDetailSkeletonDemo },
  { slug: "data-table-skeleton", title: "DataTableSkeleton", description: "表格骨架", group: "反馈", Demo: DataTableSkeletonDemo },
  { slug: "spinner", title: "Spinner", description: "加载", group: "反馈", Demo: SpinnerDemo },
  { slug: "progress", title: "Progress", description: "进度条", group: "反馈", Demo: ProgressDemo },
  { slug: "shield-loader", title: "ShieldLoader", description: "护盾加载", group: "反馈", Demo: ShieldLoaderDemo },
  { slug: "wave-grid", title: "WaveGrid", description: "波纹背景", group: "可视化", Demo: WaveGridDemo },
  { slug: "chart", title: "Chart", description: "图表", group: "可视化", Demo: ChartDemo },
  { slug: "mermaid-diagram", title: "MermaidDiagram", description: "流程图", group: "可视化", Demo: MermaidDemo },
  { slug: "terminal", title: "Terminal", description: "终端动效", group: "高级", Demo: TerminalDemo },
  { slug: "terminal-login", title: "TerminalLogin", description: "终端登录", group: "高级", Demo: TerminalLoginDemo },
  { slug: "yaml-editor", title: "YamlEditor", description: "YAML 编辑器", group: "高级", Demo: YamlEditorDemo },
  { slug: "sonner", title: "Sonner", description: "通知", group: "反馈", Demo: ToastDemo },
  { slug: "banner", title: "Banner", description: "横幅", group: "品牌", Demo: BannerDemo },
  { slug: "status", title: "Status", description: "状态指示", group: "品牌", Demo: StatusDemo },
  { slug: "field", title: "Field", description: "字段组合", group: "表单", Demo: FieldDemo },
  { slug: "form", title: "Form", description: "表单封装", group: "表单", Demo: FormDemo },
]
