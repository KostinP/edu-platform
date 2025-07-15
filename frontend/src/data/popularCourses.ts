import { Course } from "@/types/course"

export const popularCourses: Course[] = [
  {
    title: "Основы веб-разработки",
    image: "/globe.svg",
    category: "Веб",
    price: "Бесплатно",
    rating: 4.8,
    lessons: 12,
    duration: "4ч 30м",
    students: 120,
    author: "Алексей Петров",
  },
  {
    title: "Графический дизайн с нуля",
    image: "/file.svg",
    category: "Дизайн",
    price: "1290 ₽",
    rating: 4.9,
    lessons: 18,
    duration: "6ч 15м",
    students: 230,
    author: "Мария Иванова",
  },
  {
    title: "Создание 2D-игр на Unity",
    image: "/window.svg",
    category: "Игры",
    price: "990 ₽",
    rating: 4.7,
    lessons: 15,
    duration: "5ч 45м",
    students: 310,
    author: "Игорь Смирнов",
  },
]
