import './App.css'
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import Register from './pages/Register'
import MainPage from './pages/MainPage'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { DocumentProvider } from './context/DocumentContext'

export const API_URL = import.meta.env.MODE === "development" ? "http://localhost:8080/api/v1" : "/api/v1"

// Create a client for React Query
const queryClient = new QueryClient()

function App() {
  // Simple auth check - if token exists, consider user logged in
  const isAuthenticated = !!localStorage.getItem('token')

  return (
    <QueryClientProvider client={queryClient}>
      <DocumentProvider>
        <Router>
          <Routes>
            <Route path="/login" element={!isAuthenticated ? <Login /> : <Navigate to="/" replace />} />
            <Route path="/register" element={!isAuthenticated ? <Register /> : <Navigate to="/" replace />} />
            <Route path="/" element={isAuthenticated ? <MainPage /> : <Navigate to="/login" replace />} />
          </Routes>
        </Router>
      </DocumentProvider>
    </QueryClientProvider>
  )
}

export default App
