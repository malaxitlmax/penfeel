import './App.css'
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import Register from './pages/Register'

function App() {
  // Simple auth check - if token exists, consider user logged in
  const isAuthenticated = !!localStorage.getItem('token')

  return (
    <Router>
      <Routes>
        <Route path="/login" element={!isAuthenticated ? <Login /> : <Navigate to="/" replace />} />
        <Route path="/register" element={!isAuthenticated ? <Register /> : <Navigate to="/" replace />} />
        <Route path="/" element={isAuthenticated ? <div>Dashboard (To be implemented)</div> : <Navigate to="/login" replace />} />
      </Routes>
    </Router>
  )
}

export default App
