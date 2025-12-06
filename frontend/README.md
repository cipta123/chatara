# UMess Frontend

React frontend for UMess messaging application.

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create environment file (optional):
   ```bash
   # Create .env file
   VITE_API_URL=http://localhost:8080/api
   VITE_WS_URL=ws://localhost:8080/ws
   ```

3. Start development server:
   ```bash
   npm run dev
   ```

4. Build for production:
   ```bash
   npm run build
   ```

## Features

- User authentication (Login/Register)
- Real-time messaging via WebSocket
- Chat list with last message preview
- Message bubbles with status indicators
- Image upload and sharing
- Responsive design

## Project Structure

- `src/components/` - Reusable React components
- `src/pages/` - Page components
- `src/services/` - API service functions
- `src/hooks/` - Custom React hooks


