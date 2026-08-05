import { useDocumentTitle } from '@/shared/lib/use-document-title'

export function WelcomePage() {
  useDocumentTitle('Добро пожаловать')

  return (
    <main>
      <h1>Welcome</h1>
      <p>Empty page placeholder.</p>
    </main>
  )
}