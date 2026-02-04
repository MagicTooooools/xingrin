import React, { forwardRef } from "react"
import type { CarbonIconType, CarbonIconProps } from "@carbon/icons-react"
import {
  Activity as CarbonActivity,
  Add as CarbonAdd,
  Alarm as CarbonAlarm,
  ApplicationWeb as CarbonApplicationWeb,
  ArrowLeft as CarbonArrowLeft,
  ArrowRight as CarbonArrowRight,
  ArrowUp as CarbonArrowUp,
  ArrowUpRight as CarbonArrowUpRight,
  Book as CarbonBook,
  Box as CarbonBox,
  Branch as CarbonBranch,
  Building as CarbonBuilding,
  Calendar as CarbonCalendar,
  Camera as CarbonCamera,
  ChartBar as CarbonChartBar,
  ChartLine as CarbonChartLine,
  Checkmark as CarbonCheckmark,
  CheckmarkOutline as CarbonCheckmarkOutline,
  ChevronDown as CarbonChevronDown,
  ChevronLeft as CarbonChevronLeft,
  ChevronRight as CarbonChevronRight,
  ChevronSort as CarbonChevronSort,
  ChevronUp as CarbonChevronUp,
  Chip as CarbonChip,
  CircleFilled as CarbonCircleFilled,
  CircleOutline as CarbonCircleOutline,
  Close as CarbonClose,
  Cloud as CarbonCloud,
  CloudOffline as CarbonCloudOffline,
  Code as CarbonCode,
  Column as CarbonColumn,
  Copy as CarbonCopy,
  Dashboard as CarbonDashboard,
  DataBase as CarbonDataBase,
  DataCenter as CarbonDataCenter,
  Debug as CarbonDebug,
  Document as CarbonDocument,
  Download as CarbonDownload,
  Earth as CarbonEarth,
  Edit as CarbonEdit,
  Email as CarbonEmail,
  ErrorOutline as CarbonErrorOutline,
  Filter as CarbonFilter,
  FingerprintRecognition as CarbonFingerprintRecognition,
  Flash as CarbonFlash,
  Folder as CarbonFolder,
  FolderDetails as CarbonFolderDetails,
  FolderOpen as CarbonFolderOpen,
  Globe as CarbonGlobe,
  Grid as CarbonGrid,
  Hashtag as CarbonHashtag,
  Help as CarbonHelp,
  HexagonOutline as CarbonHexagonOutline,
  Image as CarbonImage,
  Information as CarbonInformation,
  Key as CarbonKey,
  Launch as CarbonLaunch,
  Layers as CarbonLayers,
  Link as CarbonLink,
  List as CarbonList,
  Location as CarbonLocation,
  Locked as CarbonLocked,
  LogoDiscord as CarbonLogoDiscord,
  LogoGithub as CarbonLogoGithub,
  LogoSlack as CarbonLogoSlack,
  Logout as CarbonLogout,
  Misuse as CarbonMisuse,
  Network_1 as CarbonNetwork_1,
  Notification as CarbonNotification,
  NotificationOff as CarbonNotificationOff,
  OverflowMenuHorizontal as CarbonOverflowMenuHorizontal,
  OverflowMenuVertical as CarbonOverflowMenuVertical,
  Package as CarbonPackage,
  PauseOutline as CarbonPauseOutline,
  Play as CarbonPlay,
  Radar as CarbonRadar,
  RadarEnhanced as CarbonRadarEnhanced,
  Radio as CarbonRadio,
  RadioButtonChecked as CarbonRadioButtonChecked,
  RecentlyViewed as CarbonRecentlyViewed,
  Renew as CarbonRenew,
  Report as CarbonReport,
  Restart as CarbonRestart,
  Rocket as CarbonRocket,
  Save as CarbonSave,
  Scan as CarbonScan,
  Screen as CarbonScreen,
  Search as CarbonSearch,
  SearchLocate as CarbonSearchLocate,
  Security as CarbonSecurity,
  SecurityServices as CarbonSecurityServices,
  Settings as CarbonSettings,
  SettingsAdjust as CarbonSettingsAdjust,
  SidePanelOpen as CarbonSidePanelOpen,
  Sight as CarbonSight,
  SignalStrength as CarbonSignalStrength,
  StoragePool as CarbonStoragePool,
  Tag as CarbonTag,
  Terminal as CarbonTerminal,
  TextAlignLeft as CarbonTextAlignLeft,
  Time as CarbonTime,
  Tools as CarbonTools,
  Translate as CarbonTranslate,
  TrashCan as CarbonTrashCan,
  Upload as CarbonUpload,
  User as CarbonUser,
  UserMultiple as CarbonUserMultiple,
  View as CarbonView,
  ViewOff as CarbonViewOff,
  Warning as CarbonWarning,
  WarningAlt as CarbonWarningAlt,
  Waveform as CarbonWaveform
} from "@carbon/icons-react"

