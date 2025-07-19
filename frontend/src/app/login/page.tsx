// app/login/page.tsx
'use client'

import { useEffect } from 'react'

export default function LoginPage() {
  useEffect(() => {
    const script = document.createElement('script')
    script.src = 'https://telegram.org/js/telegram-widget.js?22'
    script.setAttribute('data-telegram-login', 'codesignedtech_bot') // имя твоего бота
    script.setAttribute('data-size', 'large')
    script.setAttribute('data-userpic', 'false')
    script.setAttribute('data-onauth', 'onTelegramAuth(user)')
    script.setAttribute('data-request-access', 'write') // чтобы получить доступ к username, name
    script.async = true
    document.getElementById('telegram-login-btn')?.appendChild(script)
  }, [])

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div id="telegram-login-btn" />
    </div>
  )
}

// глобально в window определяем обработчик
declare global {
  interface Window {
    onTelegramAuth: (user: any) => void
  }
}

if (typeof window !== 'undefined') {
  window.onTelegramAuth = async (user) => {
    // Отправляем на backend
    const res = await fetch('/api/auth/telegram', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(user),
    })

    if (res.ok) {
      window.location.href = '/' // или куда нужно
    } else {
      alert('Ошибка авторизации')
    }
  }
}
