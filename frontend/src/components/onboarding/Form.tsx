'use client'

import { useState } from 'react'
import StepRole from './Step2Role'
import StepEmail from './Step3Email'
import { Button } from '@/components/ui/button'
import updateUser from '@/lib/api/updateUser'

export default function OnboardingForm() {
  const [step, setStep] = useState(0)
  const [role, setRole] = useState('')
  const [email, setEmail] = useState('')
  const [subscribe, setSubscribe] = useState(false)

  const next = () => setStep((s) => s + 1)
  const prev = () => setStep((s) => Math.max(0, s - 1))

  const handleSubmit = async () => {
    await updateUser({ role, email: subscribe ? email : null, subscribed: subscribe })
    window.location.href = '/' // или /dashboard
  }

  return (
    <div className="space-y-6">
      {step === 0 && <StepRole value={role} onChange={setRole} />}
      {step === 1 && <StepEmail email={email} subscribed={subscribe} onEmailChange={setEmail} onSubscribeChange={setSubscribe} />}

      <div className="flex justify-between">
        {step > 0 && <Button variant="outline" onClick={prev}>Назад</Button>}
        {step < 1 ? (
          <Button onClick={next} disabled={!role}>Далее</Button>
        ) : (
          <Button onClick={handleSubmit} disabled={!role}>Завершить</Button>
        )}
      </div>
    </div>
  )
}