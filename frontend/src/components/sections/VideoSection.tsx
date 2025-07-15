"use client"

import { useEffect, useState } from "react"
import { cn } from "@/lib/utils"

export default function VideoSection() {
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    const timeout = setTimeout(() => setVisible(true), 100)
    return () => clearTimeout(timeout)
  }, [])

  return (
    <section className="py-20 px-4 bg-[var(--cta)] text-white transition-all duration-700">
      <div
        className={cn(
          "max-w-5xl mx-auto text-center opacity-0 translate-y-4 transition-all duration-700",
          visible && "opacity-100 translate-y-0"
        )}
      >
        <h2 className="text-2xl md:text-3xl font-bold mb-4">Учитесь в удобном формате</h2>
        <p className="text-base md:text-lg text-white/80 mb-8">
          Видеообучение с примерами, практикой и автоматической проверкой
        </p>

        <div className="relative aspect-video w-full rounded-xl overflow-hidden shadow-lg border border-white/10">
          <iframe
            src="https://www.youtube.com/embed/dQw4w9WgXcQ"
            title="Demo video"
            allowFullScreen
            className="w-full h-full"
            />
        </div>
      </div>
    </section>
  )
}
