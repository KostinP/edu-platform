'use client'

import { useEffect } from 'react'

declare global {
  interface Window {
    onTelegramAuth: (user: TelegramUser) => void
  }
}

interface TelegramUser {
  id: number
  first_name: string
  last_name?: string
  username?: string
  photo_url?: string
  auth_date: string
  hash: string
}

export default function TelegramLogin() {
  useEffect(() => {
    const script = document.createElement('script')
    script.src = 'https://telegram.org/js/telegram-widget.js?7'
    script.async = true
    script.setAttribute('data-telegram-login', 'codesignedtech_bot')
    script.setAttribute('data-size', 'large')
    script.setAttribute('data-userpic', 'false')
    script.setAttribute('data-request-access', 'write')
    script.setAttribute('data-onauth', 'onTelegramAuth(user)')
    document.getElementById('telegram-login-button')?.appendChild(script)

    // Обработчик
    window.onTelegramAuth = async function (user: TelegramUser) {
      const res = await fetch('/api/auth/telegram', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user),
      })

      if (res.ok) {
        const data = await res.json()
        localStorage.setItem('token', data.token)
        window.location.href = '/' // или router.push('/dashboard')
      } else {
        alert('Ошибка авторизации через Telegram')
      }
    }
  }, [])

  return <div id="telegram-login-button" />
}
