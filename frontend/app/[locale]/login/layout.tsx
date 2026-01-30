import type { Metadata } from "next"
import { getTranslations } from "next-intl/server"

type Props = {
  params: Promise<{ locale: string }>
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { locale } = await params
  const t = await getTranslations({ locale, namespace: "auth" })

  return {
    title: t("pageTitle"),
    description: t("pageDescription"),
  }
}

/**
 * Login page layout
 * Does not include sidebar and header
 */
export default function LoginLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <>
      {/* Inline critical CSS for instant boot splash - login page only */}
      <style suppressHydrationWarning dangerouslySetInnerHTML={{ __html: `
        #boot-splash {
          position: fixed;
          inset: 0;
          z-index: 9999;
          display: flex;
          flex-direction: column;
          min-height: 100svh;
          background: #0a0a0f;
          overflow: hidden;
          transition: opacity 0.2s ease-out;
        }
        #boot-splash.hidden {
          opacity: 0;
          pointer-events: none;
        }
        /* Animated gradient background */
        #boot-splash .bg-gradient {
          position: fixed;
          inset: 0;
          overflow: hidden;
        }
        #boot-splash .bg-blob {
          position: absolute;
          top: -50%;
          left: -50%;
          width: 200%;
          height: 200%;
          opacity: 0.3;
          background: conic-gradient(from 0deg at 50% 50%, #0a0a0f 0deg, #3f3f46 60deg, #0a0a0f 120deg, #52525b 180deg, #0a0a0f 240deg, #3f3f46 300deg, #0a0a0f 360deg);
          filter: blur(80px);
          animation: boot-blob 20s linear infinite;
        }
        #boot-splash .bg-overlay {
          position: absolute;
          inset: 0;
          background: rgba(10,10,15,0.8);
        }
        /* Grid background */
        #boot-splash .bg-grid {
          position: fixed;
          inset: 0;
          opacity: 0.4;
          background-image:
            linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px),
            linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px);
          background-size: 50px 50px;
          -webkit-mask-image: radial-gradient(circle at center, black, transparent 80%);
          mask-image: radial-gradient(circle at center, black, transparent 80%);
        }
        /* Main content */
        #boot-splash .content {
          position: relative;
          z-index: 10;
          flex: 1;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        #boot-splash .center {
          text-align: center;
        }
        /* Logo container */
        #boot-splash .logo-container {
          position: relative;
          width: 200px;
          height: 200px;
          margin: 0 auto 40px;
        }
        #boot-splash .logo-spinner {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          border-radius: 50%;
          border: 4px solid transparent;
          border-top-color: rgba(255,255,255,0.8);
          animation: boot-spin 1s linear infinite;
        }
        #boot-splash .logo-spinner::before {
          content: '';
          position: absolute;
          top: -4px;
          left: -4px;
          right: -4px;
          bottom: -4px;
          border-radius: 50%;
          border: 4px solid transparent;
          border-top-color: rgba(200,200,200,0.6);
          animation: boot-spin 1.5s linear infinite reverse;
        }
        #boot-splash .logo-spinner::after {
          content: '';
          position: absolute;
          top: 8px;
          left: 8px;
          right: 8px;
          bottom: 8px;
          border-radius: 50%;
          border: 2px solid transparent;
          border-bottom-color: rgba(255,255,255,0.3);
          animation: boot-spin 2s linear infinite;
        }
        #boot-splash .logo {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          width: 120px;
          height: 120px;
          z-index: 10;
        }
        /* Title */
        #boot-splash .title {
          font-size: 32px;
          font-weight: bold;
          letter-spacing: -0.025em;
          margin-bottom: 8px;
        }
        #boot-splash .title-luna {
          background: linear-gradient(to bottom right, #d4d4d8, #f4f4f5);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }
        #boot-splash .title-fox {
          background: linear-gradient(to bottom right, #a1a1aa, #e4e4e7);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          background-clip: text;
        }
        /* Loading status */
        #boot-splash .status {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 12px;
          margin-top: 24px;
        }
        #boot-splash .status-spinner {
          width: 16px;
          height: 16px;
          border: 2px solid rgba(255,255,255,0.2);
          border-top-color: rgba(255,255,255,0.8);
          border-radius: 50%;
          animation: boot-spin 0.8s linear infinite;
        }
        #boot-splash .status-text {
          font-size: 14px;
          color: #6b7280;
          font-weight: 500;
        }
        /* Dots */
        #boot-splash .dots {
          display: flex;
          justify-content: center;
          gap: 6px;
          margin-top: 32px;
        }
        #boot-splash .dot {
          width: 6px;
          height: 6px;
          background: rgba(255,255,255,0.3);
          border-radius: 50%;
          animation: boot-dot 1.5s ease-in-out infinite;
        }
        #boot-splash .dot:nth-child(2) { animation-delay: 0.2s; }
        #boot-splash .dot:nth-child(3) { animation-delay: 0.4s; }
        #boot-splash .dot:nth-child(4) { animation-delay: 0.6s; }
        #boot-splash .dot:nth-child(5) { animation-delay: 0.8s; }
        @keyframes boot-spin {
          to { transform: rotate(360deg); }
        }
        @keyframes boot-dot {
          0%, 100% { background: rgba(255,255,255,0.3); transform: scale(1); }
          50% { background: rgba(255,255,255,0.8); transform: scale(1.3); }
        }
        @keyframes boot-blob {
          0% { transform: translate(0,0) rotate(0deg); }
          33% { transform: translate(2%,2%) rotate(120deg); }
          66% { transform: translate(-2%,2%) rotate(240deg); }
          100% { transform: translate(0,0) rotate(360deg); }
        }
      `}} />
      {/* Inline boot splash - shows only on login route */}
      <div id="boot-splash" suppressHydrationWarning>
        <div className="bg-gradient">
          <div className="bg-blob"></div>
          <div className="bg-overlay"></div>
        </div>
        <div className="bg-grid"></div>
        <div className="content">
          <div className="center">
            <div className="logo-container">
              <div className="logo-spinner"></div>
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img className="logo" src="/images/icon-256.png" alt="" />
            </div>
            <h1 className="title">
              <span className="title-luna">Luna</span>
              <span className="title-fox">Fox</span>
            </h1>
            <div className="status">
              <div className="status-spinner"></div>
              <span className="status-text">Initializing...</span>
            </div>
            <div className="dots">
              <div className="dot"></div>
              <div className="dot"></div>
              <div className="dot"></div>
              <div className="dot"></div>
              <div className="dot"></div>
            </div>
          </div>
        </div>
      </div>
      {children}
    </>
  )
}
