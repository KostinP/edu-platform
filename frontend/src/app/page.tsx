import dynamic from "next/dynamic"
import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "Онлайн-курсы для взрослых | codesign.tech",
  description: "Современные онлайн-курсы: веб-разработка, дизайн, игры, маркетинг. Интерактивно. Структурировано. Для взрослых.",
  keywords: ["онлайн-курсы", "обучение", "дизайн", "веб-разработка", "игры", "маркетинг", "для взрослых"],
  robots: "index, follow",
  openGraph: {
    title: "codesign.tech — Онлайн-курсы для взрослых",
    description: "Учитесь удобно и эффективно. Геймификация, проверка тестов, персональные рекомендации.",
    url: "https://codesign.tech", // ← заменишь на реальный адрес
    siteName: "codesign.tech",
    type: "website",
    images: [
      {
        url: "/og-image.png", // можно создать и добавить позже
        width: 1200,
        height: 630,
        alt: "codesign.tech",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "codesign.tech — Онлайн-курсы",
    description: "Образование для взрослых. Веб, дизайн, маркетинг, игры.",
    images: ["/og-image.png"],
  },
  metadataBase: new URL("https://codesign.tech"),
}


const HeroSection = dynamic(() => import("@/components/sections/HeroSection"))
const AdvantagesSection = dynamic(() => import("@/components/sections/AdvantagesSection"))
const TopicsSection = dynamic(() => import("@/components/sections/TopicsSection"))
const CoursesSection = dynamic(() => import("@/components/sections/CoursesSection"))
const VideoSection = dynamic(() => import("@/components/sections/VideoSection"))
const TestimonialsSection = dynamic(() => import("@/components/sections/TestimonialsSection"))


export default function Home() {
  return (
    <>
      <HeroSection />
      <AdvantagesSection />
      <TopicsSection />
      <VideoSection />
      <CoursesSection />
      <TestimonialsSection />
    </>
  )
}