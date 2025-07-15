import { Button } from "@/components/ui/button"
import { CoursePreviewCard } from "@/components/cards/CoursePreviewCard"
import { popularCourses } from "@/data/popularCourses"

export default function CoursesSection() {
  return (
    <section className="py-20 px-4 bg-[var(--gray-bg)]">
      <div className="max-w-7xl mx-auto flex flex-col gap-12 items-center text-center">
        <h2 className="text-2xl md:text-3xl font-bold text-foreground mb-4">Самые популярные курсы</h2>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
          {popularCourses.map((course, i) => (
            <CoursePreviewCard key={i} {...course} />
          ))}
        </div>

        <Button variant="outline" size="lg" className="rounded-full">
          Смотреть все курсы
        </Button>
      </div>
    </section>
  )
}
