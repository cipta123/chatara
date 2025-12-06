import React from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from './App'
import './index.css'

const queryClient = new QueryClient()

// Hide URL bar on mobile when in PWA mode
const isPWA = window.matchMedia('(display-mode: standalone)').matches || 
              (window.navigator as any).standalone === true ||
              document.referrer.includes('android-app://') ||
              window.location.search.includes('utm_source=pwa')

if (isPWA) {
  // Prevent URL bar from showing on scroll in Android Chrome
  let lastScrollTop = window.pageYOffset || document.documentElement.scrollTop
  
  // Prevent scroll to top which shows URL bar
  window.addEventListener('scroll', () => {
    const currentScroll = window.pageYOffset || document.documentElement.scrollTop
    // If scrolling up (which shows URL bar), prevent it
    if (currentScroll < lastScrollTop && currentScroll < 10) {
      window.scrollTo(0, 10)
    }
    lastScrollTop = currentScroll <= 0 ? 0 : currentScroll
  }, { passive: false })
  
  // Initial scroll to hide URL bar on load
  setTimeout(() => {
    if (window.pageYOffset === 0) {
      window.scrollTo(0, 1)
    }
  }, 100)
  
  // Prevent pull-to-refresh which might show URL bar
  document.body.style.overscrollBehavior = 'none'
  
  // Set viewport height to account for URL bar
  const setViewportHeight = () => {
    const vh = window.innerHeight * 0.01
    document.documentElement.style.setProperty('--vh', `${vh}px`)
  }
  setViewportHeight()
  window.addEventListener('resize', setViewportHeight)
  window.addEventListener('orientationchange', setViewportHeight)
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </React.StrictMode>,
)


