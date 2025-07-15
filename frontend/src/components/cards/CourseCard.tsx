type Props = {
  icon: React.ReactNode
  title: string
  description: string
  colorVar: string
}

export function CourseCard({ icon, title, description, colorVar }: Props) {
  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-border flex flex-col gap-4">
      <div
        className="w-12 h-12 rounded-full flex items-center justify-center text-xl"
        style={{ backgroundColor: colorVar }}
      >
        {icon}
      </div>
      <h3 className="text-lg font-semibold text-foreground">{title}</h3>
      <p className="text-muted-foreground text-sm leading-relaxed flex-1">{description}</p>
      <div className="mt-auto" />
    </div>
  )
}
