import React, { useState } from 'react'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="container">
      <header>
        <h1>[[.ProjectName]]</h1>
        <p>Welcome to your React + Vite + TypeScript frontend, bootstrapped with Claw-CLI.</p>
      </header>
      <main>
        <div className="card">
          <button onClick={() => setCount((count) => count + 1)}>
            count is {count}
          </button>
        </div>
      </main>
    </div>
  )
}

export default App
