import api from './api'

export interface Message {
  id: number
  conversation_id: number
  sender_id: number
  content: string
  type: 'text' | 'image' | 'file'
  media_url?: string
  status: 'sent' | 'delivered' | 'read'
  created_at: string
  sender?: {
    id: number
    username: string
    avatar?: string
  }
}

export interface Conversation {
  id: number
  type: string
  created_at: string
  updated_at: string
  participants: Array<{
    id: number
    username: string
    email: string
    avatar?: string
  }>
  last_message?: Message
  unread_count?: number
}

export interface SendMessageRequest {
  conversation_id: number
  content: string
  type: 'text' | 'image' | 'file'
  media_url?: string
}

export const chatService = {
  getConversations: async (): Promise<Conversation[]> => {
    const response = await api.get('/conversations')
    return response.data?.data || []
  },

  getConversation: async (conversationId: number): Promise<Conversation> => {
    const response = await api.get(`/conversations/${conversationId}`)
    return response.data.data
  },

  createDirectConversation: async (userId2: number): Promise<Conversation> => {
    const response = await api.post('/conversations', { user_id_2: userId2 })
    const conversation = response.data.data
    console.log('Created conversation response:', conversation)
    
    // Ensure participants is always an array
    if (!Array.isArray(conversation.participants)) {
      conversation.participants = []
    }
    
    return conversation
  },

  getMessages: async (
    conversationId: number,
    limit = 50,
    offset = 0
  ): Promise<Message[]> => {
    const response = await api.get(
      `/conversations/${conversationId}/messages?limit=${limit}&offset=${offset}`
    )
    return response.data?.data || []
  },

  sendMessage: async (data: SendMessageRequest): Promise<Message> => {
    const response = await api.post(
      `/conversations/${data.conversation_id}/messages`,
      data
    )
    return response.data.data
  },

  markAsRead: async (conversationId: number): Promise<void> => {
    await api.post(`/conversations/${conversationId}/read`)
  },

  uploadImage: async (file: File): Promise<{ url: string; file_path: string }> => {
    const formData = new FormData()
    formData.append('image', file)

    const response = await api.post('/media/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data.data
  },
}

export interface User {
  id: number
  username: string
  email: string
  avatar?: string
}

export const userService = {
  searchUsers: async (query: string): Promise<User[]> => {
    try {
      const response = await api.get(`/users/search?q=${encodeURIComponent(query)}`)
      return response.data.data || []
    } catch (error: any) {
      console.error('Search users error:', error)
      if (error.response?.status === 400) {
        return [] // Empty query, return empty array
      }
      throw error
    }
  },
}


