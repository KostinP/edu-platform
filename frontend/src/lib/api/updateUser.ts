export default async function updateUser(data: {
  name: string
  role: string
  subscribe: boolean
  email?: string
}) {
  const res = await fetch('/api/user/update', {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })

  if (!res.ok) throw new Error('Не удалось обновить профиль')
}
