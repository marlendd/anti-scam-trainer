import {
  useEffect,
  useId,
  useRef,
  useState,
} from 'react'
import {
  faBars,
  faXmark,
} from '@fortawesome/free-solid-svg-icons'
import {
  NavLink,
  useLocation,
} from 'react-router-dom'

import { Icon } from '@/shared/ui/icon'

import styles from './BurgerNavigation.module.scss'

export const navigationItems = [
  {
    to: '/training/role-selection',
    label: 'Сценарии',
  },
  {
    to: '/glossary',
    label: 'Глоссарий',
  },
  {
    to: '/leaderboard',
    label: 'Лидеры',
  },
  {
    to: '/dashboard',
    label: 'Статистика',
  },
] as const


export function BurgerNavigation() {
  const menuId = useId()
  const { pathname } = useLocation()

  const [isOpen, setIsOpen] = useState(false)

  const rootRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    setIsOpen(false)
  }, [pathname])

  useEffect(() => {
    if (!isOpen) {
      return
    }

    function handlePointerDown(event: PointerEvent) {
      const target = event.target

      if (
        target instanceof Node &&
        !rootRef.current?.contains(target)
      ) {
        setIsOpen(false)
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    document.addEventListener('pointerdown', handlePointerDown)
    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener(
        'pointerdown',
        handlePointerDown,
      )
      document.removeEventListener(
        'keydown',
        handleKeyDown,
      )
    }
  }, [isOpen])

  return (
    <div
      ref={rootRef}
      className={styles.root}
    >
      <button
        type="button"
        className={styles.trigger}
        aria-label={
          isOpen
            ? 'Закрыть навигацию'
            : 'Открыть навигацию'
        }
        aria-expanded={isOpen}
        aria-controls={menuId}
        onClick={() => {
          setIsOpen((currentValue) => !currentValue)
        }}
      >
        <Icon icon={isOpen ? faXmark : faBars} />
      </button>

      <nav
        id={menuId}
        className={styles.menu}
        data-open={isOpen}
        aria-label="Мобильная навигация"
      >
        {navigationItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              [
                styles.link,
                isActive ? styles.activeLink : undefined,
              ]
                .filter(Boolean)
                .join(' ')
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>
    </div>
  )
}