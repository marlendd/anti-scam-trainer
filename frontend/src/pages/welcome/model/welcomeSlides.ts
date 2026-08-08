import welcomeImage from '@/shared/assets/images/welcome.webp'
import scenariosImage from '@/shared/assets/images/scenarios.webp'
import pathsImage from '@/shared/assets/images/paths.webp'
import startImage from '@/shared/assets/images/start-game.webp'

export type WelcomeSlide = {
    id: string
    eyebrow: string
    title: string
    description: string
    illustration: string
    illustrationAlt: string
    accent: string
}

export const welcomeSlides: WelcomeSlide[] = [
    {
        id: 'introduction',
        eyebrow: 'Антискам-тренажёр',
        title: 'Распознавайте мошенников',
        description:
            'Пройдите интерактивные сценарии, основанные на реальных схемах обмана, и научитесь безопасно покупать и продавать товары.',
        illustration: welcomeImage,
        illustrationAlt: '',
        accent: '#00aaff',
    },
    {
        id: 'scenarios',
        eyebrow: 'Практика',
        title: 'Общайтесь как в настоящем чате',
        description:
            'Изучайте переписку, находите подозрительные сообщения и принимайте решения. После каждого сценария вы получите подробный разбор.',
        illustration: scenariosImage,
        illustrationAlt: '',
        accent: '#ff4053',
    },
    {
        id: 'paths',
        eyebrow: 'Два направления',
        title: 'Выберите свой путь',
        description:
            'Проходите сценарии покупателя и продавца. У каждого пути — отдельный прогресс, задания и мошеннические схемы.',
        illustration: pathsImage,
        illustrationAlt: '',
        accent: '#965eeb',
    },
    {
        id: 'start',
        eyebrow: 'Всё готово',
        title: 'Начнём тренировку?',
        description:
            'Выберите роль, пройдите первый сценарий и проверьте, сможете ли вы распознать мошенника до того, как станет слишком поздно.',
        illustration: startImage,
        illustrationAlt: '',
        accent: '#04e061',
    },
]
