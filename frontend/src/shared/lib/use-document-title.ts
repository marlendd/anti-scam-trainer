import { useEffect } from 'react'

const APP_TITLE = 'Авито'

export function useDocumentTitle(title: string) {
  useEffect(() => {
    document.title = `${title} | ${APP_TITLE}`
  }, [title])
}