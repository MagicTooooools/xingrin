"use client"

import { useState, useEffect, useRef, useCallback } from "react"
import { useTranslations } from "next-intl"
import {
  Dialog,
  DialogContent,
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
import { Button } from "@/components/ui/button"
import { IconRocket, IconEye, IconTrash, IconRefresh } from "@tabler/icons-react"
import type { WorkerNode } from "@/types/worker.types"

interface DeployTerminalDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  worker: WorkerNode | null
  onDeployComplete?: () => void
}

// Auto-generate WebSocket URL based on current page URL
const getWsBaseUrl = () => {
  if (typeof window === 'undefined') return 'ws://localhost:8080'
  
  // Prefer environment variable
  if (process.env.NEXT_PUBLIC_WS_URL) {
    return process.env.NEXT_PUBLIC_WS_URL
  }
  
  // Auto-generate based on current page protocol and domain
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}`
}

export function DeployTerminalDialog({
  open,
  onOpenChange,
  worker,
  onDeployComplete,
}: DeployTerminalDialogProps) {
  const t = useTranslations("settings.workers")
  const tCommon = useTranslations("common.actions")
  const tTerminal = useTranslations("settings.workers.terminal")
  const [isConnected, setIsConnected] = useState(false)
  const [error, setError] = useState<string | null>(null)
  // Local worker state for real-time button display updates
  const [localStatus, setLocalStatus] = useState<string | null>(null)
  const [uninstallDialogOpen, setUninstallDialogOpen] = useState(false)
  
  // Use local state or passed worker state
  const currentStatus = localStatus || worker?.status
  const terminalRef = useRef<HTMLDivElement>(null)
  const terminalInstanceRef = useRef<any>(null)
  const fitAddonRef = useRef<any>(null)
  const wsRef = useRef<WebSocket | null>(null)

  // Initialize xterm
  const initTerminal = useCallback(async () => {
    if (!terminalRef.current || terminalInstanceRef.current) return
    
    const { Terminal } = await import('@xterm/xterm')
    const { FitAddon } = await import('@xterm/addon-fit')
    const { WebLinksAddon } = await import('@xterm/addon-web-links')
    
    const terminal = new Terminal({
      cursorBlink: true,
      fontSize: 12, // Reduced font size
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
      theme: {
        background: '#1a1b26',
        foreground: '#a9b1d6',
        cursor: '#c0caf5',
        black: '#32344a',
        red: '#f7768e',
        green: '#9ece6a',
        yellow: '#e0af68',
        blue: '#7aa2f7',
        magenta: '#ad8ee6',
        cyan: '#449dab',
        white: '#787c99',
      },
    })
    
    const fitAddon = new FitAddon()
    terminal.loadAddon(fitAddon)
    terminal.loadAddon(new WebLinksAddon())
    
    terminal.open(terminalRef.current)
    fitAddon.fit()
    
    terminalInstanceRef.current = terminal
    fitAddonRef.current = fitAddon
    
    // Show connection prompt
    terminal.writeln(`\x1b[90m${tTerminal("connecting")}\x1b[0m`)

    // Auto-connect WebSocket
    connectWs()
  }, [worker])

  // Manage window resize listener separately for proper cleanup
  useEffect(() => {
    const fitAddon = fitAddonRef.current
    if (!fitAddon) return

    const handleResize = () => fitAddon.fit()
    window.addEventListener('resize', handleResize)

    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [fitAddonRef.current])

  // Connect WebSocket
  const connectWs = useCallback(() => {
    if (!worker || !terminalInstanceRef.current) return
    
    const terminal = terminalInstanceRef.current
    // Close existing connection first
    if (wsRef.current) {
        wsRef.current.close()
    }
    
    const ws = new WebSocket(`${getWsBaseUrl()}/ws/workers/${worker.id}/deploy/`)
    ws.binaryType = 'arraybuffer'
    wsRef.current = ws
    
    ws.onopen = () => {
      terminal.writeln(`\x1b[32m✓ ${tTerminal("wsConnected")}\x1b[0m`)
      // Backend will auto-start SSH connection
    }
    
    ws.onmessage = (event) => {
      if (event.data instanceof ArrayBuffer) {
        // Binary data - terminal output
        const decoder = new TextDecoder()
        terminal.write(decoder.decode(event.data))
      } else {
        // JSON message
        try {
          const data = JSON.parse(event.data)
          if (data.type === 'connected') {
            setIsConnected(true)
            setError(null)
            // Bind terminal input
            terminal.onData((data: string) => {
              if (ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({ type: 'input', data }))
              }
            })
            // Send terminal size
            ws.send(JSON.stringify({
                type: 'resize',
                cols: terminal.cols,
                rows: terminal.rows,
            }))
          } else if (data.type === 'error') {
            terminal.writeln(`\x1b[31m✗ ${data.message}\x1b[0m`)
            setError(data.message)
          } else if (data.type === 'status') {
            // Update local state to show correct buttons in real-time
            setLocalStatus(data.status)
            // Refresh parent component list on any status change
            onDeployComplete?.()
          }
        } catch {
          // Ignore parse errors
        }
      }
    }
    
    ws.onclose = () => {
      terminal.writeln('')
      terminal.writeln(`\x1b[33m${tTerminal("connectionClosed")}\x1b[0m`)
      setIsConnected(false)
    }
    
    ws.onerror = () => {
      terminal.writeln(`\x1b[31m✗ ${tTerminal("wsConnectionFailed")}\x1b[0m`)
      setError(tTerminal("connectionFailed"))
    }
  }, [worker, onDeployComplete])

  // Send terminal size change
  useEffect(() => {
    if (!isConnected || !wsRef.current || !terminalInstanceRef.current) return
    
    const terminal = terminalInstanceRef.current
    const handleResize = () => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({
          type: 'resize',
          cols: terminal.cols,
          rows: terminal.rows,
        }))
      }
    }
    
    terminal.onResize?.(handleResize)
  }, [isConnected])

  // Initialize when opened
  useEffect(() => {
    if (open && worker) {
      // Delay initialization to ensure DOM is rendered
      const timer = setTimeout(initTerminal, 100)
      return () => clearTimeout(timer)
    }
  }, [open, worker, initTerminal])

  // Cleanup when closed
  const handleClose = () => {
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
    if (terminalInstanceRef.current) {
      terminalInstanceRef.current.dispose()
      terminalInstanceRef.current = null
    }
    fitAddonRef.current = null
    setIsConnected(false)
    setError(null)
    setLocalStatus(null) // Reset local state
    // Refresh parent component list when closing to ensure state sync
    onDeployComplete?.()
    onOpenChange(false)
  }

  // Execute deploy script (runs in background)
  const handleDeploy = () => {
    if (!wsRef.current || !isConnected) return
    setLocalStatus('deploying') // Immediately update to deploying state
    onDeployComplete?.() // Refresh parent component list
    wsRef.current.send(JSON.stringify({ type: 'deploy' }))
  }

  // View deploy progress (attach to tmux session)
  const handleAttach = () => {
    if (!wsRef.current || !isConnected) return
    wsRef.current.send(JSON.stringify({ type: 'attach' }))
  }

  // Uninstall Agent (open confirmation dialog)
  const handleUninstallClick = () => {
    if (!wsRef.current || !isConnected) return
    setUninstallDialogOpen(true)
  }

  // Confirm uninstall
  const handleUninstallConfirm = () => {
    if (!wsRef.current || !isConnected) return
    setUninstallDialogOpen(false)
    wsRef.current.send(JSON.stringify({ type: 'uninstall' }))
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="w-[50vw] max-w-[50vw] h-[80vh] flex flex-col p-0 gap-0 overflow-hidden [&>button]:hidden">
        {/* Terminal title bar - macOS style */}
        <div className="flex items-center justify-between px-4 py-3 bg-[#1a1b26] border-b border-[#32344a]">
          <div className="flex items-center gap-3">
            {/* Red, yellow, green buttons */}
            <div className="flex items-center gap-1.5">
              <button 
                onClick={handleClose}
                className="w-3 h-3 rounded-full bg-[#ff5f56] hover:bg-[#ff5f56]/80 transition-colors"
                title={tCommon("close")}
              />
              <div className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
              <div className="w-3 h-3 rounded-full bg-[#27c93f]" />
            </div>
            {/* Title */}
            <span className="text-sm text-[#a9b1d6] font-medium">
              {worker?.username}@{worker?.ipAddress}
            </span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-[#9ece6a]' : 'bg-[#f7768e]'}`} />
            <span className="text-xs text-[#a9b1d6]">{isConnected ? t("terminal.connected") : t("terminal.disconnected")}</span>
          </div>
        </div>

        {/* xterm terminal container */}
        <div 
          ref={terminalRef} 
          className="flex-1 overflow-hidden bg-[#1a1b26]"
        />

        {/* Bottom action bar - show different buttons based on status */}
        <div className="flex items-center justify-between px-4 py-3 bg-[#1a1b26] border-t border-[#32344a]">
          {/* Left: Status hint */}
          <div className="text-xs text-[#565f89]">
            {!isConnected && tTerminal("waitingConnection")}
            {isConnected && currentStatus === 'pending' && tTerminal("pendingHint")}
            {isConnected && currentStatus === 'deploying' && tTerminal("deployingHint")}
            {isConnected && currentStatus === 'online' && tTerminal("onlineHint")}
            {isConnected && currentStatus === 'offline' && tTerminal("offlineHint")}
            {isConnected && currentStatus === 'updating' && tTerminal("updatingHint")}
            {isConnected && currentStatus === 'outdated' && tTerminal("outdatedHint")}
          </div>
          
          {/* Right: Action buttons */}
          <div className="flex items-center gap-2">
            {!isConnected && (
              <button 
                onClick={connectWs}
                className="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-[#32344a] text-[#a9b1d6] hover:bg-[#414868] transition-colors"
              >
                <IconRefresh className="mr-1.5 h-4 w-4" />
                {tTerminal("reconnect")}
              </button>
            )}
            {isConnected && worker && (
              <>
                {/* Not deployed -> Show "Start Deploy" */}
                {currentStatus === 'pending' && (
                  <button 
                    onClick={handleDeploy}
                    className="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-[#7aa2f7] text-[#1a1b26] hover:bg-[#7aa2f7]/80 transition-colors"
                  >
                    <IconRocket className="mr-1.5 h-4 w-4" />
                    {tTerminal("startDeploy")}
                  </button>
                )}
                
                {/* Deploying -> Show "View Progress" */}
                {currentStatus === 'deploying' && (
                  <button 
                    onClick={handleAttach}
                    className="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-[#7aa2f7] text-[#1a1b26] hover:bg-[#7aa2f7]/80 transition-colors"
                  >
                    <IconEye className="mr-1.5 h-4 w-4" />
                    {tTerminal("viewProgress")}
                  </button>
                )}
                
                {/* Updating -> Show "View Progress" */}
                {currentStatus === 'updating' && (
                  <button 
                    onClick={handleAttach}
                    className="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-[#e0af68] text-[#1a1b26] hover:bg-[#e0af68]/80 transition-colors"
                  >
                    <IconEye className="mr-1.5 h-4 w-4" />
                    {tTerminal("viewProgress")}
                  </button>
                )}
                
                {/* Version outdated -> Show "Redeploy" */}
                {currentStatus === 'outdated' && (
                  <button 
                    onClick={handleDeploy}
                    className="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-[#f7768e] text-[#1a1b26] hover:bg-[#f7768e]/80 transition-colors"
                  >
                    <IconRocket className="mr-1.5 h-4 w-4" />
                    {tTerminal("redeploy")}
                  </button>
                )}
                
                {/* Deployed (online/offline) -> Show "Redeploy" and "Uninstall" */}
                {(currentStatus === 'online' || currentStatus === 'offline') && (
                  <>
                    <button 
                      onClick={handleDeploy}
                      className="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-[#32344a] text-[#a9b1d6] hover:bg-[#414868] transition-colors"
                    >
                      <IconRocket className="mr-1.5 h-4 w-4" />
                      {tTerminal("redeploy")}
                    </button>
                    <button 
                      onClick={handleUninstallClick}
                      className="inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-[#32344a] text-[#f7768e] hover:bg-[#414868] transition-colors"
                    >
                      <IconTrash className="mr-1.5 h-4 w-4" />
                      {tTerminal("uninstall")}
                    </button>
                  </>
                )}
              </>
            )}
          </div>
        </div>
      </DialogContent>

      {/* Uninstall confirmation dialog */}
      <AlertDialog open={uninstallDialogOpen} onOpenChange={setUninstallDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t("confirmUninstall")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("confirmUninstallDesc")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("cancel")}</AlertDialogCancel>
            <AlertDialogAction 
              onClick={handleUninstallConfirm}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {tCommon("delete")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </Dialog>
  )
}
