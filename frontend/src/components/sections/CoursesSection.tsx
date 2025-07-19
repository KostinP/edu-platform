"use client"

import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { CoursePreviewCard } from "@/components/cards/CoursePreviewCard"
import { getCourses } from "@/lib/api/getCourses"
import { Course } from "@/types/course"

export default function CoursesSection() {
  const [courses, setCourses] = useState<Course[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function loadCourses() {
      try {
        const rawCourses = await getCourses()
        const adapted = rawCourses.map((c) => ({
          title: c.title,
          image: "/placeholder-course.jpg", // или c.image, если появится
          category: "Разработка", // заглушка, если нет категорий
          price: "Бесплатно", // или динамически
          rating: 4.8, // можно динамически позже
          lessons: 10, // заглушка
          duration: "4ч", // заглушка
          students: 120, // заглушка
          author: "Автор курса", // заглушка
        }))
        setCourses(adapted)
      } catch (e) {
        console.error("Не удалось загрузить курсы", e)
      } finally {
        setLoading(false)
      }
    }

    loadCourses()
  }, [])

  return (
    <section className="py-20 px-4 bg-[var(--gray-bg)]">
      <div className="max-w-7xl mx-auto flex flex-col gap-12 items-center text-center">
        <h2 className="text-2xl md:text-3xl font-bold text-foreground mb-4">Самые популярные курсы</h2>

        {loading ? (
          <p className="text-muted-foreground">Загрузка курсов...</p>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
            {courses.map((course, i) => (
              <CoursePreviewCard key={i} {...course} />
            ))}
          </div>
        )}

        <Button variant="outline" size="lg" className="rounded-full">
          Смотреть все курсы
        </Button>
      </div>
    </section>
  )
}
