import { FeatureBox } from "@/components/shared/FeatureBox"
import { SectionTitle } from "@/components/shared/SectionTitle"

export default function AdvantagesSection() {
  return (
    <section className="bg-[var(--gray-bg)] py-16 px-4">
      <div className="max-w-7xl mx-auto flex flex-col gap-12">
        <SectionTitle
          title="Почему стоит учиться у нас?"
          description="Потому что мы объединяем академическую строгость с современными технологиями и пониманием того как учатся взрослые. Наша платформа — это не просто набор лекций, а продуманная образовательная среда."
          className="max-w-3xl"
        />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <FeatureBox icon="🎮" title="Геймификация с научным подходом" text="Геймификация мотивирует — бейджи, рейтинги и наглядный прогресс делают обучение осознанным." />
          <FeatureBox icon="🔬" title="Качественные курсы, а не массовый контент" text="Каждый курс создан экспертами и основан на проверенных научных методиках." />
          <FeatureBox icon="👨‍🏫" title="Персональная работа с каждым студентом" text="Автоматизированная проверка экономит время, а обратная связь от преподавателя помогает расти." />
        </div>
      </div>
    </section>
  )
}
