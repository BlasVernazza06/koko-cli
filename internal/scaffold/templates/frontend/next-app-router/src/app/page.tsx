import React from 'react';

export default function Home() {
  return (
    <main style={{ padding: '2rem', fontFamily: 'sans-serif', textAlign: 'center' }}>
      <h1>{{.ProjectName}}</h1>
      <p>Welcome to your Next.js frontend, bootstrapped with Claw-CLI.</p>
    </main>
  );
}