const withCarbon = (Icon: CarbonIconType) => {
  const Wrapped = forwardRef<React.ReactSVGElement, CarbonIconProps>(({ size = 16, ...props }, ref) => (
    <Icon ref={ref} size={size} {...props} />
  ))
  Wrapped.displayName = Icon.displayName ?? Icon.name
  return Wrapped
}

export type Icon = CarbonIconType
export type LucideIcon = CarbonIconType

export const Activity = withCarbon(CarbonActivity)
export const AlarmClock = withCarbon(CarbonAlarm)
export const AlertCircle = withCarbon(CarbonWarningAlt)
export const AlertTriangle = withCarbon(CarbonWarning)
export const AlignLeft = withCarbon(CarbonTextAlignLeft)
export const ArrowLeft = withCarbon(CarbonArrowLeft)
export const ArrowRight = withCarbon(CarbonArrowRight)
export const ArrowUpRight = withCarbon(CarbonArrowUpRight)
export const Ban = withCarbon(CarbonMisuse)
export const BarChart3 = withCarbon(CarbonChartBar)
export const Bell = withCarbon(CarbonNotification)
export const BellOff = withCarbon(CarbonNotificationOff)
export const Box = withCarbon(CarbonBox)
export const Building2 = withCarbon(CarbonBuilding)
export const Calendar = withCarbon(CarbonCalendar)
export const Camera = withCarbon(CarbonCamera)
export const Check = withCarbon(CarbonCheckmark)
export const CheckCircle = withCarbon(CarbonCheckmarkOutline)
export const CheckCircle2 = withCarbon(CarbonCheckmarkOutline)
export const CheckIcon = withCarbon(CarbonCheckmark)
export const ChevronDown = withCarbon(CarbonChevronDown)
export const ChevronDownIcon = withCarbon(CarbonChevronDown)
export const ChevronLeft = withCarbon(CarbonChevronLeft)
export const ChevronLeftIcon = withCarbon(CarbonChevronLeft)
export const ChevronRight = withCarbon(CarbonChevronRight)
export const ChevronRightIcon = withCarbon(CarbonChevronRight)
export const ChevronUp = withCarbon(CarbonChevronUp)
export const ChevronUpIcon = withCarbon(CarbonChevronUp)
export const ChevronsUpDown = withCarbon(CarbonChevronSort)
export const Circle = withCarbon(CarbonCircleOutline)
export const CircleCheckIcon = withCarbon(CarbonCheckmarkOutline)
export const CircleIcon = withCarbon(CarbonCircleFilled)
export const Clock = withCarbon(CarbonTime)
export const Copy = withCarbon(CarbonCopy)
export const Cpu = withCarbon(CarbonChip)
export const Crosshair = withCarbon(CarbonSight)
export const Database = withCarbon(CarbonDataBase)
export const Download = withCarbon(CarbonDownload)
export const Edit = withCarbon(CarbonEdit)
export const ExternalLink = withCarbon(CarbonLaunch)
export const Eye = withCarbon(CarbonView)
export const FileCode = withCarbon(CarbonCode)
export const FileText = withCarbon(CarbonDocument)
export const Filter = withCarbon(CarbonFilter)
export const Fingerprint = withCarbon(CarbonFingerprintRecognition)
export const Folder = withCarbon(CarbonFolder)
export const FolderOpen = withCarbon(CarbonFolderOpen)
export const FolderSearch = withCarbon(CarbonFolderDetails)
export const GitBranch = withCarbon(CarbonBranch)
export const Globe = withCarbon(CarbonGlobe)
export const HardDrive = withCarbon(CarbonStoragePool)
export const Hash = withCarbon(CarbonHashtag)
export const HelpCircle = withCarbon(CarbonHelp)
export const Hexagon = withCarbon(CarbonHexagonOutline)
export const History = withCarbon(CarbonRecentlyViewed)
export const Image = withCarbon(CarbonImage)
export const Info = withCarbon(CarbonInformation)
export const InfoIcon = withCarbon(CarbonInformation)
export const Layers = withCarbon(CarbonLayers)
export const LayoutDashboard = withCarbon(CarbonDashboard)
export const LayoutGrid = withCarbon(CarbonGrid)
export const Link = withCarbon(CarbonLink)
export const Link2 = withCarbon(CarbonLink)
export const Loader2 = withCarbon(CarbonRenew)
export const Loader2Icon = withCarbon(CarbonRenew)
export const Lock = withCarbon(CarbonLocked)
export const Monitor = withCarbon(CarbonScreen)
export const MoreHorizontal = withCarbon(CarbonOverflowMenuHorizontal)
export const Network = withCarbon(CarbonNetwork_1)
export const OctagonXIcon = withCarbon(CarbonMisuse)
export const Package = withCarbon(CarbonPackage)
export const PackageOpen = withCarbon(CarbonPackage)
export const PanelLeftIcon = withCarbon(CarbonSidePanelOpen)
export const PauseCircle = withCarbon(CarbonPauseOutline)
export const Pencil = withCarbon(CarbonEdit)
export const Play = withCarbon(CarbonPlay)
export const Plus = withCarbon(CarbonAdd)
export const Radar = withCarbon(CarbonRadar)
export const Radio = withCarbon(CarbonRadio)
export const RefreshCw = withCarbon(CarbonRestart)
export const Save = withCarbon(CarbonSave)
export const Scan = withCarbon(CarbonScan)
export const Search = withCarbon(CarbonSearch)
export const Server = withCarbon(CarbonDataCenter)
export const Settings = withCarbon(CarbonSettings)
export const Shield = withCarbon(CarbonSecurity)
export const ShieldAlert = withCarbon(CarbonSecurity)
export const ShieldCheck = withCarbon(CarbonSecurityServices)
export const Signal = withCarbon(CarbonSignalStrength)
export const Sliders = withCarbon(CarbonSettingsAdjust)
export const Tag = withCarbon(CarbonTag)
export const Target = withCarbon(CarbonSight)
export const Terminal = withCarbon(CarbonTerminal)
export const Trash2 = withCarbon(CarbonTrashCan)
export const TriangleAlertIcon = withCarbon(CarbonWarning)
export const Upload = withCarbon(CarbonUpload)
export const UploadIcon = withCarbon(CarbonUpload)
export const User = withCarbon(CarbonUser)
export const Waves = withCarbon(CarbonWaveform)
export const Wrench = withCarbon(CarbonTools)
export const X = withCarbon(CarbonClose)
export const XCircle = withCarbon(CarbonErrorOutline)
export const XIcon = withCarbon(CarbonClose)
export const Zap = withCarbon(CarbonFlash)
export const IconActivity = withCarbon(CarbonActivity)
export const IconAlertTriangle = withCarbon(CarbonWarning)
export const IconArrowUp = withCarbon(CarbonArrowUp)
export const IconBan = withCarbon(CarbonMisuse)
export const IconBook = withCarbon(CarbonBook)
export const IconBrandDiscord = withCarbon(CarbonLogoDiscord)
export const IconBrandGithub = withCarbon(CarbonLogoGithub)
export const IconBrandSlack = withCarbon(CarbonLogoSlack)
export const IconBrowser = withCarbon(CarbonApplicationWeb)
export const IconBug = withCarbon(CarbonDebug)
export const IconBuilding = withCarbon(CarbonBuilding)
export const IconCheck = withCarbon(CarbonCheckmark)
export const IconChevronDown = withCarbon(CarbonChevronDown)
export const IconChevronLeft = withCarbon(CarbonChevronLeft)
export const IconChevronRight = withCarbon(CarbonChevronRight)
export const IconChevronUp = withCarbon(CarbonChevronUp)
export const IconChevronsLeft = withCarbon(CarbonChevronLeft)
export const IconChevronsRight = withCarbon(CarbonChevronRight)
export const IconCircleCheck = withCarbon(CarbonCheckmarkOutline)
export const IconCircleDot = withCarbon(CarbonRadioButtonChecked)
export const IconCircleX = withCarbon(CarbonErrorOutline)
export const IconClock = withCarbon(CarbonTime)
export const IconCloud = withCarbon(CarbonCloud)
export const IconCloudOff = withCarbon(CarbonCloudOffline)
export const IconCode = withCarbon(CarbonCode)
export const IconCpu = withCarbon(CarbonChip)
export const IconDashboard = withCarbon(CarbonDashboard)
export const IconDatabase = withCarbon(CarbonDataBase)
export const IconDotsVertical = withCarbon(CarbonOverflowMenuVertical)
export const IconDownload = withCarbon(CarbonDownload)
export const IconEdit = withCarbon(CarbonEdit)
export const IconExternalLink = withCarbon(CarbonLaunch)
export const IconEye = withCarbon(CarbonView)
export const IconEyeOff = withCarbon(CarbonViewOff)
export const IconFileText = withCarbon(CarbonDocument)
export const IconFolder = withCarbon(CarbonFolder)
export const IconHeartbeat = withCarbon(CarbonActivity)
export const IconInfoCircle = withCarbon(CarbonInformation)
export const IconKey = withCarbon(CarbonKey)
export const IconLanguage = withCarbon(CarbonTranslate)
export const IconLayoutColumns = withCarbon(CarbonColumn)
export const IconLink = withCarbon(CarbonLink)
export const IconListDetails = withCarbon(CarbonList)
export const IconLoader2 = withCarbon(CarbonRenew)
export const IconLogout = withCarbon(CarbonLogout)
export const IconMail = withCarbon(CarbonEmail)
export const IconMapPin = withCarbon(CarbonLocation)
export const IconMessageReport = withCarbon(CarbonReport)
export const IconPlayerPlay = withCarbon(CarbonPlay)
export const IconPlus = withCarbon(CarbonAdd)
export const IconRadar = withCarbon(CarbonRadar)
export const IconRadar2 = withCarbon(CarbonRadarEnhanced)
export const IconRefresh = withCarbon(CarbonRestart)
export const IconRocket = withCarbon(CarbonRocket)
export const IconScan = withCarbon(CarbonScan)
export const IconSearch = withCarbon(CarbonSearch)
export const IconServer = withCarbon(CarbonDataCenter)
export const IconSettings = withCarbon(CarbonSettings)
export const IconShield = withCarbon(CarbonSecurity)
export const IconShieldCheck = withCarbon(CarbonSecurityServices)
export const IconStack2 = withCarbon(CarbonLayers)
export const IconTarget = withCarbon(CarbonSight)
export const IconTerminal2 = withCarbon(CarbonTerminal)
export const IconTool = withCarbon(CarbonTools)
export const IconTrash = withCarbon(CarbonTrashCan)
export const IconTrendingDown = withCarbon(CarbonChartLine)
export const IconTrendingUp = withCarbon(CarbonChartLine)
export const IconUpload = withCarbon(CarbonUpload)
export const IconUsers = withCarbon(CarbonUserMultiple)
export const IconWorld = withCarbon(CarbonEarth)
export const IconWorldSearch = withCarbon(CarbonSearchLocate)
export const IconX = withCarbon(CarbonClose)
export const MagnifyingGlassIcon = withCarbon(CarbonSearch)
