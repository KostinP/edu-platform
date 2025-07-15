import Image from "next/image"
import { Button } from "@/components/ui/button"
import { Star } from "lucide-react"
import type { Course } from "@/types/course"

export function CoursePreviewCard(props: Course) {
  const {
    title, image, category, price, rating, lessons,
    duration, students, author
  } = props

  return (
    <div className="bg-white border border-border rounded-2xl overflow-hidden shadow-sm flex flex-col">
      <div className="relative h-48 w-full">
        <Image src={image} alt={title} fill className="object-cover" />
        <div className="absolute top-2 left-2 bg-white/90 text-xs px-3 py-1 rounded-full">{category}</div>
        <div className="absolute top-2 right-2 bg-white/90 text-xs px-3 py-1 rounded-full">{price}</div>
        <div className="absolute bottom-2 left-2 bg-yellow-400 text-xs px-2 py-1 rounded-full flex items-center gap-1">
          <Star className="w-3 h-3" />
          {rating.toFixed(1)}
        </div>
      </div>

      <div className="p-4 flex flex-col gap-3 flex-1">
        <h3 className="text-base font-semibold text-foreground">{title}</h3>
        <ul className="text-sm text-muted-foreground space-y-1">
          <li>{lessons} уроков · {duration}</li>
          <li>{students}+ студентов</li>
          <li>Автор: {author}</li>
        </ul>
        <Button variant="outline" size="sm" className="mt-auto rounded-full">Посмотреть</Button>
      </div>
    </div>
  )
}
