import { Link, Route, Routes } from 'react-router-dom'
import ApiTest from './pages/ApiTest'

function Home() {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Home</h1>
      <nav>
        <ul>
          <li>
            <Link to="/api-test" className="text-blue-500 hover:text-blue-700">
              API Test Dashboard
            </Link>
          </li>
        </ul>
      </nav>
    </div>
  )
}

function App() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/api-test" element={<ApiTest />} />
      </Routes>
    </div>
  )
}

export default App
