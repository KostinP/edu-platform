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

// ✅ Добавим локальное расширение Window:
interface CustomWindow extends Window {
  onTelegramAuth: (user: TelegramUser) => void
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

    // ✅ Приводим window к нужному типу
    const customWindow = window as CustomWindow
    customWindow.onTelegramAuth = async function (user: TelegramUser) {
      const res = await fetch('/api/auth/telegram', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user),
      })

      if (res.ok) {
        const data = await res.json()
        localStorage.setItem('token', data.token)
        window.location.href = '/'
      } else {
        alert('Ошибка авторизации через Telegram')
      }
    }
  }, [])

  return <div id="telegram-login-button" />
}
