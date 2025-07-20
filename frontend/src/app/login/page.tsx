// app/login/page.tsx
'use client'

import { useEffect } from 'react'

interface TelegramUser {
  id: number
  first_name: string
  last_name?: string
  username?: string
  photo_url?: string
  auth_date: number
  hash: string
}


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
    onTelegramAuth: (user: TelegramUser) => void
  }
}

if (typeof window !== 'undefined') {
  window.onTelegramAuth = async (user) => {
    const payload = {
      telegram_id: String(user.id),
      first_name: user.first_name,
      last_name: user.last_name || '',
      username: user.username || '',
      photo_url: user.photo_url || '',
      email: null, // Telegram не предоставляет email
      subscribe_to_newsletter: false, // дефолт
      role: 'student', // по умолчанию, можно кастомизировать
    }

    const res = await fetch('/api/auth/telegram', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    if (res.ok) {
      window.location.href = '/'
    } else {
      alert('Ошибка авторизации')
    }
  }
}

