import { useEffect, useRef, useState, useCallback } from 'react'

export interface WebSocketMessage {
  type: string
  message?: any
  [key: string]: any
}

export const useWebSocket = (
  url: string,
  token: string | null,
  onMessage?: (message: WebSocketMessage) => void
) => {
  const [isConnected, setIsConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const reconnectAttempts = useRef(0)
  const maxReconnectAttempts = 5
  const isConnectingRef = useRef(false)
  const onMessageRef = useRef(onMessage)

  // Keep onMessage ref updated
  useEffect(() => {
    onMessageRef.current = onMessage
  }, [onMessage])

  const connect = useCallback(() => {
    // Prevent multiple simultaneous connections
    if (!token) {
      console.log('WebSocket: No token, skipping connection')
      return
    }

    if (isConnectingRef.current) {
      console.log('WebSocket: Already connecting, skipping')
      return
    }

    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log('WebSocket: Already connected')
      return
    }

    if (wsRef.current?.readyState === WebSocket.CONNECTING) {
      console.log('WebSocket: Connection in progress')
      return
    }

    // Close any existing connection first
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }

    isConnectingRef.current = true

    try {
      const wsUrl = `${url}?token=${encodeURIComponent(token)}`
      console.log('WebSocket: Connecting to', url)
      const ws = new WebSocket(wsUrl)

      ws.onopen = () => {
        isConnectingRef.current = false
        setIsConnected(true)
        reconnectAttempts.current = 0
        console.log('WebSocket: Connected successfully')
      }

      ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          if (onMessageRef.current) {
            onMessageRef.current(message)
          }
        } catch (error) {
          console.error('WebSocket: Error parsing message:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('WebSocket: Error', error)
        isConnectingRef.current = false
      }

      ws.onclose = (event) => {
        console.log('WebSocket: Connection closed', event.code, event.reason)
        isConnectingRef.current = false
        setIsConnected(false)
        wsRef.current = null

        // Only attempt to reconnect if not a normal close and under max attempts
        if (event.code !== 1000 && reconnectAttempts.current < maxReconnectAttempts) {
          reconnectAttempts.current++
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.current), 30000)
          console.log(`WebSocket: Reconnecting in ${delay}ms (attempt ${reconnectAttempts.current})`)
          
          if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current)
          }
          
          reconnectTimeoutRef.current = setTimeout(() => {
            connect()
          }, delay)
        }
      }

      wsRef.current = ws
    } catch (error) {
      console.error('WebSocket: Connection error:', error)
      isConnectingRef.current = false
    }
  }, [url, token])

  const disconnect = useCallback(() => {
    console.log('WebSocket: Disconnecting')
    
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
    
    reconnectAttempts.current = maxReconnectAttempts // Prevent reconnection
    isConnectingRef.current = false
    
    if (wsRef.current) {
      wsRef.current.close(1000, 'User disconnect')
      wsRef.current = null
    }
    
    setIsConnected(false)
  }, [])

  const sendMessage = useCallback((message: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    } else {
      console.error('WebSocket: Cannot send - not connected')
    }
  }, [])

  // Connect on mount, disconnect on unmount
  useEffect(() => {
    if (token) {
      // Small delay to prevent rapid reconnections during React StrictMode
      const timeoutId = setTimeout(() => {
        connect()
      }, 100)
      
      return () => {
        clearTimeout(timeoutId)
        disconnect()
      }
    }
  }, [token]) // Only depend on token, not connect/disconnect

  return {
    isConnected,
    sendMessage,
    connect,
    disconnect,
  }
}
