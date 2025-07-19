import { NextRequest, NextResponse } from 'next/server'

export async function POST(req: NextRequest) {
  const body = await req.json()

  const res = await fetch(`${process.env.API_URL}/auth/telegram`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    return NextResponse.json({ error: 'Auth failed' }, { status: 401 })
  }

  const data = await res.json()

  // здесь можно сохранить токен в куки, если нужно
  return NextResponse.redirect('/')
}
