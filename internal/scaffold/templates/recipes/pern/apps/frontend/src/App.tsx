import React, { useEffect, useState } from 'react'
import axios from 'axios'

function App() {
  const [health, setHealth] = useState<any>(null)

  useEffect(() => {
    axios.get('http://localhost:8080/api/health')
      .then(res => setHealth(res.data))
      .catch(err => console.error(err))
  }, [])

  return (
    <div style={{ fontFamily: 'sans-serif', textAlign: 'center', padding: '2rem' }}>
      <h1>[[.ProjectName]] (PERN Stack)</h1>
      <p>API Status: {health ? health.status : 'Connecting...'}</p>
    </div>
  )
}

export default App