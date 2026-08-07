// src/shared/ui/chat-date-separator/ChatDateSeparator.tsx

import styles from './ChatDateSeparator.module.scss'

type ChatDateSeparatorProps = {
  children: string
}

export function ChatDateSeparator({
  children,
}: ChatDateSeparatorProps) {
  return (
    <div className={styles.separator} role="separator">
      <span className={styles.label}>{children}</span>
    </div>
  )
}