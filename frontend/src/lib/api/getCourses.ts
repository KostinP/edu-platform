import { Course } from "@/types/course"
import { API_URL } from "@/lib/constants";

export async function getCourses(): Promise<Course[]> {
  console.log("API_URL:", API_URL)

  const res = await fetch(`${API_URL}/courses`, {
    cache: "no-store",
  })

  if (!res.ok) {
    throw new Error("Failed to fetch courses")
  }

  return res.json()
}
