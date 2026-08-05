import { useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { Avatar } from '@/shared/ui/avatar'
import styles from './AvatarDropdown.module.scss'

const menuItems = [
  { to: '/home', label: 'Home' },
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/glossary', label: 'Glossary' },
  { to: '/leaderboard', label: 'Leaderboard' },
]

export function AvatarDropdown() {
  const [isOpen, setIsOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handlePointerDown(event: PointerEvent) {
      const container = containerRef.current
      if (!container || container.contains(event.target as Node)) {
        return
      }

      setIsOpen(false)
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    document.addEventListener('pointerdown', handlePointerDown)
    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('pointerdown', handlePointerDown)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [])

  return (
    <div ref={containerRef} className={styles.dropdown}>
      <button
        type="button"
        className={styles.trigger}
        aria-label="Open sections menu"
        aria-haspopup="menu"
        aria-expanded={isOpen}
        onClick={() => setIsOpen((value) => !value)}
      >
        <Avatar />
      </button>

      <nav
        className={`${styles.menu} ${isOpen ? styles.menuOpen : ''}`}
        aria-label="Sections"
      >
        {menuItems.map((item) => (
          <Link
            key={item.to}
            to={item.to}
            className={styles.link}
            onClick={() => setIsOpen(false)}
          >
            {item.label}
          </Link>
        ))}
      </nav>
    </div>
  )
}
