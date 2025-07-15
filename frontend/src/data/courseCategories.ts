import { CourseCategory } from "@/types/category"
import { COLORS } from "@/lib/constants"

export const courseCategories: CourseCategory[] = [
  {
    icon: "🎨",
    title: "Дизайн",
    description: "Визуальная эстетика — создавайте продукты, которые впечатляют",
    colorVar: COLORS.yellow,
  },
  {
    icon: "🌐",
    title: "Веб-разработка",
    description: "От верстки до fullstack — освоите технологии современных сайтов",
    colorVar: COLORS.violet,
  },
  {
    icon: "🕹️",
    title: "Разработка игр",
    description: "От идеи до релиза — программируйте игры на Unity, Unreal и не только",
    colorVar: COLORS.blue,
  },
  {
    icon: "📈",
    title: "Маркетинг",
    description: "От SMM до аналитики — учитесь продвигать продукты в цифровой среде",
    colorVar: COLORS.orange,
  },
]
