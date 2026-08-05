import { useDocumentTitle } from '@/shared/lib/use-document-title'

export function HomePage() {
  useDocumentTitle('Главная')

  return (
    <section>
      <h1>Home</h1>
      <p>Anti-scam trainer</p>
    </section>
  )
}