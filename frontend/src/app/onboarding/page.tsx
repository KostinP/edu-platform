'use client'

import OnboardingForm from '@/components/onboarding/Form'

export default function OnboardingPage() {
  return (
    <div className="max-w-xl mx-auto mt-10 p-4">
      <h1 className="text-2xl font-bold mb-4">Добро пожаловать!</h1>
      <p className="text-muted-foreground mb-6">Ответьте на несколько вопросов, чтобы мы лучше понимали ваш профиль.</p>
      <OnboardingForm />
    </div>
  )
}