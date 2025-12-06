import api from './api'

export interface User {
  id: number
  username: string
  email: string
  avatar?: string
  created_at: string
}

export interface AuthResponse {
  user: User
  token: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export const authService = {
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/register', data)
    return response.data.data
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', data)
    return response.data.data
  },

  getMe: async (): Promise<User> => {
    const response = await api.get('/auth/me')
    return response.data.data
  },
}


