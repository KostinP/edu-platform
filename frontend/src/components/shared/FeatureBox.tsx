import { Button } from "@/components/ui/button"

type Props = {
  icon: string
  title: string
  text: string
}

export function FeatureBox({ icon, title, text }: Props) {
  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-border flex flex-col gap-4">
      <div className="text-3xl">{icon}</div>
      <h3 className="text-lg font-semibold text-foreground">{title}</h3>
      <p className="text-muted-foreground text-sm leading-relaxed">{text}</p>
      <Button variant="outline" size="sm" className="self-start rounded-full mt-auto">
        Узнать больше
      </Button>
    </div>
  )
}
